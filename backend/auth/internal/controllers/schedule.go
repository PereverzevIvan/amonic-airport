package controllers

import (
	"net/http"
	"path/filepath"
	"strings"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type ScheduleController struct {
	jwtUseCase      jwtUseCase
	scheduleService scheduleService
}

func AddScheduleControllerRoutes(
	api *fiber.Router,
	jwtUseCase jwtUseCase,
	scheduleService scheduleService,
	authMiddleware AuthMiddleware,
) {
	controller := &ScheduleController{
		jwtUseCase:      jwtUseCase,
		scheduleService: scheduleService,
	}

	(*api).Get("/schedules", controller.Schedules, authMiddleware.IsActive)
	(*api).Get("/schedule/:id", controller.GetByID, authMiddleware.IsActive)
	(*api).Put("/schedule/:id", controller.UpdateConfirmed, authMiddleware.IsActive)
	(*api).Patch("/schedule/:id", controller.UpdateByID, authMiddleware.IsActive)

	(*api).Post("/schedules/upload", controller.SchedulesUpload, authMiddleware.IsAdmin)
}

// Get All Schedules (Список полетов)
// @Summary      Get All Schedules
// @Description  Получение cписка полетов
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Param        SchedulesParams query  models.SchedulesParams false "example"
// @Success      200  {object}  []models.Schedule
// @Failure      400
// @Failure      404
// @Router       /schedules [get]
func (controller *ScheduleController) Schedules(ctx fiber.Ctx) error {
	var params models.SchedulesParams

	if err := ctx.Bind().Query(&params); err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	// Проверка параметров
	err := params.Validate()
	if err != nil {
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	// Получаем список полетов
	schedules, err := controller.scheduleService.GetAll(&params)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusOK).JSON(schedules)
}

// Get Schedules by id ()
// @Summary      Get Schedules by id
// @Description  Получение cписка полетов
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Schedule
// @Failure      400
// @Failure      404
// @Router       /schedule/{id} [get]
func (controller *ScheduleController) GetByID(ctx fiber.Ctx) error {
	schedule_id := fiber.Params[int](ctx, "id")

	if schedule_id < 1 {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Получаем список полетов
	schedule, err := controller.scheduleService.GetByID(schedule_id)
	if err != nil {
		if err == models.ErrNotFound {
			return ctx.SendStatus(http.StatusNotFound)
		}

		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusOK).JSON(schedule)
}

// Update Schedule confirmed  by id (обновить статус confirmed)
// @Summary      Update Schedule confirmed  by id
// @Description  Обновление статуса confirmed
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Param        SchedulesParams body  models.ScheduleUpdateConfirmedParams true "example"
// @Success      200  string    "ok"
// @Failure      400
// @Failure      404
// @Router       /schedule/{id} [put]
func (controller *ScheduleController) UpdateConfirmed(ctx fiber.Ctx) error {
	schedule_id := fiber.Params[int](ctx, "id")

	if schedule_id < 1 {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	var params models.ScheduleUpdateConfirmedParams
	if err := ctx.Bind().Body(&params); err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	err := params.Validate()
	if err != nil {
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	// Получаем список полетов
	err = controller.scheduleService.UpdateConfirmed(schedule_id, *params.Confirmed)
	if err != nil {
		if err == models.ErrNotFound {
			return ctx.SendStatus(http.StatusNotFound)
		}

		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}

// Update Schedule  by id (обновить)
// @Summary      Update Schedule  by id
// @Description  Обновление
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Param        SchedulesParams body  models.ScheduleUpdateParams true "example"
// @Success      200  string    "ok"
// @Failure      400
// @Failure      404
// @Router       /schedule/{id} [patch]
func (controller *ScheduleController) UpdateByID(ctx fiber.Ctx) error {
	schedule_id := fiber.Params[int](ctx, "id")

	if schedule_id < 1 {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	var params models.ScheduleUpdateParams
	if err := ctx.Bind().Body(&params); err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	err := params.Validate()
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	// Получаем список полетов
	err = controller.scheduleService.UpdateByID(schedule_id, &params)
	if err != nil {
		if err == models.ErrNotFound {
			return ctx.SendStatus(http.StatusNotFound)
		}

		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}

// @Summary      Загрузить файл CSV для обновления или добавления / обновления списка полетов
// @Description  Загрузить CSV файл по ключу "file" (name="file")
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Success      200  {object} models.SchedulesUploadResult
// @Failure      400
// @Failure      404
// @Router       /schedules/upload [post]
func (controller *ScheduleController) SchedulesUpload(ctx fiber.Ctx) error {
	// Get the file from the form input
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Failed to get file")
	}

	if file.Size > 1024*1024*5 {
		return ctx.Status(fiber.StatusBadRequest).SendString("File size must be less than 5 MB")
	}

	if ext := strings.ToLower(filepath.Ext(file.Filename)); ext != ".csv" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Only CSV files are allowed")
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		log.Error(err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to open file")
	}
	defer src.Close()

	// Apply changes from SCV
	res, err := controller.scheduleService.ApplyChangesFromSCV(&src)
	if err != nil {
		log.Error(err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(res)
}
