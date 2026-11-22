package service

import (
	"errors"
	"strings"

	"fe_exam/model"
	"fe_exam/pkg"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentSrv struct {
}

func toCommentInfo(cm model.Comment) CommentInfo {
	return CommentInfo{
		ID:        cm.ID,
		ArticleID: cm.ArticleID,
		UserID:    cm.UserID,
		User: UserBrief{
			ID:     cm.User.ID,
			Name:   cm.User.Name,
			Avatar: cm.User.Avatar,
		},
		Content:    cm.Content,
		LikeCount:  cm.LikeCount,
		CreateTime: FormatTime(cm.CreatedAt),
	}
}

func ensureArticleExists(id int64) error {
	var count int64
	if err := model.DB.Model(&model.Article{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return ErrInternal.With(err)
	}
	if count == 0 {
		return ErrArticleNotFound
	}
	return nil
}

// List 评论列表：{list,total}，不分页（服务端上限 500）
func (c *CommentSrv) List(articleID int64, p ArticleCommentListParams) (CommentListResponse, error) {
	if err := ensureArticleExists(articleID); err != nil {
		return CommentListResponse{}, err
	}

	order := "comment.like_count DESC, comment.id DESC"
	if p.Sort == Sort("time-desc") {
		order = "comment.created_at DESC, comment.id DESC"
	}

	var total int64
	err := model.DB.Model(&model.Comment{}).
		Where("comment.article_id = ?", articleID).
		Count(&total).Error
	if err != nil {
		return CommentListResponse{}, ErrInternal.With(err)
	}

	var rows []model.Comment
	err = model.DB.Model(&model.Comment{}).
		Where("comment.article_id = ?", articleID).
		Order(order).
		Limit(500).
		Preload("User").
		Find(&rows).Error
	if err != nil {
		return CommentListResponse{}, ErrInternal.With(err)
	}

	return CommentListResponse{
		List:  pkg.Map(rows, toCommentInfo),
		Total: total,
	}, nil
}

func (c *CommentSrv) Create(uid, articleID int64, p CommentCreateParams) (CommentCreateResponse, error) {
	content := strings.TrimSpace(p.Content)
	if content == "" {
		return CommentCreateResponse{}, ErrParam.WithMsg("评论内容不能为空")
	}
	if err := ensureArticleExists(articleID); err != nil {
		return CommentCreateResponse{}, err
	}

	comment := model.Comment{
		ArticleID: articleID,
		UserID:    uid,
		Content:   content,
	}
	if err := model.DB.Create(&comment).Error; err != nil {
		return CommentCreateResponse{}, ErrInternal.With(err)
	}

	var user model.User
	if err := model.DB.Select("id", "name", "avatar").Where("id = ?", uid).First(&user).Error; err != nil {
		return CommentCreateResponse{}, ErrInternal.With(err)
	}
	comment.User = user

	return CommentCreateResponse{CommentInfo: toCommentInfo(comment)}, nil
}

func (c *CommentSrv) Delete(uid, id int64) error {
	var comment model.Comment
	err := model.DB.Where("id = ?", id).First(&comment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCommentNotFound
	}
	if err != nil {
		return ErrInternal.With(err)
	}
	if comment.UserID != uid {
		return ErrCommentForbidden
	}

	txErr := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("target_type = ? AND target_id = ?", LikeTargetComment, id).
			Delete(&model.LikeRecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Comment{}, id).Error
	})
	if txErr != nil {
		return ErrInternal.With(txErr)
	}
	return nil
}

func (c *CommentSrv) Like(uid, id int64) (CommentLikeResponse, error) {
	return c.setLike(uid, id, true)
}

func (c *CommentSrv) Unlike(uid, id int64) (CommentLikeResponse, error) {
	return c.setLike(uid, id, false)
}

func (c *CommentSrv) setLike(uid, id int64, liked bool) (CommentLikeResponse, error) {
	var resp CommentLikeResponse
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		var comment model.Comment
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&comment).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		if err != nil {
			return err
		}

		var record model.LikeRecord
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND target_id = ? AND target_type = ?", uid, id, LikeTargetComment).
			First(&record).Error

		delta := 0
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			if liked {
				record = model.LikeRecord{UserID: uid, TargetID: id, TargetType: LikeTargetComment, Status: 1}
				if err := tx.Create(&record).Error; err != nil {
					if isDuplicateKey(err) {
						return nil
					}
					return err
				}
				delta = 1
			}
		case err != nil:
			return err
		default:
			if record.Status == 1 && !liked {
				if err := tx.Model(&record).Update("status", 0).Error; err != nil {
					return err
				}
				delta = -1
			} else if record.Status == 0 && liked {
				if err := tx.Model(&record).Update("status", 1).Error; err != nil {
					return err
				}
				delta = 1
			}
		}

		if delta != 0 {
			if err := tx.Model(&model.Comment{}).Where("id = ?", id).
				UpdateColumn("like_count", gorm.Expr("GREATEST(like_count + ?, 0)", delta)).Error; err != nil {
				return err
			}
		}

		var fresh model.Comment
		if err := tx.Select("like_count").Where("id = ?", id).First(&fresh).Error; err != nil {
			return err
		}
		resp = CommentLikeResponse{Liked: liked, LikeCount: fresh.LikeCount}
		return nil
	})
	if err != nil {
		var appErr *AppError
		if errors.As(err, &appErr) {
			return CommentLikeResponse{}, appErr
		}
		return CommentLikeResponse{}, ErrInternal.With(err)
	}
	return resp, nil
}
