package httputils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/utils/resterrors"
)

//JSONParamInvalidResponse shortcut to response invalid params
func JSONParamInvalidResponse(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusBadRequest).JSON(resterrors.NewBadRequestError(err.Error()))
}
