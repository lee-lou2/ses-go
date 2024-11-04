package handlers

import (
	"ses-go/api/schemas"
	"ses-go/config"
	"ses-go/models"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// InitRecipientsDataHandler 수신자 데이터 초기화 핸들러
func InitRecipientsDataHandler(c fiber.Ctx) error {
	// 기본 컬럼은 email
	columns := []string{"email"}

	// DB 연결 가져오기
	db := config.GetDB()

	// 템플릿 모델 변수 선언
	var template models.Template

	// URL 파라미터에서 템플릿 ID를 가져와 DB에서 해당 템플릿 조회
	templateId := c.Params("templateId")
	templateIdUint, _ := strconv.Atoi(templateId)
	if err := db.Where("id = ?", templateIdUint).First(&template).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 템플릿에서 파라미터 가져오기
	params := template.GetParams()

	// 기본 컬럼에 템플릿 파라미터 추가
	columns = append(columns, *params...)

	// 수신자 데이터 초기화 (헤더 포함 11행)
	recipients := [][]string{
		columns,
		make([]string, len(columns)),
		make([]string, len(columns)),
		make([]string, len(columns)),
		make([]string, len(columns)),
		make([]string, len(columns)),
		make([]string, len(columns)),
		make([]string, len(columns)),
		make([]string, len(columns)),
		make([]string, len(columns)),
		make([]string, len(columns)),
	}
	return c.JSON(schemas.RespInitRecipientsData{
		Data: recipients,
	})
}

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
	recipient := models.Recipient{
		TemplateId: uint(templateIdUint),
		Data:       body.Data,
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
