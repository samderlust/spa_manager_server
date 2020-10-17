package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/handlers"
)

func PingRoutes(app fiber.Router) {
	app.Get("/ping", handlers.PingHandler.Get)
}
