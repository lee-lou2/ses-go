package models

import (
	"gorm.io/gorm"
	"ses-go/config"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique;not null;type:varchar(255)"`
	GoogleId string `json:"google_id" gorm:"unique;not null;type:varchar(100)"`
}

func (u *User) TableName() string {
	return "users"
}

type UserSession struct {
	gorm.Model
	UserId uint   `json:"user_id" gorm:"index;not null"`
	User   User   `json:"user" gorm:"foreignKey:UserId;references:ID"`
	UUID   string `json:"uuid" gorm:"unique;not null;type:varchar(100)"`
}

func (us *UserSession) TableName() string {
	return "user_sessions"
}

func init() {
	db := config.GetDB()
	_ = db.AutoMigrate(&User{})
	_ = db.AutoMigrate(&UserSession{})
}
