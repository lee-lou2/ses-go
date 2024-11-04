package models

import (
	"ses-go/config"
	"time"

	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null;type:varchar(255)"`
	TemplateId  uint      `json:"template_id" gorm:"index;not null"`
	Template    Template  `json:"template" gorm:"foreignKey:TemplateId;references:ID"`
	RecipientId uint      `json:"recipient_id" gorm:"index;not null"`
	Recipient   Recipient `json:"recipient" gorm:"foreignKey:RecipientId;references:ID"`
	// Status
	// 0: 생성 완료, 1: 준비 완료, 2: 메세지 생성 완료, 3: 메세지 전달,
	// 4: 메세지 전송 완료, 5: 실패(이메일 조회 실패), 6: 실패(메세지 생성 실패),
	// 7: 실패(메세지 전송 실패), 8: 실패(기타)
	Status      int        `json:"status" gorm:"default:0;not null;type:tinyint"`
	ScheduledAt *time.Time `json:"scheduled_at" gorm:"index;null;type:datetime;default:null"`
}

func (p *Plan) TableName() string {
	return "plans"
}

func init() {
	db := config.GetDB()
	_ = db.AutoMigrate(&Plan{})
}
