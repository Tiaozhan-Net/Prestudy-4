package service

import (
	"errors"
	"strings"

	"fe_exam/model"
	"fe_exam/pkg"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleSrv struct {
}

// ---------------- 通用工具 ----------------

// escapeLike 转义 LIKE 通配符
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// isDuplicateKey MySQL 唯一键冲突（错误号 1062）
func isDuplicateKey(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}

// toArticleInfo model.Article -> ArticleInfo（Author/Images/Tags 需已 Preload）
func toArticleInfo(at model.Article) ArticleInfo {
	info := ArticleInfo{
		ID:           at.ID,
		AuthorID:     at.AuthorID,
		AuthorName:   at.Author.Name,
		AuthorAvatar: at.Author.Avatar,
		Title:        at.Title,
		Content:      at.Content,
		Category:     Category(at.Category),
		LikeCount:    at.LikeCount,
		CreateTime:   FormatTime(at.CreatedAt),
		Images:       make([]string, 0, len(at.Images)),
		Tags:         make([]string, 0, len(at.Tags)),
	}
	for _, img := range at.Images {
		info.Images = append(info.Images, img.URL)
	}
	for _, tag := range at.Tags {
		info.Tags = append(info.Tags, tag.Name)
	}
	if len(info.Images) > 0 {
		info.Cover = info.Images[0]
	}
	return info
}

// ---------------- 查询 ----------------

func articleBaseQuery(p ArticleListParams) *gorm.DB {
	query := model.DB.Model(&model.Article{})
	if p.Category != "" && p.Category != Category("all") {
		query = query.Where("article.category = ?", string(p.Category))
	}
	if keyword := strings.TrimSpace(p.Keyword); keyword != "" {
		like := "%" + escapeLike(keyword) + "%"
		query = query.Where("(article.title LIKE ? OR article.content LIKE ? OR article.category LIKE ? OR EXISTS ("+
			"SELECT 1 FROM `user` u WHERE u.id = article.author_id AND u.deleted_at IS NULL AND u.name LIKE ?))",
			like, like, like, like)
	}
	return query
}

func articleOrder(query *gorm.DB, s Sort) *gorm.DB {
	switch s {
	case SortDateAsc:
		return query.Order("article.created_at ASC, article.id ASC")
	case SortLikesDesc:
		return query.Order("article.like_count DESC, article.id DESC")
	default:
		return query.Order("article.created_at DESC, article.id DESC")
	}
}

func (a *ArticleSrv) List(p ArticleListParams) (ArticleListResponse, error) {
	p.PagerForm = p.PagerForm.Normalize()

	var total int64
	if err := articleBaseQuery(p).Count(&total).Error; err != nil {
		return ArticleListResponse{}, ErrInternal.With(err)
	}

	var rows []model.Article
	err := articleOrder(articleBaseQuery(p), p.Sort).
		Preload("Author").
		Preload("Images").
		Preload("Tags").
		Scopes(model.Paginate(p.PagerForm)).
		Find(&rows).Error
	if err != nil {
		return ArticleListResponse{}, ErrInternal.With(err)
	}

	return ArticleListResponse{
		PagerForm: p.PagerForm,
		List:      pkg.Map(rows, toArticleInfo),
		Total:     total,
	}, nil
}

func (a *ArticleSrv) Detail(id int64) (ArticleDetailResponse, error) {
	var row model.Article
	err := model.DB.Where("id = ?", id).
		Preload("Author").
		Preload("Images").
		Preload("Tags").
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ArticleDetailResponse{}, ErrArticleNotFound
	}
	if err != nil {
		return ArticleDetailResponse{}, ErrInternal.With(err)
	}
	return ArticleDetailResponse{ArticleInfo: toArticleInfo(row)}, nil
}

// ---------------- 写入 ----------------

func (a *ArticleSrv) Create(uid int64, p ArticleCreateParams) (ArticleCreateResponse, error) {
	if strings.TrimSpace(p.Title) == "" || strings.TrimSpace(p.Content) == "" || p.Category == "" {
		return ArticleCreateResponse{}, ErrParam.WithMsg("标题、分类和正文不能为空")
	}
	for _, img := range p.Images {
		if len(img) > 2<<20 {
			return ArticleCreateResponse{}, ErrParam.WithMsg("单张图片过大（base64 后需小于 2MB）")
		}
	}

	var createdID int64
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		tags := make([]model.Tag, 0, len(p.Tags))
		seen := make(map[string]struct{}, len(p.Tags))
		for _, name := range p.Tags {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			var tag model.Tag
			if err := tx.Where("name = ?", name).FirstOrCreate(&tag, model.Tag{Name: name}).Error; err != nil {
				return err
			}
			tags = append(tags, tag)
		}

		images := make([]model.ImageURL, 0, len(p.Images))
		for _, url := range p.Images {
			url = strings.TrimSpace(url)
			if url == "" {
				continue
			}
			images = append(images, model.ImageURL{URL: url})
		}

		article := model.Article{
			AuthorID: uid,
			Title:    strings.TrimSpace(p.Title),
			Content:  strings.TrimSpace(p.Content),
			Category: string(p.Category),
			Tags:     tags,
			Images:   images,
		}
		if err := tx.Create(&article).Error; err != nil {
			return err
		}
		createdID = article.ID
		return nil
	})
	if err != nil {
		return ArticleCreateResponse{}, ErrInternal.With(err)
	}

	detail, err := a.Detail(createdID)
	if err != nil {
		return ArticleCreateResponse{}, err
	}
	return ArticleCreateResponse{ArticleInfo: detail.ArticleInfo}, nil
}

func (a *ArticleSrv) Delete(uid, id int64) error {
	var article model.Article
	err := model.DB.Where("id = ?", id).First(&article).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrArticleNotFound
	}
	if err != nil {
		return ErrInternal.With(err)
	}
	if article.AuthorID != uid {
		return ErrArticleForbidden
	}

	txErr := model.DB.Transaction(func(tx *gorm.DB) error {
		var commentIDs []int64
		if err := tx.Model(&model.Comment{}).Where("article_id = ?", id).Pluck("id", &commentIDs).Error; err != nil {
			return err
		}
		if len(commentIDs) > 0 {
			if err := tx.Unscoped().Where("target_type = ? AND target_id IN ?", LikeTargetComment, commentIDs).
				Delete(&model.LikeRecord{}).Error; err != nil {
				return err
			}
			if err := tx.Where("id IN ?", commentIDs).Delete(&model.Comment{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Unscoped().Where("target_type = ? AND target_id = ?", LikeTargetArticle, id).
			Delete(&model.LikeRecord{}).Error; err != nil {
			return err
		}
		if err := tx.Where("article_id = ?", id).Delete(&model.ImageURL{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&article).Association("Tags").Clear(); err != nil {
			return err
		}
		return tx.Delete(&model.Article{}, id).Error
	})
	if txErr != nil {
		return ErrInternal.With(txErr)
	}
	return nil
}

// ---------------- 点赞 ----------------

func (a *ArticleSrv) Like(uid, id int64) (ArticleLikeResponse, error) {
	return a.setLike(uid, id, true)
}

func (a *ArticleSrv) Unlike(uid, id int64) (ArticleLikeResponse, error) {
	return a.setLike(uid, id, false)
}

func (a *ArticleSrv) setLike(uid, id int64, liked bool) (ArticleLikeResponse, error) {
	var resp ArticleLikeResponse
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		var article model.Article
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&article).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrArticleNotFound
		}
		if err != nil {
			return err
		}

		var record model.LikeRecord
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND target_id = ? AND target_type = ?", uid, id, LikeTargetArticle).
			First(&record).Error

		delta := 0
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			if liked {
				record = model.LikeRecord{UserID: uid, TargetID: id, TargetType: LikeTargetArticle, Status: 1}
				if err := tx.Create(&record).Error; err != nil {
					if isDuplicateKey(err) { // 并发下已被插入，视为已点赞
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
			if err := tx.Model(&model.Article{}).Where("id = ?", id).
				UpdateColumn("like_count", gorm.Expr("GREATEST(like_count + ?, 0)", delta)).Error; err != nil {
				return err
			}
		}

		var fresh model.Article
		if err := tx.Select("like_count").Where("id = ?", id).First(&fresh).Error; err != nil {
			return err
		}
		resp = ArticleLikeResponse{Liked: liked, LikeCount: fresh.LikeCount}
		return nil
	})
	if err != nil {
		var appErr *AppError
		if errors.As(err, &appErr) {
			return ArticleLikeResponse{}, appErr
		}
		return ArticleLikeResponse{}, ErrInternal.With(err)
	}
	return resp, nil
}
