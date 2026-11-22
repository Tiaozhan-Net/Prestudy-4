package model

import "gorm.io/gorm"

type LikeRecord struct {
	UserID     int64  `gorm:"type:bigint not null;uniqueIndex:uk_user_target;comment:用户ID" json:"userId"`
	TargetID   int64  `gorm:"type:bigint not null;uniqueIndex:uk_user_target;comment:被点赞对象ID" json:"targetId"`
	TargetType string `gorm:"type:varchar(20) not null;uniqueIndex:uk_user_target;comment:对象类型:article/comment" json:"targetType"`
	Status     int8   `gorm:"type:tinyint not null;default:1;comment:1=点赞 0=取消" json:"status"`

	BaseModel
}

func (LikeRecord) TableName() string {
	return "like_record"
}

func (e *LikeRecord) BeforeCreate(_ *gorm.DB) error {
	return nil
}

func (e *LikeRecord) BeforeUpdate(_ *gorm.DB) error {
	return nil
}

func (e *LikeRecord) AfterFind(_ *gorm.DB) error {
	return nil
}
