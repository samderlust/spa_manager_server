package handlers

import (
	"net/http"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/resources"
	"github.com/samderlust/spa_manager/utils/httputils"
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
	}
	return c.JSON(technicians)

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
	// cxt := context.Background()
	session, _ := resources.MgoClient.StartSession()
	session.StartTransaction()
	//Finding the store
	lUser := c.Locals("user").(*jwt.Token)
	claims := lUser.Claims.(jwt.MapClaims)
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

	user := new(models.User)
	if err := c.BodyParser(&user); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}
	user.Role = "technician"
	user.Password = "12345"
	if err := user.Validate(); err != nil {
		return httputils.JSONResponseModelError(c, err)
	}
	userID, err := user.Save()
	if err != nil {
		session.AbortTransaction(c.Context())
		return httputils.JSONResponseModelError(c, err)
	}
	user.ID = *userID

	var technician models.Technician
	technician.UserID = *userID
	if err := c.BodyParser(&technician); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}

	if err := technician.Validate(); err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	ID, err := technician.Save()
	if err != nil {
		session.AbortTransaction(c.Context())
		return c.Status(err.Status).JSON(err)
	}
	technician.ID = *ID

	//ADD TECHNICIAN ID INTO STORE
	if err := store.AddTechnician(technician); err != nil {
		session.AbortTransaction(c.Context())
		return httputils.JSONResponseModelError(c, err)
	}

	session.CommitTransaction(c.Context())
	session.EndSession(c.Context())
	popTech, _ := technician.Marshall()
	return c.Status(http.StatusCreated).JSON(popTech)
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
