package model

import "gorm.io/gorm"

type Article struct {
	AuthorID  int64  `gorm:"type:bigint not null;index;comment:作者ID" json:"authorId"`
	Title     string `gorm:"type:varchar(128) not null;comment:标题" json:"title"`
	Content   string `gorm:"type:text not null;comment:内容" json:"content"`
	Category  string `gorm:"type:varchar(128) not null;comment:类型" json:"category"`
	LikeCount int64  `gorm:"type:bigint not null;default:0;comment:点赞数量" json:"likeCount"`

	BaseModel
	Author User       `gorm:"foreignKey:AuthorID" json:"author"`
	Images []ImageURL `gorm:"foreignKey:ArticleID" json:"images,omitempty"`
	Tags   []Tag      `gorm:"many2many:article_tags" json:"tags,omitempty"`
}

func (Article) TableName() string {
	return "article"
}

func (e *Article) BeforeCreate(_ *gorm.DB) error {
	return nil
}

func (e *Article) BeforeUpdate(_ *gorm.DB) error {
	return nil
}

func (e *Article) AfterFind(_ *gorm.DB) error {
	return nil
}
