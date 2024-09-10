package controllers

import (
	service "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/services"
	"github.com/gofiber/fiber/v3"
)

func InitControllers(app *fiber.App) {
	api := app.Group("/api")

	userService := service.NewUserService()
	NewUserController(&api, userService)
}
