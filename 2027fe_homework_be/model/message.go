package model

import "gorm.io/gorm"

type Message struct {
	ConversationID int64  `gorm:"type:bigint not null;index;comment:会话ID" json:"conversationId"`
	SenderID       int64  `gorm:"type:bigint not null;index;comment:发送者ID" json:"senderId"`
	ReceiverID     int64  `gorm:"type:bigint not null;index;comment:接收者ID" json:"receiverId"`
	Content        string `gorm:"type:text not null;comment:内容" json:"content"`

	BaseModel
	Sender User `gorm:"foreignKey:SenderID" json:"sender"`
}

func (Message) TableName() string {
	return "message"
}

func (e *Message) BeforeCreate(_ *gorm.DB) error {
	return nil
}

func (e *Message) BeforeUpdate(_ *gorm.DB) error {
	return nil
}

func (e *Message) AfterFind(_ *gorm.DB) error {
	return nil
}
