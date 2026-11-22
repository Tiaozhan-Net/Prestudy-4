package model

import "gorm.io/gorm"

type ImageURL struct {
	ArticleID int64  `gorm:"type:bigint not null;index;comment:笔记ID" json:"articleId"`
	URL       string `gorm:"type:longtext not null;comment:图片URL（可为 base64 data URL）" json:"url"`

	BaseModel
}

func (ImageURL) TableName() string {
	return "image_url"
}

func (e *ImageURL) BeforeCreate(_ *gorm.DB) error {
	return nil
}

func (e *ImageURL) BeforeUpdate(_ *gorm.DB) error {
	return nil
}

func (e *ImageURL) AfterFind(_ *gorm.DB) error {
	return nil
}
