package handlers

import (
	"encoding/json"
	"ses-go/config"
	"ses-go/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

// LoginHTMLRenderHandler 로그인 템플릿 핸들러
func LoginHTMLRenderHandler(c fiber.Ctx) error {
	return c.Render("accounts/login", nil, "layouts/main")
}

// IndexHTMLRenderHandler 인덱스 템플릿 핸들러
func IndexHTMLRenderHandler(c fiber.Ctx) error {
	db := config.GetDB()
	var plans []models.Plan
	db.Find(&plans).Order("id desc")
	return c.Render("index", fiber.Map{
		"Plans": plans,
	}, "layouts/main")
}

// PlanCreateHTMLRenderHandler 플랜 상세 템플릿 핸들러
func PlanCreateHTMLRenderHandler(c fiber.Ctx) error {
	db := config.GetDB()
	var templates []models.Template
	db.Find(&templates).Order("id desc")
	return c.Render("plans/create", fiber.Map{
		"Templates": templates,
	}, "layouts/main")
}

// PlanDetailHTMLRenderHandler 플랜 상세 템플릿 핸들러
func PlanDetailHTMLRenderHandler(c fiber.Ctx) error {
	db := config.GetDB()
	var plan models.Plan
	planId := c.Params("planId")
	planIdUint, _ := strconv.Atoi(planId)
	if err := db.Preload("Template").Preload("Recipient").First(&plan, planIdUint).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Render("plans/detail", fiber.Map{
		"Plan": plan,
	}, "layouts/main")
}

// PlanResultHTMLRenderHandler 플랜 결과 템플릿 핸들러
func PlanResultHTMLRenderHandler(c fiber.Ctx) error {
	db := config.GetDB()
	planId := c.Params("planId")
	planIdUint, _ := strconv.Atoi(planId)
	var messagesWithResults []struct {
		ID        uint      `json:"id"`
		To        string    `json:"to"`
		Params    string    `json:"params"`
		Status    int       `json:"status"`
		Results   string    `json:"results"`
		CreatedAt time.Time `json:"created_at"`
	}
	err := db.Table("messages").
		Select(`messages.id, messages."to", messages.params, messages.status, messages.created_at, 
                GROUP_CONCAT(CONCAT(message_results.status, '(', strftime('%Y-%m-%d %H:%M', message_results.created_at), ')')) as results`).
		Joins("LEFT JOIN message_results ON messages.id = message_results.message_id").
		Where("messages.plan_id = ?", planIdUint).
		Group("messages.id").
		Scan(&messagesWithResults).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Render("plans/result", fiber.Map{
		"Messages": messagesWithResults,
	}, "layouts/main")
}

// TemplateDetailHTMLRenderHandler 템플릿 상세 템플릿 핸들러
func TemplateDetailHTMLRenderHandler(c fiber.Ctx) error {
	db := config.GetDB()
	var template models.Template
	templateId := c.Params("templateId")
	templateIdUint, _ := strconv.Atoi(templateId)
	db.First(&template, templateIdUint)
	tinymceApiKey := config.GetEnv("TINYMCE_API_KEY")
	return c.Render("plans/template", fiber.Map{
		"Template":      template,
		"TinymceApiKey": tinymceApiKey,
	}, "layouts/main")
}

// GetRecipientsHTMLRenderHandler 발송 대상 상세 템플릿 핸들러
func GetRecipientsHTMLRenderHandler(c fiber.Ctx) error {
	// DB 연결 가져오기
	db := config.GetDB()
	// 수신자 모델 변수 선언
	var recipients models.Recipient
	// URL 파라미터에서 수신자 ID와 템플릿 ID 가져오기
	recipientId := c.Params("recipientId")
	recipientIdUint, _ := strconv.Atoi(recipientId)
	templateId := c.Params("templateId")
	templateIdUint, _ := strconv.Atoi(templateId)
	query := db.Where("id = ?", recipientIdUint).Where("template_id = ?", templateIdUint)
	// 생성자이거나 뷰어인 경우 조회 가능
	user := fiber.Locals[models.User](c, "user")
	query = query.Where("creator_id = ? OR id IN (SELECT recipient_id FROM recipient_viewers WHERE user_id = ?)", user.ID, user.ID)
	if err := query.First(&recipients).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// data 를 json 으로 변환해서 리스트로 반환
	var data [][]string
	// JSON 문자열을 2차원 문자열 배열로 변환
	if err := json.Unmarshal([]byte(recipients.Data), &data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// 템플릿 렌더링하여 응답 반환
	return c.Render("plans/recipients", fiber.Map{
		"Recipients": data,
		"TemplateId": templateIdUint,
	}, "layouts/main")
}

// TokenHTMLRenderHandler 토큰 템플릿 핸들러
func TokenHTMLRenderHandler(c fiber.Ctx) error {
	user := fiber.Locals[models.User](c, "user")
	var tokens []models.UserToken
	db := config.GetDB()
	db.Where("user_id = ?", user.ID).Find(&tokens)
	return c.Render("tokens/index", fiber.Map{
		"Tokens": tokens,
	}, "layouts/main")
}
