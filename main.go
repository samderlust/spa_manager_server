package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/samderlust/spa_manager/middlewares"
	"github.com/samderlust/spa_manager/routers"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	apiGroup := app.Group("api/v1")
	app.Use(logger.New())
	routers.AuthRouter(apiGroup)

	app.Use(middlewares.AuthMiddleware())

	routers.PingRoutes(apiGroup)
	routers.TechnicianRouter(apiGroup)
	routers.CustomerRouter(apiGroup)
	routers.ServiceRouter(apiGroup)
	routers.GeneralRouter(apiGroup)

	app.Listen(":" + os.Getenv("PORT"))

}
