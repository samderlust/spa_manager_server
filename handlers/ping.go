package handlers

import "github.com/gofiber/fiber/v2"

var PingHandler pingHandlerI = &pingHandler{}

type pingHandler struct{}
type pingHandlerI interface {
	Get(*fiber.Ctx) error
}

func (p pingHandler) Get(c *fiber.Ctx) error {
	return c.SendString("pong")
}
