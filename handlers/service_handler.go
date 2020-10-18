package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/utils/httputils"
)

//ServiceHandler handlers for service route
var ServiceHandler serviceHandlerI = &serviceHandler{}

type serviceHandler struct {
}
type serviceHandlerI interface {
	GetAll(*fiber.Ctx) error
	// GetOne(*fiber.Ctx) error
	Create(*fiber.Ctx) error
}

func (h serviceHandler) GetAll(c *fiber.Ctx) error {
	service := new(models.Service)
	services, err := service.GetAll()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}
	return httputils.JSONSuccessResponse(c, services)
}

func (h serviceHandler) Create(c *fiber.Ctx) error {
	service := new(models.Service)
	if err := c.BodyParser(service); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}
	ID, err := service.Save()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}
	service.ID = *ID

	return httputils.JSONCreatedResponse(c, service)
}
