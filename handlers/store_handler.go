package handlers

import (
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/bookstore_utils-go/resterrors"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/utils/httputils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//StoreHandler handlers for StoreHandler routes
var StoreHandler storeHandlerI = &storeHandler{}

type storeHandler struct{}
type storeHandlerI interface {
	Create(*fiber.Ctx) error
	GetAll(*fiber.Ctx) error
	SetOwner(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	GetOwnerOverView(*fiber.Ctx) error
	AddEmployee(*fiber.Ctx) error
	AddService(*fiber.Ctx) error
}

func (h storeHandler) Create(c *fiber.Ctx) error {
	store := new(models.Store)
	if err := c.BodyParser(store); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}
	if err := store.Validate(); err != nil {
		return httputils.JSONResponseModelError(c, err)
	}
	ID, err := store.Save()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}
	store.ID = *ID
	return httputils.JSONCreatedResponse(c, store)
}

func (h storeHandler) GetAll(c *fiber.Ctx) error {
	store := new(models.Store)

	stores, err := store.GetAll()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	return httputils.JSONSuccessResponse(c, stores)
}
func (h storeHandler) SetOwner(c *fiber.Ctx) error {
	storeID := c.Params("storeId")
	userID := c.Params("userId")
	if storeID == "" || userID == "" {
		return httputils.JSONParamInvalidResponse(c, resterrors.NewBadRequestError("missing storeId or userId"))
	}
	return httputils.JSONSuccessResponse(c, c.Params)
}

func (h storeHandler) Update(c *fiber.Ctx) error {
	updatingStore := new(models.Store)
	if err := c.BodyParser(updatingStore); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}
	if err := updatingStore.Validate(); err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	if err := updatingStore.Update(); err != nil {
		return httputils.JSONResponseModelError(c, err)
	}
	return httputils.JSONSuccessResponse(c, *updatingStore.Puplate())

}
func (h storeHandler) GetOwnerOverView(c *fiber.Ctx) error {
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

	var apt models.Appointment
	year, month, day := time.Now().Date()

	startTime := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
	endTime := time.Date(year, month, day, 23, 59, 59, 0, time.Now().Location())

	apts, err := apt.FindInTimeRange(startTime, endTime)
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	return httputils.JSONSuccessResponse(c, fiber.Map{
		"store": store.Puplate(),
		"apt":   apts,
	})
}

// AddEmployee
func (h storeHandler) AddEmployee(c *fiber.Ctx) error {
	return c.Status(404).SendString("to be implemented")

}
func (h storeHandler) AddService(c *fiber.Ctx) error {
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
	if err := c.BodyParser(service); err != nil {
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
