package routes

import (
	"ses-go/api/handlers"
	"ses-go/api/middlewares"

	"github.com/gofiber/fiber/v3"
)

// SetTemplateRoutes 템플릿 라우터
func SetTemplateRoutes(app *fiber.App) {
	app.Get("/accounts/login", handlers.LoginHTMLRenderHandler)
	app.Use(middlewares.SessionAuthenticate).Get("/", handlers.IndexHTMLRenderHandler)
	plans := app.Group("/plans", middlewares.SessionAuthenticate)
	{
		plans.Get("", handlers.PlanCreateHTMLRenderHandler)
		plans.Get("/templates/:templateId/recipients/:recipientId", handlers.GetRecipientsHTMLRenderHandler)
		plans.Get("/templates/:templateId", handlers.TemplateDetailHTMLRenderHandler)
		plans.Get("/:planId", handlers.PlanDetailHTMLRenderHandler)
		plans.Get("/:planId/result", handlers.PlanResultHTMLRenderHandler)
	}
	token := app.Group("/tokens", middlewares.SessionAuthenticate)
	{
		token.Get("", handlers.TokenHTMLRenderHandler)
	}
}
