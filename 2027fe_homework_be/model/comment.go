package model

import "gorm.io/gorm"

type Comment struct {
	ArticleID int64  `gorm:"type:bigint not null;index;comment:笔记ID" json:"articleId"`
	UserID    int64  `gorm:"type:bigint not null;index;comment:用户ID" json:"userId"`
	Content   string `gorm:"type:text not null;comment:内容" json:"content"`
	LikeCount int64  `gorm:"type:bigint not null;default:0;comment:点赞数量" json:"likeCount"`

	BaseModel
	Article Article `gorm:"foreignKey:ArticleID" json:"article"`
	User    User    `gorm:"foreignKey:UserID" json:"user"`
}

func (Comment) TableName() string {
	return "comment"
}

func (e *Comment) BeforeCreate(_ *gorm.DB) error {
	return nil
}

func (e *Comment) BeforeUpdate(_ *gorm.DB) error {
	return nil
}

func (e *Comment) AfterFind(_ *gorm.DB) error {
	return nil
}
