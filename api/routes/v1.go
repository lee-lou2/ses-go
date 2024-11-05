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
			auth.Get("/logout", middlewares.SessionAuthenticate, handlers.LogoutHandler)
			auth.Post("/tokens", middlewares.SessionAuthenticate, handlers.CreateTokenHandler)
		}
		event := v1.Group("/events")
		{
			event.Get("/open", handlers.AddOpenEventHandler)
			event.Post("/send", handlers.AddSendEventHandler)
		}
		plan := v1.Group("/plans")
		{
			plan.Post("/", middlewares.SessionOrTokenAuthenticate, handlers.CreatePlanHandler)
			template := plan.Group("/templates")
			{
				template.Post("/", middlewares.SessionAuthenticate, handlers.CreateTemplateHandler)
				template.Put("/:templateId", middlewares.SessionAuthenticate, handlers.UpdateTemplateHandler)
				template.Get("/:templateId/fields", middlewares.SessionAuthenticate, handlers.GetTemplateFieldsHandler)
				template.Post("/:templateId/recipients", middlewares.SessionOrTokenAuthenticate, handlers.CreateRecipientHandler)
			}
		}
	}
}
