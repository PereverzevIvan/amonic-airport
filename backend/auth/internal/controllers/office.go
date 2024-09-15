package controllers

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
)

type OfficeService interface {
	GetByID(id int) (*models.Office, error)
	GetByTitle(title string) (*models.Office, error)
}

type OfficeController struct {
	OfficeService OfficeService
}

func AddOfficeControllerRoutes(router *fiber.Router, s OfficeService) {
	api := (*router).Group("/office")

	controller := OfficeController{OfficeService: s}

	api.Get("/:id", controller.GetByID)
	api.Get("/title/:title", controller.GetByTitle)
}

// Get Office By ID
// @Summary      Get Office by ID
// @Description  Получение информации об офисе по его идентификатору
// @Tags         Office
// @Accept       json
// @Produce      json
// @Param        id path  int  true  "Office ID"
// @Success      200  {object}  models.Office
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Router       /office/{id} [get]
func (con OfficeController) GetByID(c fiber.Ctx) error {
	id := fiber.Params(c, "id", 0)
	if id < 1 {
		c.Status(fiber.StatusBadRequest).SendString("Неверный id")
	}

	office, err := con.OfficeService.GetByID(id)
	if err != nil {
		// TODO: Добавить обработку разных ошибок
		return c.Status(fiber.StatusNotFound).SendString("Не удалось получить офис")
	}

	return c.Status(fiber.StatusOK).JSON(office)
}
