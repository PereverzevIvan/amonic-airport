package controllers

import (
	"net/http"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
)

type OfficeService interface {
	GetByID(id int) (*models.Office, error)
	GetByTitle(title string) (*models.Office, error)
	GetAll() (*[]models.Office, error)
}

type OfficeController struct {
	OfficeService OfficeService
}

func AddOfficeControllerRoutes(router *fiber.Router, s OfficeService) {
	controller := OfficeController{OfficeService: s}

	(*router).Get("/office/:id", controller.GetByID)
	(*router).Get("/offices", controller.GetAll)
}

// Get Office By ID
// @Summary      Get Office by ID
// @Description  Получение информации об офисе по его идентификатору
// @Tags         Office
// @Accept       json
// @Produce      json
// @Param        id path  int  true  "Office ID"
// @Success      200  {object}  models.Office
// @Failure      400
// @Failure      500
// @Router       /office/{id} [get]
func (con OfficeController) GetByID(c fiber.Ctx) error {
	id := fiber.Params(c, "id", 0)
	if id < 1 {
		c.Status(http.StatusBadRequest).SendString("Неверный id")
	}

	office, err := con.OfficeService.GetByID(id)
	if err != nil {
		// TODO: Добавить обработку разных ошибок
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(office)
}

// @Summary      Get all offices
// @Description  Получение информации о всех офисах
// @Tags         Office
// @Success      200  {object}  []models.Office
// @Failure      500
// @Router       /offices [get]
func (con OfficeController) GetAll(c fiber.Ctx) error {
	offices, err := con.OfficeService.GetAll()
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(offices)
}
