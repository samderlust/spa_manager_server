package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/utils/httputils"
)

//GeneralHandler handler for all general routes
var GeneralHandler generalHandlerI = &generalHandler{}

type generalHandler struct{}
type generalHandlerI interface {
	GetGeneralInfo(*fiber.Ctx) error
}

func (h generalHandler) GetGeneralInfo(c *fiber.Ctx) error {
	appointment := new(models.Appointment)
	appointments, err := appointment.GetAll()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	customer := new(models.Customer)
	customers, err := customer.GetAll()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	technician := models.Technician{}
	technicians, err := technician.GetAll()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	service := new(models.Service)
	services, err := service.GetAll()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	return httputils.JSONSuccessResponse(c, &fiber.Map{
		"appointments": appointments,
		"customers":    customers,
		"technicians":  technicians,
		"services":     services,
	})
}
