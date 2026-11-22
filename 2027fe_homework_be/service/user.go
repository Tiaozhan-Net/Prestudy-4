package service

import (
	"errors"
	"strings"

	"fe_exam/model"
	"fe_exam/pkg"

	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
)

type UserSrv struct {
}

func ensureUserExists(id int64) error {
	var count int64
	if err := model.DB.Model(&model.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return ErrInternal.With(err)
	}
	if count == 0 {
		return ErrUserNotFound
	}
	return nil
}

// info 组装本人资料（含关注/点赞聚合）
func (u *UserSrv) info(uid int64) (UserInfo, error) {
	var user model.User
	err := model.DB.Preload("Interests").Where("id = ?", uid).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserInfo{}, ErrUserNotFound
	}
	if err != nil {
		return UserInfo{}, ErrInternal.With(err)
	}

	following := make([]int64, 0)
	if err := model.DB.Table("user_follows").
		Where("user_id = ?", uid).
		Order("follow_id ASC").
		Pluck("follow_id", &following).Error; err != nil {
		return UserInfo{}, ErrInternal.With(err)
	}

	likeArticles := make([]int64, 0)
	if err := model.DB.Model(&model.LikeRecord{}).
		Where("user_id = ? AND target_type = ? AND status = 1", uid, LikeTargetArticle).
		Order("target_id ASC").
		Pluck("target_id", &likeArticles).Error; err != nil {
		return UserInfo{}, ErrInternal.With(err)
	}

	likeComments := make([]int64, 0)
	if err := model.DB.Model(&model.LikeRecord{}).
		Where("user_id = ? AND target_type = ? AND status = 1", uid, LikeTargetComment).
		Order("target_id ASC").
		Pluck("target_id", &likeComments).Error; err != nil {
		return UserInfo{}, ErrInternal.With(err)
	}

	return UserInfo{
		ID:           user.ID,
		Name:         user.Name,
		Avatar:       user.Avatar,
		CreateTime:   FormatTime(user.CreatedAt),
		Slogan:       user.Slogan,
		Description:  user.Description,
		Location:     user.Location,
		Interests:    pkg.Map(user.Interests, func(t model.Tag) string { return t.Name }),
		Following:    following,
		LikeArticles: likeArticles,
		LikeComments: likeComments,
	}, nil
}

func (u *UserSrv) Me(uid int64) (UserMeResponse, error) {
	info, err := u.info(uid)
	if err != nil {
		return UserMeResponse{}, err
	}
	return UserMeResponse{UserInfo: info}, nil
}

// UpdateMe 部分更新；第二返回值表示是否修改了密码
func (u *UserSrv) UpdateMe(uid int64, p UserUpdateParams) (UserUpdateResponse, bool, error) {
	var user model.User
	err := model.DB.Where("id = ?", uid).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserUpdateResponse{}, false, ErrUserNotFound
	}
	if err != nil {
		return UserUpdateResponse{}, false, ErrInternal.With(err)
	}

	updates := make(map[string]any)
	passwordChanged := false

	if p.Name != nil {
		name := strings.TrimSpace(*p.Name)
		if name == "" {
			return UserUpdateResponse{}, false, ErrParam.WithMsg("用户名不能为空")
		}
		if name != user.Name {
			var count int64
			if err := model.DB.Model(&model.User{}).Where("name = ? AND id <> ?", name, uid).Count(&count).Error; err != nil {
				return UserUpdateResponse{}, false, ErrInternal.With(err)
			}
			if count > 0 {
				return UserUpdateResponse{}, false, ErrUserNameExists
			}
			updates["name"] = name
		}
	}
	if p.Avatar != nil {
		updates["avatar"] = *p.Avatar
	}
	if p.Slogan != nil {
		updates["slogan"] = *p.Slogan
	}
	if p.Description != nil {
		updates["description"] = *p.Description
	}
	if p.Location != nil {
		updates["location"] = *p.Location
	}
	if p.Password != nil && *p.Password != "" {
		hash, err := argon2id.CreateHash(*p.Password, argon2id.DefaultParams)
		if err != nil {
			return UserUpdateResponse{}, false, ErrInternal.With(err)
		}
		updates["password"] = hash
		passwordChanged = true
	}

	if len(updates) > 0 {
		if err := model.DB.Model(&model.User{}).Where("id = ?", uid).Updates(updates).Error; err != nil {
			if isDuplicateKey(err) {
				return UserUpdateResponse{}, false, ErrUserNameExists
			}
			return UserUpdateResponse{}, false, ErrInternal.With(err)
		}
	}

	if p.Interests != nil {
		tags := make([]model.Tag, 0, len(p.Interests))
		seen := make(map[string]struct{}, len(p.Interests))
		for _, name := range p.Interests {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			var tag model.Tag
			if err := model.DB.Where("name = ?", name).FirstOrCreate(&tag, model.Tag{Name: name}).Error; err != nil {
				return UserUpdateResponse{}, false, ErrInternal.With(err)
			}
			tags = append(tags, tag)
		}
		if err := model.DB.Model(&user).Association("Interests").Replace(tags); err != nil {
			return UserUpdateResponse{}, false, ErrInternal.With(err)
		}
	}

	if passwordChanged {
		// 密码已修改时，由 controller 清理当前会话
	}

	info, err := u.info(uid)
	if err != nil {
		return UserUpdateResponse{}, passwordChanged, err
	}
	return UserUpdateResponse{UserInfo: info}, passwordChanged, nil
}

func (u *UserSrv) Public(uid, viewer int64) (UserPublicResponse, error) {
	var user model.User
	err := model.DB.Preload("Interests").Where("id = ?", uid).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserPublicResponse{}, ErrUserNotFound
	}
	if err != nil {
		return UserPublicResponse{}, ErrInternal.With(err)
	}

	isFollowing := false
	if viewer > 0 && viewer != uid {
		var count int64
		if err := model.DB.Table("user_follows").
			Where("user_id = ? AND follow_id = ?", viewer, uid).
			Count(&count).Error; err != nil {
			return UserPublicResponse{}, ErrInternal.With(err)
		}
		isFollowing = count > 0
	}

	return UserPublicResponse{
		ID:          user.ID,
		Name:        user.Name,
		Avatar:      user.Avatar,
		CreateTime:  FormatTime(user.CreatedAt),
		Slogan:      user.Slogan,
		Description: user.Description,
		Location:    user.Location,
		Interests:   pkg.Map(user.Interests, func(t model.Tag) string { return t.Name }),
		IsFollowing: isFollowing,
	}, nil
}

func (u *UserSrv) Follow(uid, target int64) (FollowResponse, error) {
	if uid == target {
		return FollowResponse{}, ErrFollowSelf
	}
	if err := ensureUserExists(target); err != nil {
		return FollowResponse{}, err
	}
	err := model.DB.Exec("INSERT IGNORE INTO user_follows (user_id, follow_id) VALUES (?, ?)", uid, target).Error
	if err != nil {
		return FollowResponse{}, ErrInternal.With(err)
	}
	return FollowResponse{Following: true}, nil
}

func (u *UserSrv) Unfollow(uid, target int64) (FollowResponse, error) {
	if err := ensureUserExists(target); err != nil {
		return FollowResponse{}, err
	}
	err := model.DB.Exec("DELETE FROM user_follows WHERE user_id = ? AND follow_id = ?", uid, target).Error
	if err != nil {
		return FollowResponse{}, ErrInternal.With(err)
	}
	return FollowResponse{Following: false}, nil
}

// Articles 某用户的文章列表：{list,total,page,size}
func (u *UserSrv) Articles(target int64, p UserArticlesParams) (UserArticlesResponse, error) {
	if err := ensureUserExists(target); err != nil {
		return UserArticlesResponse{}, err
	}
	p.PagerForm = p.PagerForm.Normalize()

	base := func() *gorm.DB {
		return model.DB.Model(&model.Article{}).Where("article.author_id = ?", target)
	}

	var total int64
	if err := base().Count(&total).Error; err != nil {
		return UserArticlesResponse{}, ErrInternal.With(err)
	}

	var rows []model.Article
	err := articleOrder(base(), p.Sort).
		Preload("Author").
		Preload("Images").
		Preload("Tags").
		Scopes(model.Paginate(p.PagerForm)).
		Find(&rows).Error
	if err != nil {
		return UserArticlesResponse{}, ErrInternal.With(err)
	}

	return UserArticlesResponse{
		PagerForm: p.PagerForm,
		List:      pkg.Map(rows, toArticleInfo),
		Total:     total,
	}, nil
}
