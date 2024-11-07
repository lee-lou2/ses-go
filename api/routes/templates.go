package routes

import (
	"ses-go/api/handlers"
	"ses-go/api/middlewares"

	"github.com/gofiber/fiber/v3"
)

// SetTemplateRoutes 템플릿 라우터
func SetTemplateRoutes(app *fiber.App) {
	app.Get("/accounts/login", handlers.LoginHTMLRenderHandler)
	template := app.Group("", middlewares.SessionAuthenticate)
	{
		template.Get("/", handlers.IndexHTMLRenderHandler)
		template.Get("/plans", handlers.PlanCreateHTMLRenderHandler)
		template.Get("/plans/templates/:templateId/recipients/:recipientId", handlers.GetRecipientsHTMLRenderHandler)
		template.Get("/plans/templates/:templateId", handlers.TemplateDetailHTMLRenderHandler)
		template.Get("/plans/:planId", handlers.PlanDetailHTMLRenderHandler)
		template.Get("/plans/:planId/result", handlers.PlanResultHTMLRenderHandler)
		template.Get("/tokens", handlers.TokenHTMLRenderHandler)
	}
}
