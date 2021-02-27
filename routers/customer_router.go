package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/handlers"
)

func CustomerRouter(api fiber.Router) {
	cusGroup := api.Group("/customers")

	cusGroup.Get("/", handlers.CustomerHandler.GetAll)
	cusGroup.Post("/", handlers.CustomerHandler.Create)

}
