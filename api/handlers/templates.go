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
	user := c.Locals("user").(*models.User)
	template := models.Template{
		Subject:   body.Subject,
		Body:      "",
		CreatorId: user.ID,
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
	user := c.Locals("user").(*models.User)
	var template models.Template
	templateId := c.Params("templateId")
	templateIdUint, _ := strconv.Atoi(templateId)
	if err := db.First(&template, templateIdUint).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if template.CreatorId != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "권한이 없습니다."})
	}
	template.Body = body.Body
	if err := db.Save(&template).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(schemas.RespUpdateTemplate{
		Id: template.ID,
	})
}

// GetTemplateFieldsHandler 템플릿 필드 조회 핸들러
func GetTemplateFieldsHandler(c fiber.Ctx) error {
	db := config.GetDB()
	templateId := c.Params("templateId")
	templateIdUint, _ := strconv.Atoi(templateId)
	var template models.Template
	if err := db.First(&template, templateIdUint).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	fields := []string{"email"}
	params := template.GetParams()
	fields = append(fields, *params...)
	return c.JSON(schemas.RespGetTemplateFields{
		Fields: fields,
	})
}
