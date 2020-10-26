package middlewares

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func AuthMiddleware() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte("klsajflkdasdsfdsfsdjf"),
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			return c.Status(401).JSON(&fiber.Map{"message": "Unauthorized"})
		},
	})
}
