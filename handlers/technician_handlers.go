package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var TechnicianHandler technicianHandlerI = &technicianHandler{}

type technicianHandler struct{}
type technicianHandlerI interface {
	GetAll(*fiber.Ctx) error
	GetOne(*fiber.Ctx) error
	Create(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
}

func (h technicianHandler) GetAll(c *fiber.Ctx) error {
	technician := models.Technician{}

	technicians, err := technician.GetAll()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	} else {
		return c.JSON(technicians)
	}
}

func (h technicianHandler) GetOne(c *fiber.Ctx) error {
	techID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	technician := models.Technician{}
	technician.ID = techID
	if err := technician.GetByID(); err != nil {
		return c.Status(err.Status).JSON(err)
	}
	return c.JSON(technician)
}

func (h technicianHandler) Create(c *fiber.Ctx) error {
	technician := new(models.Technician)
	if err := c.BodyParser(&technician); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "Invalid Request",
			},
		)
	}
	ID, err := technician.Save()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}
	technician.ID = *ID
	return c.Status(http.StatusCreated).JSON(technician)
}

func (h technicianHandler) Delete(c *fiber.Ctx) error {
	techID, _ := primitive.ObjectIDFromHex(c.Params("id"))

	technician := new(models.Technician)

	technician.ID = techID
	if err := technician.Delete(); err != nil {
		return c.Status(err.Status).JSON(err)
	}
	return c.JSON(technician)

}

func (h technicianHandler) Update(c *fiber.Ctx) error {
	technician := new(models.Technician)

	if err := c.BodyParser(&technician); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	if err := technician.Update(); err != nil {
		return c.Status(err.Status).JSON(err)
	}
	return c.Status(http.StatusOK).JSON(technician)
}
