package handlers

import (
	"encoding/json"
	"ses-go/api/schemas"
	"ses-go/config"
	"ses-go/models"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// CreateRecipientHandler 수신자 생성 핸들러
func CreateRecipientHandler(c fiber.Ctx) error {
	// URL 파라미터에서 템플릿 ID를 가져옴
	templateId := c.Params("templateId")
	// 요청 바디를 파싱
	body := new(schemas.ReqCreateRecipients)
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// DB 연결 가져오기
	db := config.GetDB()
	// 문자열 템플릿 ID를 uint로 변환
	templateIdUint, _ := strconv.Atoi(templateId)
	// 수신자 모델 생성
	user := fiber.Locals[models.User](c, "user")
	data := body.Data
	// email 이 없는 경우 제외
	recipients := [][]string{}
	for _, row := range data {
		if row[0] == "" {
			continue
		}
		recipients = append(recipients, row)
	}
	jsonData, _ := json.Marshal(recipients)
	recipient := models.Recipient{
		TemplateId: uint(templateIdUint),
		Data:       string(jsonData),
		CreatorId:  user.ID,
	}
	// DB에 수신자 데이터 저장
	err := db.Create(&recipient).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// 생성된 수신자의 ID를 응답으로 반환
	return c.JSON(fiber.Map{
		"id": recipient.ID,
	})
}
