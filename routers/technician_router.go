package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/handlers"
)

func TechnicianRouter(api fiber.Router) {
	techicianGroup := api.Group("/technician")

	techicianGroup.Get("/", handlers.TechnicianHandler.GetAll)
	techicianGroup.Get("/:id", handlers.TechnicianHandler.GetOne)
	techicianGroup.Post("/", handlers.TechnicianHandler.Create)
	techicianGroup.Delete("/:id", handlers.TechnicianHandler.Delete)

}
