package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/utils/httputils"
)

var CustomerHandler customerHandlerI = &customerHandler{}

type customerHandler struct{}
type customerHandlerI interface {
	GetAll(*fiber.Ctx) error
	Create(*fiber.Ctx) error
	// Update(*fiber.Ctx) error
	// GetOne(*fiber.Ctx) error
}

func (h customerHandler) GetAll(c *fiber.Ctx) error {
	customer := new(models.Customer)
	customers, err := customer.GetAll()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}
	return httputils.JSONSuccessResponse(c, customers)
}

func (h customerHandler) Create(c *fiber.Ctx) error {
	customer := new(models.Customer)

	if err := c.BodyParser(customer); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}
	ID, err := customer.Save()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	customer.ID = *ID

	return httputils.JSONCreatedResponse(c, customer)
}
