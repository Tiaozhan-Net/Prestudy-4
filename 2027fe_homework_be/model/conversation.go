package model

import "gorm.io/gorm"

type Conversation struct {
	UserLowID  int64 `gorm:"type:bigint not null;uniqueIndex:uk_pair;comment:较小用户ID" json:"userLowId"`
	UserHighID int64 `gorm:"type:bigint not null;uniqueIndex:uk_pair;comment:较大用户ID" json:"userHighId"`

	BaseModel
	Messages []Message `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
	UserLow  User      `gorm:"foreignKey:UserLowID" json:"userLow"`
	UserHigh User      `gorm:"foreignKey:UserHighID" json:"userHigh"`
}

func (Conversation) TableName() string {
	return "conversation"
}

func (e *Conversation) BeforeCreate(_ *gorm.DB) error {
	return nil
}

func (e *Conversation) BeforeUpdate(_ *gorm.DB) error {
	return nil
}

func (e *Conversation) AfterFind(_ *gorm.DB) error {
	return nil
}
