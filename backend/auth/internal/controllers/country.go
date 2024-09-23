package controllers

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
)

type CountryService interface {
	GetByID(id int) (*models.Country, error)
	GetByName(title string) (*models.Country, error)
}

type CountryController struct {
	CountryService CountryService
}

func AddCountryControllerRoutes(router *fiber.Router, s CountryService) {
	controller := CountryController{CountryService: s}

	(*router).Get("/country/:id", controller.GetByID)
}

// Get Country By ID
// @Summary      Get Country by id
// @Description  Получение информации о стране по ее числовому ID
// @Tags         Country
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Country ID"
// @Success      200  {object}  models.Country
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Router       /country/{id} [get]
func (con CountryController) GetByID(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id")
	if id == 0 {
		c.Status(fiber.StatusBadRequest).SendString("id страны обязателен")
	}

	country, err := con.CountryService.GetByID(id)
	if err != nil {
		// TODO: Добавить обработку разных ошибок
		return c.Status(fiber.StatusNotFound).SendString("Не удалось получить страну")
	}

	return c.Status(fiber.StatusOK).JSON(country)
}
