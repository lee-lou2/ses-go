package handlers

import (
	"ses-go/config"
	"ses-go/models"
	"ses-go/pkg/google"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// GoogleAuthHandler 구글 인증 핸들러
func GoogleAuthHandler(c fiber.Ctx) error {
	path := google.ConfigGoogle()
	url := path.AuthCodeURL("state")
	return c.Redirect().To(url)
}

// GoogleCallbackHandler 구글 콜백 핸들러
func GoogleCallbackHandler(c fiber.Ctx) error {
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
		Secure:   true,
		SameSite: "Lax",
	})
	return c.Redirect().To("/")
}

// LogoutHandler 로그아웃 핸들러
func LogoutHandler(c fiber.Ctx) error {
	db := config.GetDB()
	session := fiber.Locals[models.UserSession](c, "session")
	db.Delete(&session)
	c.ClearCookie("session")
	return c.Redirect().To("/")
}
