package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/handlers"
)

func StoreRouter(api fiber.Router) {
	storeGroup := api.Group("/stores")
	storeGroup.Post("/", handlers.StoreHandler.Create)
	storeGroup.Get("/", handlers.StoreHandler.GetAll)
	storeGroup.Put("/", handlers.StoreHandler.Update)
}
