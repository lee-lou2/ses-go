package handlers

import (
	"fmt"
	"ses-go/cmd/accounts"
	"ses-go/models"

	"github.com/gofiber/fiber/v3"
)

// CreateTokenHandler 토큰 생성 핸들러
func CreateTokenHandler(c fiber.Ctx) error {
	user := fiber.Locals[models.User](c, "user")
	fmt.Println(user)
	token, err := accounts.GetToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"token": token})
}

// DeleteTokenHandler 토큰 삭제 핸들러
func DeleteTokenHandler(c fiber.Ctx) error {
	user := fiber.Locals[models.User](c, "user")
	accounts.DeleteToken(user.ID)
	return c.JSON(fiber.Map{})
}
