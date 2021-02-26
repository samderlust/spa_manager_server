package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/bookstore_utils-go/resterrors"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/utils/httputils"
)

//StoreHandler handlers for StoreHandler routes
var StoreHandler storeHandlerI = &storeHandler{}

type storeHandler struct{}
type storeHandlerI interface {
	Create(*fiber.Ctx) error
	GetAll(*fiber.Ctx) error
	SetOwner(*fiber.Ctx) error
	Update(*fiber.Ctx) error
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
