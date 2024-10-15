package controllers

import (
	"net/http"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type AmenityController struct {
	jwtUseCase     jwtUseCase
	amenityService amenityService
}

func AddAmenityControllerRoutes(
	api *fiber.Router,
	jwtUseCase jwtUseCase,
	amenityService amenityService,

	authMiddleware AuthMiddleware,
) {
	controller := &AmenityController{
		jwtUseCase:     jwtUseCase,
		amenityService: amenityService,
	}

	(*api).Get("/amenities", controller.GetAll, authMiddleware.IsActive)
	(*api).Get("/amenities/count", controller.CountAll, authMiddleware.IsActive)

	(*api).Get("/cabin-type-default-amenities", controller.GetCabinTypeDefaultAmenities, authMiddleware.IsActive)

	(*api).Get("/ticket-amenities", controller.GetTicketAmenities, authMiddleware.IsActive)
	(*api).Post("/ticket-amenities/edit", controller.EditTicketAmenities, authMiddleware.IsActive)
}

// @Summary      Получить все услуги
// @Tags         Amenity
// @Accept       json
// @Produce      json
// @Success      200 {object}  []models.Amenity
// @Router       /amenities [get]
func (controller *AmenityController) GetAll(ctx fiber.Ctx) error {
	amenities, err := controller.amenityService.GetAll()
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusOK).JSON(amenities)
}

// @Summary      Получить включенные услуги для типа кабины
// @Tags         Amenity
// @Accept       json
// @Produce      json
// @Param        GetCabinTypeDefaultAmenitiesParams query  models.GetCabinTypeDefaultAmenitiesParams true "example"
// @Success      200 {object}  []int
// @Router       /cabin-type-default-amenities [get]
func (controller *AmenityController) GetCabinTypeDefaultAmenities(ctx fiber.Ctx) error {
	params := models.GetCabinTypeDefaultAmenitiesParams{}

	if err := ctx.Bind().Query(&params); err != nil {
		log.Error(err)
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if err := params.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	amenities_ids, err := controller.amenityService.GetCabinTypeDefaultAmenities(&params)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusOK).JSON(amenities_ids)
}

// @Summary      Получить купленные услуги для билета
// @Tags         Amenity
// @Accept       json
// @Produce      json
// @Param        GetTicketAmenitiesParams query  models.GetTicketAmenitiesParams true "example"
// @Success      200 {object}  []models.AmenityTicket
// @Router       /ticket-amenities [get]
func (controller *AmenityController) GetTicketAmenities(ctx fiber.Ctx) error {
	params := models.GetTicketAmenitiesParams{}

	if err := ctx.Bind().Query(&params); err != nil {
		log.Error(err)
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if err := params.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	ticket_amenities, err := controller.amenityService.GetTicketAmenities(&params)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusOK).JSON(ticket_amenities)
}
func (controller *AmenityController) EditTicketAmenities(ctx fiber.Ctx) error {
	params := models.EditTicketAmenitiesParams{}

	if err := ctx.Bind().Body(&params); err != nil {
		log.Error(err)
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if err := params.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err := controller.amenityService.EditTicketAmenities(&params)
	if err != nil {
		log.Error(err)
		if err.Error() == models.ErrCantEditAmenitiesTimeExpired.Error() {
			return ctx.Status(http.StatusBadRequest).SendString("Время на редактирование услуг истекло")
		}
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}

func (controller *AmenityController) CountAll(ctx fiber.Ctx) error {
	params := models.AmenityCountAllParams{}
	if err := ctx.Bind().Query(&params); err != nil {
		log.Error(err)
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if err := params.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	count, err := controller.amenityService.CountAll(&params)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusOK).JSON(count)
}
