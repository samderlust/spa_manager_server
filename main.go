package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/samderlust/spa_manager/routers"
)

func main() {
	app := fiber.New()

	apiGroup := app.Group("api/v1")
	app.Use(logger.New())
	routers.PingRoutes(apiGroup)
	routers.TechnicianRouter(apiGroup)
	app.Listen(":8080")
}
