package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"image"
	"image/color"
	"image/png"
	"ses-go/config"
	"ses-go/models"
	"ses-go/pkg/google"
	"strconv"
	"time"
)

// createPlanHandler 플랜 생성 핸들러
func createPlanHandler(c fiber.Ctx) error {
	body := new(reqCreatePlan)
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	db := config.GetDB()
	plan := models.Plan{
		Title:       body.Title,
		TemplateId:  body.TemplateId,
		SheetId:     body.SheetId,
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

// createSheetAndShareHandler 시트 생성 및 공유 핸들러
func createSheetAndShareHandler(c fiber.Ctx) error {
	body := new(reqCreateSheetAndShare)
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// 파라미터 추출
	db := config.GetDB()
	var template models.Template
	if err := db.Where("id = ?", body.TemplateId).First(&template).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// 컬럼 생성
	var columns []string
	params := template.GetParams()
	columns = append(columns, "email")
	for _, param := range *params {
		if param != "email" {
			columns = append(columns, param)
		}
	}

	// 시트 생성 및 공유
	email := fiber.Locals[models.User](c, "user").Email
	sheetId, err := google.CreateSheetAndShare(email, &columns)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(respCreateSheetAndShare{
		SheetId: sheetId,
	})
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
	db.First(&plan, c.Params("id"))
	return c.Render("plan/detail", fiber.Map{
		"Plan": plan,
	}, "layouts/main")
}

// planResultTemplateHandler 플랜 결과 템플릿 핸들러
func planResultTemplateHandler(c fiber.Ctx) error {
	db := config.GetDB()
	planId, _ := strconv.Atoi(c.Params("id"))
	var messages []models.Message
	if err := db.Where("plan_id = ?", planId).Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	var results []models.MessageResult
	if len(messages) > 0 {
		messageIDs := make([]uint, len(messages))
		for i, message := range messages {
			messageIDs[i] = message.ID
		}
		if err := db.Where("message_id IN ?", messageIDs).Find(&results).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.Render("plan/result", fiber.Map{
		"Messages": messages,
		"Results":  results,
	}, "layouts/main")
}

// templateDetailTemplateHandler 템플릿 상세 템플릿 핸들러
func templateDetailTemplateHandler(c fiber.Ctx) error {
	db := config.GetDB()
	var template models.Template
	db.First(&template, c.Params("id"))
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
	if err := db.First(&template, c.Params("id")).Error; err != nil {
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
	fmt.Println(c.Body())
	body := new(reqAddSendEvent)
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println(body)
	if body.Type != "Notification" {
		return c.JSON(fiber.Map{})
	}

	// 메세지 조회
	db := config.GetDB()
	var message models.Message
	if err := db.Where("message_id = ?", body.MessageId).First(&message).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 결과 저장
	var result models.MessageResult
	result.MessageId = message.ID
	var detail map[string]interface{}
	if err := json.Unmarshal([]byte(body.Message), &detail); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	result.Status = detail["eventType"].(string)
	result.Raw = body.Message
	_ = db.Create(&result).Error
	return c.JSON(fiber.Map{})
}
