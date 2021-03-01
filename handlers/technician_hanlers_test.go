package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/routers"
	"github.com/stretchr/testify/assert"
)

var (
	app = fiber.New()
)

func TestCreateTechnician(t *testing.T) {

	app.Use(cors.New())

	apiGroup := app.Group("api/v1")
	app.Use(logger.New())
	routers.TechnicianRouter(apiGroup)
	// app.Get("/tech", handlers.AppointmentHandler.GetAll)

	req := httptest.NewRequest("GET", "/api/v1/technician", nil)

	resp, _ := app.Test(req)

	var techs []models.Technician

	json.NewDecoder(resp.Body).Decode(&techs)

	assert.Equal(t, 200, resp.StatusCode)
	assert.NotEmpty(t, techs)

	// Do something with results:
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body)) // => Hello, World!
	}
}
