package service

import (
	"errors"
	"strings"
	"time"

	"fe_exam/common"
	"fe_exam/model"
	"fe_exam/pkg"

	"gorm.io/gorm"
)

type ConversationSrv struct {
}

func toMessageInfo(m model.Message) MessageInfo {
	return MessageInfo{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		ReceiverID:     m.ReceiverID,
		Content:        m.Content,
		CreateTime:     FormatTime(m.CreatedAt),
	}
}

// load 校验会话存在且当前用户是成员，返回会话与对方用户 ID
func (cv *ConversationSrv) load(uid, cid int64) (model.Conversation, int64, error) {
	var conv model.Conversation
	err := model.DB.Where("id = ?", cid).First(&conv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return conv, 0, ErrConversationNotFound
	}
	if err != nil {
		return conv, 0, ErrInternal.With(err)
	}
	if conv.UserLowID != uid && conv.UserHighID != uid {
		return conv, 0, ErrConversationForbidden
	}
	peerID := conv.UserHighID
	if conv.UserHighID == uid {
		peerID = conv.UserLowID
	}
	return conv, peerID, nil
}

func reverseMessages(rows []model.Message) {
	for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
		rows[i], rows[j] = rows[j], rows[i]
	}
}

func (cv *ConversationSrv) hasOlder(cid int64, rows []model.Message) (bool, error) {
	if len(rows) == 0 {
		return false, nil
	}
	var count int64
	err := model.DB.Model(&model.Message{}).
		Where("conversation_id = ? AND id < ?", cid, rows[0].ID).
		Count(&count).Error
	if err != nil {
		return false, ErrInternal.With(err)
	}
	return count > 0, nil
}

// List 会话列表：{list,total,page,size}，含对方信息与最后一条消息
func (cv *ConversationSrv) List(uid int64, p common.PagerForm) (ConversationListResponse, error) {
	p = p.Normalize()

	base := func() *gorm.DB {
		return model.DB.Model(&model.Conversation{}).
			Where("user_low_id = ? OR user_high_id = ?", uid, uid)
	}

	var total int64
	if err := base().Count(&total).Error; err != nil {
		return ConversationListResponse{}, ErrInternal.With(err)
	}

	var rows []model.Conversation
	err := base().
		Preload("UserLow").
		Preload("UserHigh").
		Order("updated_at DESC, id DESC").
		Scopes(model.Paginate(p)).
		Find(&rows).Error
	if err != nil {
		return ConversationListResponse{}, ErrInternal.With(err)
	}

	// 批量取每个会话的最后一条消息，避免 N+1
	last := make(map[int64]MessageInfo, len(rows))
	if len(rows) > 0 {
		ids := make([]int64, 0, len(rows))
		for _, row := range rows {
			ids = append(ids, row.ID)
		}

		type maxRow struct {
			ConversationID int64
			MaxID          int64
		}
		var maxRows []maxRow
		err = model.DB.Model(&model.Message{}).
			Select("conversation_id, MAX(id) AS max_id").
			Where("conversation_id IN ?", ids).
			Group("conversation_id").
			Scan(&maxRows).Error
		if err != nil {
			return ConversationListResponse{}, ErrInternal.With(err)
		}

		msgIDs := make([]int64, 0, len(maxRows))
		for _, r := range maxRows {
			msgIDs = append(msgIDs, r.MaxID)
		}
		if len(msgIDs) > 0 {
			var msgs []model.Message
			if err := model.DB.Where("id IN ?", msgIDs).Find(&msgs).Error; err != nil {
				return ConversationListResponse{}, ErrInternal.With(err)
			}
			for _, m := range msgs {
				last[m.ConversationID] = toMessageInfo(m)
			}
		}
	}

	list := make([]ConversationInfo, 0, len(rows))
	for _, row := range rows {
		peer := row.UserHigh
		if row.UserHighID == uid {
			peer = row.UserLow
		}
		item := ConversationInfo{
			ID:         row.ID,
			Peer:       UserBrief{ID: peer.ID, Name: peer.Name, Avatar: peer.Avatar},
			CreateTime: FormatTime(row.CreatedAt),
			UpdateTime: FormatTime(row.UpdatedAt),
		}
		if m, ok := last[row.ID]; ok {
			msg := m
			item.LastMessage = &msg
		}
		list = append(list, item)
	}

	return ConversationListResponse{PagerForm: p, List: list, Total: total}, nil
}

// Create 创建或复用两人会话
func (cv *ConversationSrv) Create(uid int64, p ConversationCreateParams) (ConversationCreateResponse, error) {
	if uid == p.UserID {
		return ConversationCreateResponse{}, ErrChatSelf
	}

	var peer model.User
	err := model.DB.Where("id = ?", p.UserID).First(&peer).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ConversationCreateResponse{}, ErrUserNotFound
	}
	if err != nil {
		return ConversationCreateResponse{}, ErrInternal.With(err)
	}

	low, high := uid, p.UserID
	if low > high {
		low, high = high, low
	}

	var conv model.Conversation
	err = model.DB.Where("user_low_id = ? AND user_high_id = ?", low, high).First(&conv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		conv = model.Conversation{UserLowID: low, UserHighID: high}
		if err := model.DB.Create(&conv).Error; err != nil {
			if !isDuplicateKey(err) {
				return ConversationCreateResponse{}, ErrInternal.With(err)
			}
			if err := model.DB.Where("user_low_id = ? AND user_high_id = ?", low, high).First(&conv).Error; err != nil {
				return ConversationCreateResponse{}, ErrInternal.With(err)
			}
		}
	} else if err != nil {
		return ConversationCreateResponse{}, ErrInternal.With(err)
	}

	return ConversationCreateResponse{ConversationInfo: ConversationInfo{
		ID:         conv.ID,
		Peer:       UserBrief{ID: peer.ID, Name: peer.Name, Avatar: peer.Avatar},
		CreateTime: FormatTime(conv.CreatedAt),
		UpdateTime: FormatTime(conv.UpdatedAt),
	}}, nil
}

// Messages 会话消息：afterId 拉新、beforeId 翻历史、都不传取最新一页
func (cv *ConversationSrv) Messages(uid, cid int64, q MessageListQuery) (MessageListResponse, error) {
	_, peerID, err := cv.load(uid, cid)
	if err != nil {
		return MessageListResponse{}, err
	}

	var peer model.User
	if err := model.DB.Select("id", "name", "avatar").Where("id = ?", peerID).First(&peer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return MessageListResponse{}, ErrUserNotFound
		}
		return MessageListResponse{}, ErrInternal.With(err)
	}

	limit := q.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var rows []model.Message
	var queryErr error
	switch {
	case q.AfterID > 0:
		queryErr = model.DB.Where("conversation_id = ? AND id > ?", cid, q.AfterID).
			Order("id ASC").Limit(limit).Find(&rows).Error
	case q.BeforeID > 0:
		queryErr = model.DB.Where("conversation_id = ? AND id < ?", cid, q.BeforeID).
			Order("id DESC").Limit(limit).Find(&rows).Error
		reverseMessages(rows)
	default:
		queryErr = model.DB.Where("conversation_id = ?", cid).
			Order("id DESC").Limit(limit).Find(&rows).Error
		reverseMessages(rows)
	}
	if queryErr != nil {
		return MessageListResponse{}, ErrInternal.With(queryErr)
	}

	hasMore := false
	if q.AfterID == 0 {
		hasMore, err = cv.hasOlder(cid, rows)
		if err != nil {
			return MessageListResponse{}, err
		}
	}

	return MessageListResponse{
		List:    pkg.Map(rows, toMessageInfo),
		HasMore: hasMore,
		Peer:    UserBrief{ID: peer.ID, Name: peer.Name, Avatar: peer.Avatar},
	}, nil
}

// Send 发送消息
func (cv *ConversationSrv) Send(uid, cid int64, p MessageSendParams) (MessageSendResponse, error) {
	conv, peerID, err := cv.load(uid, cid)
	if err != nil {
		return MessageSendResponse{}, err
	}

	content := strings.TrimSpace(p.Content)
	if content == "" {
		return MessageSendResponse{}, ErrParam.WithMsg("消息内容不能为空")
	}

	msg := model.Message{
		ConversationID: conv.ID,
		SenderID:       uid,
		ReceiverID:     peerID,
		Content:        content,
	}
	txErr := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}
		return tx.Model(&model.Conversation{}).
			Where("id = ?", conv.ID).
			Update("updated_at", time.Now()).Error
	})
	if txErr != nil {
		return MessageSendResponse{}, ErrInternal.With(txErr)
	}

	return MessageSendResponse{MessageInfo: toMessageInfo(msg)}, nil
}
