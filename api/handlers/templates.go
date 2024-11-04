package handlers

import (
	"ses-go/api/schemas"
	"ses-go/config"
	"ses-go/models"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// CreateTemplateHandler 템플릿 생성 핸들러
func CreateTemplateHandler(c fiber.Ctx) error {
	body := new(schemas.ReqCreateTemplate)
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	db := config.GetDB()
	template := models.Template{
		Subject: body.Subject,
		Body:    "",
	}
	// 데이터 생성
	if err := db.Create(&template).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(schemas.RespCreateTemplate{
		Id: template.ID,
	})
}

// UpdateTemplateHandler 템플릿 업데이트 핸들러
func UpdateTemplateHandler(c fiber.Ctx) error {
	body := new(schemas.ReqUpdateTemplate)
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	db := config.GetDB()
	var template models.Template
	templateId := c.Params("templateId")
	templateIdUint, _ := strconv.Atoi(templateId)
	if err := db.First(&template, templateIdUint).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	template.Body = body.Body
	if err := db.Save(&template).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(schemas.RespUpdateTemplate{
		Id: template.ID,
	})
}
