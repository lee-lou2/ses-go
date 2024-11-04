package handlers

import (
	"fmt"
	"ses-go/api/schemas"
	"ses-go/config"
	"ses-go/models"
	"time"

	"github.com/gofiber/fiber/v3"
)

// CreatePlanHandler 플랜 생성 핸들러
func CreatePlanHandler(c fiber.Ctx) error {
	body := new(schemas.ReqCreatePlan)
	if err := c.Bind().JSON(&body); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	db := config.GetDB()
	plan := models.Plan{
		Title:       body.Title,
		TemplateId:  body.TemplateId,
		RecipientId: body.RecipientId,
		ScheduledAt: nil,
	}
	// 스케줄링 시간이 있으면 파싱
	if body.ScheduledAt != "" {
		// "2024-09-23T11:23" 형식으로 파싱
		scheduledAt, err := time.Parse("2006-01-02T15:04", body.ScheduledAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		plan.ScheduledAt = &scheduledAt
	}
	// 데이터 생성
	if err := db.Create(&plan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(schemas.RespCreatePlan{
		Id: plan.ID,
	})
}
