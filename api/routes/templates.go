package routes

import (
	"ses-go/api/handlers"
	"ses-go/api/middlewares"

	"github.com/gofiber/fiber/v3"
)

// SetTemplateRoutes 템플릿 라우터
func SetTemplateRoutes(app *fiber.App) {
	template := app.Group("")
	{
		template.Get("/accounts/login", handlers.LoginHTMLRenderHandler)
		template.Use(middlewares.SessionAuthenticate).Get("/", handlers.IndexHTMLRenderHandler)
		template.Use(middlewares.SessionAuthenticate).Get("/plans", handlers.PlanCreateHTMLRenderHandler)
		template.Use(middlewares.SessionAuthenticate).Get("/plans/templates/:templateId/recipients/:recipientId", handlers.GetRecipientsHTMLRenderHandler)
		template.Use(middlewares.SessionAuthenticate).Get("/plans/templates/:templateId", handlers.TemplateDetailHTMLRenderHandler)
		template.Use(middlewares.SessionAuthenticate).Get("/plans/:planId", handlers.PlanDetailHTMLRenderHandler)
		template.Use(middlewares.SessionAuthenticate).Get("/plans/:planId/result", handlers.PlanResultHTMLRenderHandler)
	}
}
