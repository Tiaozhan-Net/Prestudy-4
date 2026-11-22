package model

import "gorm.io/gorm"

type User struct {
	Name        string `gorm:"type:varchar(20) not null;uniqueIndex:uk_user_name;comment:姓名" json:"name"`
	Password    string `gorm:"type:varchar(255);comment:密码" json:"password"`
	Avatar      string `gorm:"type:varchar(512);comment:头像URL" json:"avatar"`
	Slogan      string `gorm:"type:varchar(32);comment:个性签名" json:"slogan"`
	Description string `gorm:"type:varchar(120);comment:个人简介" json:"description"`
	Location    string `gorm:"type:varchar(40) not null;comment:所在地" json:"location"`

	BaseModel
	Interests   []Tag        `gorm:"many2many:user_interests" json:"interests,omitempty"`
	Followings  []User       `gorm:"many2many:user_follows;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:FollowID" json:"followings,omitempty"`
	LikeRecords []LikeRecord `gorm:"foreignKey:UserID" json:"likeRecords,omitempty"`
}

func (User) TableName() string {
	return "user"
}

func (e *User) BeforeCreate(_ *gorm.DB) error {
	return nil
}

func (e *User) BeforeUpdate(_ *gorm.DB) error {
	return nil
}

func (e *User) AfterFind(_ *gorm.DB) error {
	return nil
}
