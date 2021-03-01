package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/handlers"
)

func TechnicianRouter(api fiber.Router) {
	technicianGroup := api.Group("/technician")

	technicianGroup.Get("/", handlers.TechnicianHandler.GetAll)
	technicianGroup.Get("/:id", handlers.TechnicianHandler.GetOne)
	technicianGroup.Post("/", handlers.TechnicianHandler.Create)
	technicianGroup.Delete("/:id", handlers.TechnicianHandler.Delete)

}
