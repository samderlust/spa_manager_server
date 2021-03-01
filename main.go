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

	print("start main")
	app := fiber.New()
	app.Use(cors.New())
	apiGroup := app.Group("api/v1")
	app.Use(logger.New())
	routers.AuthRouter(apiGroup)

	app.Use(middlewares.AuthMiddleware())

	routers.PingRoutes(apiGroup)
	routers.TechnicianRouter(apiGroup)
	routers.StoreRouter(apiGroup)
	routers.CustomerRouter(apiGroup)
	routers.ServiceRouter(apiGroup)
	routers.GeneralRouter(apiGroup)
	routers.AppointmentRouter(apiGroup)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	app.Listen(":" + port)

}
