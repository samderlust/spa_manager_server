package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/handlers"
)

//AppointmentRouter router for appointment
func AppointmentRouter(api fiber.Router) {
	aptGroup := api.Group("/appointment")

	aptGroup.Get("/", handlers.AppointmentHandler.GetAll)
	aptGroup.Post("/", handlers.AppointmentHandler.Create)
}
