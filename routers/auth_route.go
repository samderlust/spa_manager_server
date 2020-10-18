package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/handlers"
)

func AuthRouter(api fiber.Router) {
	authGroup := api.Group("/auth")

	authGroup.Post("/signup", handlers.AuthHandler.SignUp)
	authGroup.Post("/signin", handlers.AuthHandler.SignIn)
}
