package controllers

import (
	"net/http"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type TicketController struct {
	jwtUseCase      jwtUseCase
	scheduleService scheduleService
	ticketService   ticketService
}

func AddTicketControllerRoutes(
	api *fiber.Router,
	jwtUseCase jwtUseCase,
	scheduleService scheduleService,
	ticketService ticketService,

	authMiddleware AuthMiddleware,
) {
	controller := &TicketController{
		jwtUseCase:      jwtUseCase,
		scheduleService: scheduleService,
		ticketService:   ticketService,
	}

	(*api).Post("/tickets/remaining-seats-count", controller.RemainingSeatsCount, authMiddleware.IsActive)
	(*api).Post("/tickets/booking", controller.BookTickets, authMiddleware.IsActive)
	(*api).Post("/tickets/confirm", controller.ConfirmTickets, authMiddleware.IsActive)
}

// @Summary      Check enough tickets
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Param        TicketsCountRemainingSeatsParams body  models.TicketsCountRemainingSeatsParams true "example"
// @Success      200  {object}  models.TicketsRemainingSeatsCount
// @Failure      400
// @Failure      404
// @Router       /tickets/remaining-seats-count [post]
func (controller *TicketController) RemainingSeatsCount(ctx fiber.Ctx) error {
	var params models.TicketsCountRemainingSeatsParams
	if err := ctx.Bind().Body(&params); err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	// Проверка параметров
	err := params.Validate()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	/*
		1. obtain schedules with airplanes, which containes count and type of seats,
		2. obtain count confirmed tickets of each type for each schedule
	*/

	result, err := controller.ticketService.CountRemainingSeats(&params)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	return ctx.Status(http.StatusOK).JSON(result)
}

// @Summary      Book tickets
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Param        TicketsBookingParams body  models.TicketsBookingParams true "example"
// @Success      200  {object}  models.TicketsBookResult
// @Failure      400
// @Failure      404
// @Router       /tickets/booking [post]
func (controller *TicketController) BookTickets(ctx fiber.Ctx) error {
	var params models.TicketsBookingParams
	if err := ctx.Bind().Body(&params); err != nil {
		log.Error(err)
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Проверка параметров
	err := params.Validate()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	params.UserID, err = controller.jwtUseCase.GetUserIdFromToken(ctx, false)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	tickets, err := controller.ticketService.BookTickets(&params)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	var total_cost float64 = 0
	for _, ticket := range tickets {
		total_cost += ticket.Schedule.EconomyPrice
	}

	switch params.CabinType {
	case int(models.KCabinTypeBusiness):
		total_cost = total_cost * models.KBusinessRatio
	case int(models.KCabinTypeFirstClass):
		total_cost = total_cost * models.KFirstClassRation
	}

	// Удалить лишние данные перед отправкой
	for _, ticket := range tickets {
		ticket.Schedule = nil
	}

	tickets_and_total_cost := models.TicketsBookResult{
		Tickets:   tickets,
		TotalCost: total_cost,
	}

	return ctx.Status(http.StatusCreated).JSON(tickets_and_total_cost)
}

// @Summary      Confirm tickets
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Param        TicketIDsParams body  models.TicketIDsParams true "example"
// @Success      200
// @Failure      400
// @Failure      404
// @Router       /tickets/confirm [post]
func (controller *TicketController) ConfirmTickets(ctx fiber.Ctx) error {
	var params models.TicketIDsParams

	if err := ctx.Bind().Body(&params); err != nil {
		log.Error(err)
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Проверка параметров
	err := params.Validate()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err = controller.ticketService.ChangeTicketsStatus(params.TicketeIDs, true)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}

// func (controller *TicketController) CancelBooking(ctx fiber.Ctx) error {
// 	return ctx.SendStatus(http.StatusNotImplemented)
// }
