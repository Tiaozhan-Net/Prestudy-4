package model

import "gorm.io/gorm"

type Tag struct {
	Name string `gorm:"type:varchar(40) not null;uniqueIndex:uk_tag_name;comment:标签名" json:"name"`

	BaseModel
}

func (Tag) TableName() string {
	return "tag"
}

func (e *Tag) BeforeCreate(_ *gorm.DB) error {
	return nil
}

func (e *Tag) BeforeUpdate(_ *gorm.DB) error {
	return nil
}

func (e *Tag) AfterFind(_ *gorm.DB) error {
	return nil
}
