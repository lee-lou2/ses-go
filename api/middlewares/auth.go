package middlewares

import (
	"fmt"
	"ses-go/cmd/accounts"
	"ses-go/config"
	"ses-go/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

// ValidateSession 세션 유효성 검사
func ValidateSession(c fiber.Ctx) bool {
	session := c.Cookies("session")
	// 세션 확인
	db := config.GetDB()
	var userSession models.UserSession
	if err := db.Where("uuid = ?", session).Preload("User").First(&userSession).Error; err != nil {
		if session != "" {
			c.ClearCookie("session")
		}
		return false
	}
	// 생성일이 24일 이상이면 세션 삭제
	if userSession.CreatedAt.AddDate(0, 0, 24).Before(time.Now()) {
		db.Delete(&userSession)
		return false
	}
	fiber.Locals[models.User](c, "user", userSession.User)
	fiber.Locals[models.UserSession](c, "session", userSession)
	return true
}

// SessionAuthenticate 세션 인증
func SessionAuthenticate(c fiber.Ctx) error {
	fmt.Println("SessionAuthenticate")
	if !ValidateSession(c) {
		return c.Redirect().To("/accounts/login")
	}
	return c.Next()
}

// ValidateToken 토큰 유효성 검사
func ValidateToken(c fiber.Ctx) bool {
	// Header에서 Authorization 토큰 확인
	token := c.Get("Authorization")
	if token == "" {
		return false
	}
	token = strings.Split(token, " ")[1]
	userId, err := accounts.ValidateToken(token)
	if err != nil {
		return false
	}
	db := config.GetDB()
	var user models.User
	if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
		return false
	}
	fiber.Locals[models.User](c, "user", user)
	return true
}

// SessionOrTokenAuthenticate 세션 또는 토큰 인증
func SessionOrTokenAuthenticate(c fiber.Ctx) error {
	if !ValidateSession(c) {
		if !ValidateToken(c) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}
	return c.Next()
}
