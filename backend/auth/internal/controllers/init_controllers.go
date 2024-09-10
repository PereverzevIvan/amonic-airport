package controllers

import (
	service "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/services"
	usecase "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/usecases"
	"github.com/gofiber/fiber/v3"
)

func InitControllers(app *fiber.App) {
	api := app.Group("/api")

	userService := service.NewUserService()
	userUseCase := usecase.NewUserUseCase(userService)
	NewUserController(&api, userUseCase)
}
