package models

import (
	"gorm.io/gorm"
	"ses-go/config"
)

type Message struct {
	gorm.Model
	MessageId string `json:"message_id" gorm:"index;not null;type:varchar(255)"`
	PlanId    uint   `json:"plan_id" gorm:"index;not null"`
	Plan      Plan   `json:"plan" gorm:"foreignKey:PlanId;references:ID"`
	To        string `json:"to" gorm:"index;not null;type:varchar(255)"`
	Params    string `json:"params" gorm:"null;type:json"`
	// Status 0: 생성 완료, 1: 전송 완료, 2: 실패, 3: 중지
	Status int    `json:"status" gorm:"default:0;not null;type:tinyint"`
	Error  string `json:"error" gorm:"null;type:varchar(255)"`
}

func (m *Message) TableName() string {
	return "messages"
}

type MessageResult struct {
	gorm.Model
	MessageId uint    `json:"message_id" gorm:"index;not null"`
	Message   Message `json:"message" gorm:"foreignKey:MessageId;references:ID"`
	Status    string  `json:"status" gorm:"not null;type:varchar(50)"`
	Raw       string  `json:"raw" gorm:"null;type:json"`
}

func (m *MessageResult) TableName() string {
	return "message_results"
}

func init() {
	db := config.GetDB()
	_ = db.AutoMigrate(&Message{})
	_ = db.AutoMigrate(&MessageResult{})
}
