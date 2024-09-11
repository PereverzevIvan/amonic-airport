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
	api := (*router).Group("/country")

	controller := CountryController{CountryService: s}

	api.Get("/:id", controller.GetByID)
	api.Get("/name/:name", controller.GetByName)
}

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

func (con CountryController) GetByName(c fiber.Ctx) error {
	name := fiber.Params[string](c, "name")
	if name == "" {
		c.Status(fiber.StatusBadRequest).SendString("title страны обязателен")
	}

	country, err := con.CountryService.GetByName(name)
	if err != nil {
		// TODO: Добавить обработку разных ошибок
		return c.Status(fiber.StatusNotFound).SendString("Не удалось получить страну")
	}

	return c.Status(fiber.StatusOK).JSON(country)
}
