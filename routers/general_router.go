package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/handlers"
)

//GeneralRouter setup router for general routes
func GeneralRouter(api fiber.Router) {
	genGroup := api.Group("/general")

	genGroup.Get("/", handlers.GeneralHandler.GetGeneralInfo)
}
