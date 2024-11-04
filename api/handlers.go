package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"ses-go/config"
	"ses-go/models"
	"ses-go/pkg/google"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// createPlanHandler 플랜 생성 핸들러
func createPlanHandler(c fiber.Ctx) error {
	body := new(reqCreatePlan)
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
	return c.JSON(respCreatePlan{
		Id: plan.ID,
	})
}

// googleAuthHandler 구글 인증 핸들러
func googleAuthHandler(c fiber.Ctx) error {
	path := google.ConfigGoogle()
	url := path.AuthCodeURL("state")
	return c.Redirect().To(url)
}

// googleCallbackHandler 구글 콜백 핸들러
func googleCallbackHandler(c fiber.Ctx) error {
	token, err := google.ConfigGoogle().Exchange(c.Context(), c.FormValue("code"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 유저 정보 조회
	userInfo, err := google.GetUserInfo(token.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 유저 정보가 없으면 생성
	db := config.GetDB()
	var user models.User
	if err := db.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
		user = models.User{
			Email:    userInfo.Email,
			GoogleId: userInfo.ID,
		}
		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// 세션 생성
	uuidStr := uuid.New().String()
	session := models.UserSession{
		UserId: user.ID,
		UUID:   uuidStr,
	}
	if err := db.Create(&session).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// 세션 쿠키 설정
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    uuidStr,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})
	return c.Redirect().To("/")
}

// logoutHandler 로그아웃 핸들러
func logoutHandler(c fiber.Ctx) error {
	db := config.GetDB()
	session := fiber.Locals[models.UserSession](c, "session")
	db.Delete(&session)
	c.ClearCookie("session")
	return c.Redirect().To("/")
}

// loginTemplateHandler 로그인 템플릿 핸들러
func loginTemplateHandler(c fiber.Ctx) error {
	return c.Render("accounts/login", nil, "layouts/main")
}

// indexTemplateHandler 인덱스 템플릿 핸들러
func indexTemplateHandler(c fiber.Ctx) error {
	db := config.GetDB()
	var plans []models.Plan
	db.Find(&plans).Order("id desc")
	return c.Render("index", fiber.Map{
		"Plans": plans,
	}, "layouts/main")
}

// planCreateTemplateHandler 플랜 상세 템플릿 핸들러
func planCreateTemplateHandler(c fiber.Ctx) error {
	db := config.GetDB()
	var templates []models.Template
	db.Find(&templates).Order("id desc")
	return c.Render("plan/create", fiber.Map{
		"Templates": templates,
	}, "layouts/main")
}

// planDetailTemplateHandler 플랜 상세 템플릿 핸들러
func planDetailTemplateHandler(c fiber.Ctx) error {
	db := config.GetDB()
	var plan models.Plan
	planId := c.Params("planId")
	planIdUint, _ := strconv.Atoi(planId)
	if err := db.Preload("Template").Preload("Recipient").First(&plan, planIdUint).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Render("plan/detail", fiber.Map{
		"Plan": plan,
	}, "layouts/main")
}

// planResultTemplateHandler 플랜 결과 템플릿 핸들러
func planResultTemplateHandler(c fiber.Ctx) error {
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
	return c.Render("plan/result", fiber.Map{
		"Messages": messagesWithResults,
	}, "layouts/main")
}

// templateDetailTemplateHandler 템플릿 상세 템플릿 핸들러
func templateDetailTemplateHandler(c fiber.Ctx) error {
	db := config.GetDB()
	var template models.Template
	templateId := c.Params("templateId")
	templateIdUint, _ := strconv.Atoi(templateId)
	db.First(&template, templateIdUint)
	tinymceApiKey := config.GetEnv("TINYMCE_API_KEY")
	return c.Render("plan/template", fiber.Map{
		"Template":      template,
		"TinymceApiKey": tinymceApiKey,
	}, "layouts/main")
}

// createTemplateHandler 템플릿 생성 핸들러
func createTemplateHandler(c fiber.Ctx) error {
	body := new(reqCreateTemplate)
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
	return c.JSON(respCreateTemplate{
		Id: template.ID,
	})
}

// updateTemplateHandler 템플릿 업데이트 핸들러
func updateTemplateHandler(c fiber.Ctx) error {
	body := new(reqUpdateTemplate)
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
	return c.JSON(respUpdateTemplate{
		Id: template.ID,
	})
}

// addOpenEventHandler 오픈 이벤트 핸들러
func addOpenEventHandler(c fiber.Ctx) error {
	msgId := c.Query("message_id")
	if msgId != "" {
		// 데이터 생성
		db := config.GetDB()
		var message models.MessageResult
		msgIdInt, _ := strconv.Atoi(msgId)
		message.MessageId = uint(msgIdInt)
		message.Status = "Open"
		_ = db.Create(&message).Error
	}
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 0, G: 0, B: 0, A: 0})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	c.Set("Content-Type", "image/png")
	return c.Send(buf.Bytes())
}

// addSendEventHandler 전송 이벤트 핸들러
func addSendEventHandler(c fiber.Ctx) error {
	body := new(reqAddSendEvent)
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if body.Type == "SubscriptionConfirmation" {
		log.Println(body.SubscribeURL)
		return c.JSON(fiber.Map{})
	} else if body.Type != "Notification" {
		return c.JSON(fiber.Map{})
	}

	// 메세지 조회
	var bodyMessage struct {
		EventType string `json:"eventType"`
		Mail      struct {
			MessageId string `json:"messageId"`
		} `json:"mail"`
	}
	if err := json.Unmarshal([]byte(body.Message), &bodyMessage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 메세지 조회
	db := config.GetDB()
	var message models.Message
	if err := db.Where("message_id = ?", bodyMessage.Mail.MessageId).First(&message).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 결과 저장
	result := models.MessageResult{
		MessageId: message.ID,
		Status:    bodyMessage.EventType,
		Raw:       body.Message,
	}
	_ = db.Create(&result).Error
	return c.JSON(fiber.Map{})
}

// initRecipientsDataHandler 수신자 데이터 초기화 핸들러
func initRecipientsDataHandler(c fiber.Ctx) error {
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
	return c.JSON(respInitRecipientsData{
		Data: recipients,
	})
}

// getRecipientsTemplateHandler 발송 대상 상세 템플릿 핸들러
func getRecipientsTemplateHandler(c fiber.Ctx) error {
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
	return c.Render("plan/recipients", fiber.Map{
		"Recipients": data,
		"TemplateId": templateIdUint,
	}, "layouts/main")
}

// createRecipientsHandler 수신자 생성 핸들러
func createRecipientsHandler(c fiber.Ctx) error {
	// URL 파라미터에서 템플릿 ID를 가져옴
	templateId := c.Params("templateId")
	// 요청 바디를 파싱
	body := new(reqCreateRecipients)
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
