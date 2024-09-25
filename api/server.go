package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/template/html/v2"
	"log"
)

// Run 서버 실행
func Run() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))
	app.Use(pprof.New())

	// API
	v1 := app.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.Get("/google/", googleAuthHandler)
			auth.Get("/google/callback/", googleCallbackHandler)
			auth.Use(sessionAuthenticate).Get("/logout", logoutHandler)
		}
		event := v1.Group("/events")
		{
			event.Get("/open", addOpenEventHandler)
			event.Post("/send", addSendEventHandler)
		}
		plan := v1.Group("/plans")
		{
			plan.Use(sessionAuthenticate).Post("/", createPlanHandler)
			plan.Use(sessionAuthenticate).Post("/users/share", createSheetAndShareHandler)
			template := plan.Group("/templates")
			{
				template.Use(sessionAuthenticate).Post("/", createTemplateHandler)
				template.Use(sessionAuthenticate).Put("/:id", updateTemplateHandler)
			}
		}
	}

	// Templates
	template := app.Group("")
	{
		template.Get("/accounts/login", loginTemplateHandler)
		template.Use(sessionAuthenticate).Get("/", indexTemplateHandler)
		template.Use(sessionAuthenticate).Get("/plans", planCreateTemplateHandler)
		template.Use(sessionAuthenticate).Get("/plans/:id", planDetailTemplateHandler)
		template.Use(sessionAuthenticate).Get("/plans/:id/result", planResultTemplateHandler)
		template.Use(sessionAuthenticate).Get("/plans/templates/:id", templateDetailTemplateHandler)
	}

	log.Fatal(app.Listen(":3000"))
}
