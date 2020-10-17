package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/utils/httputils"
)

var AuthHandler authHandlerI = &authHandler{}

type authHandler struct{}
type authHandlerI interface {
	SignUp(*fiber.Ctx) error
}

func (h authHandler) SignUp(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}
	ID, err := user.Save()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}
	user.ID = *ID
	user = user.Marshall()

	return c.Status(http.StatusCreated).JSON(user)
}
