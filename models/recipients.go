package models

import (
	"ses-go/config"

	"gorm.io/gorm"
)

type Recipient struct {
	gorm.Model
	TemplateId uint     `json:"template_id" gorm:"index;not null"`
	Template   Template `json:"template" gorm:"foreignKey:TemplateId;references:ID"`
	Data       string   `json:"data" gorm:"not null;type:text"`
	CreatorId  uint     `json:"creator_id" gorm:"index;not null"`
	Creator    User     `json:"creator" gorm:"foreignKey:CreatorId;references:ID"`
}

func (p *Recipient) TableName() string {
	return "recipients"
}

type RecipientViewer struct {
	gorm.Model
	RecipientId uint      `json:"recipient_id" gorm:"index;not null"`
	Recipient   Recipient `json:"recipient" gorm:"foreignKey:RecipientId;references:ID"`
	UserId      uint      `json:"user_id" gorm:"index;not null"`
	User        User      `json:"user" gorm:"foreignKey:UserId;references:ID"`
}

func (p *RecipientViewer) TableName() string {
	return "recipient_viewers"
}

func init() {
	db := config.GetDB()
	_ = db.AutoMigrate(&Recipient{})
	_ = db.AutoMigrate(&RecipientViewer{})
}
