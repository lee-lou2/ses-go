package pre_send

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"ses-go/cmd/post_send"
	"ses-go/cmd/send"
	"ses-go/config"
	"ses-go/models"
	"ses-go/pkg/google"
	"time"
)

// Run 메세지 전 처리
func Run() {
	for {
		<-time.After(time.Second * 10)
		db := config.GetDB()

		tx := db.Begin()
		var plan models.Plan

		// Plan 조회 및 상태 업데이트
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("(status = 0 AND scheduled_at is null) OR (status = 0 AND scheduled_at <= ?)", time.Now()).First(&plan).Error; err != nil {
			tx.Rollback()
			continue
		}
		plan.Status = 1
		if err := tx.Save(&plan).Error; err != nil {
			tx.Rollback()
			continue
		}
		tx.Commit()

		// 메시지 생성
		var messages []models.Message
		if err := google.GetEmailsFromSheet(plan.SheetId, &messages); err != nil {
			// 메시지 생성 실패 시 상태를 업데이트하기 위해 트랜잭션 사용
			updatePlanStatus(db, &plan, 5)
			continue
		}

		// messages bulk create
		for m := range messages {
			messages[m].PlanId = plan.ID
		}
		for i := 0; i < len(messages); i += 1000 {
			minLimit := len(messages)
			if minLimit > i+1000 {
				minLimit = i + 1000
			}
			// plan id 추가
			if err := db.Create(messages[i:minLimit]).Error; err != nil {
				// 메시지 생성 중 오류 발생 시 상태 업데이트
				updatePlanStatus(db, &plan, 6)
				continue
			}
		}

		// 전송 전 상태로 변경
		updatePlanStatus(db, &plan, 2)

		var template models.Template
		if err := db.First(&template, plan.TemplateId).Error; err != nil {
			continue
		}
		subject := template.Subject
		body := template.Body
		ctx := context.Background()

		// 메시지 처리
		for _, msg := range messages {
			params := make(map[string]string)
			paramsStr := msg.Params
			if err := json.Unmarshal([]byte(paramsStr), &params); err != nil {
				post_send.AddMessage("", int(msg.ID), 2, err.Error())
				continue
			}
			send.AddMessage(int(msg.ID), msg.To, subject, body, params, ctx)
		}

		// 메시지 전송 완료 시 상태 업데이트
		updatePlanStatus(db, &plan, 3)
	}
}

// updatePlanStatus Plan 상태 업데이트
func updatePlanStatus(db *gorm.DB, plan *models.Plan, status int) {
	tx := db.Begin()
	defer tx.Commit()
	plan.Status = status
	if err := tx.Save(plan).Error; err != nil {
		log.Printf("failed to update plan status: %v", err)
	}
}
