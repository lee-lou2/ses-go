package routes

import (
	"ses-go/api/handlers"
	"ses-go/api/middlewares"

	"github.com/gofiber/fiber/v3"
)

// SetV1Routes V1 라우터
func SetV1Routes(app *fiber.App) {
	v1 := app.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.Get("/google/", handlers.GoogleAuthHandler)
			auth.Get("/google/callback/", handlers.GoogleCallbackHandler)
			auth.Use(middlewares.SessionAuthenticate).Get("/logout", handlers.LogoutHandler)
		}
		event := v1.Group("/events")
		{
			event.Get("/open", handlers.AddOpenEventHandler)
			event.Post("/send", handlers.AddSendEventHandler)
		}
		plan := v1.Group("/plans")
		{
			plan.Use(middlewares.SessionAuthenticate).Post("/", handlers.CreatePlanHandler)
			template := plan.Group("/templates")
			{
				template.Use(middlewares.SessionAuthenticate).Post("/", handlers.CreateTemplateHandler)
				template.Use(middlewares.SessionAuthenticate).Put("/:templateId", handlers.UpdateTemplateHandler)
				template.Use(middlewares.SessionAuthenticate).Get("/:templateId/recipients", handlers.InitRecipientsDataHandler)
				template.Use(middlewares.SessionAuthenticate).Post("/:templateId/recipients", handlers.CreateRecipientHandler)
			}
		}
	}
}
