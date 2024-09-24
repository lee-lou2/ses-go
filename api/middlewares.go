package api

import (
	"github.com/gofiber/fiber/v3"
	"ses-go/config"
	"ses-go/models"
	"time"
)

func sessionAuthenticate(c fiber.Ctx) error {
	session := c.Cookies("session")
	// 세션 확인
	db := config.GetDB()
	var userSession models.UserSession
	if err := db.Where("uuid = ?", session).Preload("User").First(&userSession).Error; err != nil {
		return c.Redirect().To("/accounts/login")
	}
	// 생성일이 24일 이상이면 세션 삭제
	if userSession.CreatedAt.AddDate(0, 0, 24).Before(time.Now()) {
		db.Delete(&userSession)
		return c.Redirect().To("/accounts/login")
	}
	fiber.Locals[models.User](c, "user", userSession.User)
	fiber.Locals[models.UserSession](c, "session", userSession)
	return c.Next()
}
