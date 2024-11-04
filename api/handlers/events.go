package handlers

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"log"
	"ses-go/api/schemas"
	"ses-go/config"
	"ses-go/models"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// AddOpenEventHandler 오픈 이벤트 핸들러
func AddOpenEventHandler(c fiber.Ctx) error {
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

// AddSendEventHandler 전송 이벤트 핸들러
func AddSendEventHandler(c fiber.Ctx) error {
	body := new(schemas.ReqAddSendEvent)
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
