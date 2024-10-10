package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type AirportController struct {
	jwtUseCase     jwtUseCase
	airportService airportService
}

func AddAirportControllerRoutes(
	api *fiber.Router,
	jwtUseCase jwtUseCase,
	airportService airportService,
	authMiddleware AuthMiddleware,
) {
	controller := &AirportController{
		jwtUseCase:     jwtUseCase,
		airportService: airportService,
	}

	(*api).Get("/airports", controller.Airports, authMiddleware.IsActive)
}

// Get All Airports (Список самолётов)
// @Summary      Get All Airports
// @Description  Получение cписка самолётов
// @Tags         Airports
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Airport
// @Failure      400
// @Failure      404
// @Router       /airports [get]
func (controller *AirportController) Airports(ctx fiber.Ctx) error {
	// Получаем список самолётов
	airports, err := controller.airportService.GetAll()
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusOK).JSON(airports)
}
