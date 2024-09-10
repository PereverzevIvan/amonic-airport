package controllers

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/middleware"
	mysql_repo "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/repositories/mysql"
	service "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/services"
	usecase "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/usecases"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func InitControllers(app *fiber.App, conn *gorm.DB) {
	api := app.Group("/api")

	userRepo := mysql_repo.NewUserRepo(conn)
	userService := service.NewUserService(userRepo)
	userUseCase := usecase.NewUserUseCase(userService)

	authMiddleware := middleware.NewAuthMiddleware(userUseCase)
	NewUserController(&api, userUseCase, authMiddleware)
}
