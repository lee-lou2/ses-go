package pre_send

import (
	"context"
	"encoding/json"
	"log"
	"ses-go/cmd/post_send"
	"ses-go/cmd/send"
	"ses-go/config"
	"ses-go/models"
	"time"

	"gorm.io/gorm"
)

// Run 메세지 전 처리
func Run() {
	for {
		<-time.After(time.Second * 10)
		db := config.GetDB()

		tx := db.Begin()
		var plan models.Plan

		// Plan 조회 및 상태 업데이트
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Preload("Recipient").Where("(status = 0 AND scheduled_at is null) OR (status = 0 AND scheduled_at <= ?)", time.Now()).First(&plan).Error; err != nil {
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
		var data [][]string
		if err := json.Unmarshal([]byte(plan.Recipient.Data), &data); err != nil {
			updatePlanStatus(db, &plan, 5)
			continue
		}
		columns := data[0]
		for i, row := range data {
			// 첫번째 행은 컬럼명이므로 스킵
			if i == 0 {
				continue
			}
			// 이메일이 없는 행은 스킵
			if row[0] == "" {
				continue
			}
			params := make(map[string]string)
			params["email"] = row[0]
			for j := 1; j < len(row); j++ {
				params[columns[j]] = row[j]
			}
			paramsStr, _ := json.Marshal(params)
			messages = append(messages, models.Message{
				To:     row[0],
				Params: string(paramsStr),
				PlanId: plan.ID,
			})
		}

		// messages bulk create
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
