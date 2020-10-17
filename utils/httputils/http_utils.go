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

//JSONResponseModelError shortcut to response to invalid error from db
func JSONResponseModelError(c *fiber.Ctx, err *resterrors.RestError) error {
	return c.Status(err.Status).JSON(&err)
}

//JSONCreatedResponse shortcut to response to create request
func JSONCreatedResponse(c *fiber.Ctx, res interface{}) error {
	return c.Status(http.StatusCreated).JSON(res)
}

//JSONSuccessResponse shortcut to restponse to successful request
func JSONSuccessResponse(c *fiber.Ctx, res interface{}) error {
	return c.Status(http.StatusOK).JSON(res)
}
