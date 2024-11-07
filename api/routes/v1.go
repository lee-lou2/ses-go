package routes

import (
	"ses-go/api/handlers"
	"ses-go/api/middlewares"

	"github.com/gofiber/fiber/v3"
)

// SetV1Routes V1 라우터
func SetV1Routes(app *fiber.App) {
	app.Get("/v1/auth/google", handlers.GoogleAuthHandler)
	app.Get("/v1/auth/google/callback", handlers.GoogleCallbackHandler)
	app.Get("/v1/events/open", handlers.AddOpenEventHandler)
	app.Post("/v1/events/send", handlers.AddSendEventHandler)

	auth := app.Group("/v1/auth", middlewares.SessionAuthenticate)
	{
		auth.Get("/logout", handlers.LogoutHandler)
		auth.Post("/tokens", handlers.CreateTokenHandler)
	}
	plan := app.Group("/v1/plans", middlewares.SessionOrTokenAuthenticate)
	{
		plan.Post("", handlers.CreatePlanHandler)
		plan.Post("/templates", handlers.CreateTemplateHandler)
		plan.Put("/templates/:templateId", handlers.UpdateTemplateHandler)
		plan.Get("/templates/:templateId/fields", handlers.GetTemplateFieldsHandler)
		plan.Post("/templates/:templateId/recipients", handlers.CreateRecipientHandler)
	}
}
