package controllers

import (
	"net/http"
	"time"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type SummaryController struct {
	jwtUseCase     jwtUseCase
	summaryService summaryService
	// summaryService summaryService
}

func AddSummaryControllerRoutes(
	api *fiber.Router,
	jwtUseCase jwtUseCase,
	summaryService summaryService,

	authMiddleware AuthMiddleware,
) {
	controller := &SummaryController{
		jwtUseCase:     jwtUseCase,
		summaryService: summaryService,
	}

	(*api).Get("/summary", controller.Generate30DaysSummary, authMiddleware.IsActive)
}

// @Summary      Generate 30 days summary
// @Description  Generate 30 days summary
// @Tags         Summary
// @Accept       json
// @Produce      json
// @Param        SummaryParams query  models.SummaryParams true "example"
// @Success      200 {object}  models.Summary
// @Failure      400
// @Failure      404
func (controller *SummaryController) Generate30DaysSummary(ctx fiber.Ctx) error {
	params := models.SummaryParams{}

	if err := ctx.Bind().Query(&params); err != nil {
		log.Error(err)
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if err := params.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	start_date, err := time.Parse("2006-01-02", params.StartDate)
	if err != nil {
		log.Error(err)
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	params.EndDate = start_date.Add(30 * 24 * time.Hour).Format("2006-01-02")
	log.Info(params)

	time_start_generating := time.Now()

	summary := &models.Summary{}
	err = controller.summaryService.GetFlightsInfo(&params, summary)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	err = controller.summaryService.GetTopCustomersInfo(&params, summary)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	err = controller.summaryService.GetTopFlightsInfo(&params, summary)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	err = controller.summaryService.GetTopOfficesInfo(&params, summary)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	err = controller.summaryService.GetRevenueFromTicketSales(&params, summary)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	err = controller.summaryService.GetWeeklyReportOfPercentageOfEmptySeats(&params, summary)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	summary.TimeTakenToGenerateReport = int(time.Since(time_start_generating).Milliseconds())

	return ctx.Status(http.StatusOK).JSON(summary)
}
