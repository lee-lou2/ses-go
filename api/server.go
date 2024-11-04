package api

import (
	"log"
	"ses-go/api/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/template/html/v2"
)

// Run 서버 실행
func Run() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Middleware
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))
	app.Use(pprof.New())

	// Routes
	routes.SetV1Routes(app)
	routes.SetTemplateRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
