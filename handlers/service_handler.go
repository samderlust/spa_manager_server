package handlers

import (
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/utils/httputils"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	if role != "admin" && role != "owner" {
		return c.Status(401).JSON(&fiber.Map{"message": "Unauthorized"})
	}

	var store models.Store
	ownerID, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	store.Owner = ownerID
	if err := store.FindByOwnerID(); err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	var service models.Service
	if err := c.BodyParser(&service); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}

	ID, err := service.Save()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	service.ID = *ID
	if err := store.AddService(service); err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	return httputils.JSONSuccessResponse(c, service)
}
