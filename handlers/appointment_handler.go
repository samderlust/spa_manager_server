package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/utils/httputils"
	"github.com/samderlust/spa_manager/utils/logger"
)

//AppointmentHandler handlers for appointmenr routes
var AppointmentHandler appointmentHanderI = &appointmentHander{}

type appointmentHander struct{}
type appointmentHanderI interface {
	GetAll(*fiber.Ctx) error
	// GetOne(*fiber.Ctx) error
	Create(*fiber.Ctx) error
	// Delete(*fiber.Ctx) error
}

func (h appointmentHander) GetAll(c *fiber.Ctx) error {
	appointment := new(models.Appointment)
	appointments, err := appointment.GetAll()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	return httputils.JSONSuccessResponse(c, appointments)
}
func (h appointmentHander) Create(c *fiber.Ctx) error {
	appointment := new(models.Appointment)

	if err := c.BodyParser(appointment); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}
	appointment.Time = time.Now()

	logger.Info(time.Now().String())
	ID, err := appointment.Save()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	appointment.ID = *ID
	return httputils.JSONCreatedResponse(c, appointment)
}
