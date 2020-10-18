package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/handlers"
)

//ServiceRouter router for service
func ServiceRouter(api fiber.Router) {
	serviceGroup := api.Group("/service")
	serviceGroup.Get("/", handlers.ServiceHandler.GetAll)
	serviceGroup.Post("/", handlers.ServiceHandler.Create)
}
