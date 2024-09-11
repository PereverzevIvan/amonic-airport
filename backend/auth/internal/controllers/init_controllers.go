package controllers

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/middleware"
	mysql_repo "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/repositories/mysql"
	service "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/services"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func InitControllers(app *fiber.App, conn *gorm.DB) {
	api := app.Group("/api")

	userRepo := mysql_repo.NewUserRepo(conn)
	userService := service.NewUserService(userRepo)

	authMiddleware := middleware.NewAuthMiddleware(userService)
	AddUserControllerRoutes(&api, userService, authMiddleware)

	countryRepo := mysql_repo.NewCountryRepo(conn)
	countryService := service.NewCountryService(countryRepo)
	AddCountryControllerRoutes(&api, countryService)

	officeRepo := mysql_repo.NewOfficeRepo(conn)
	officeService := service.NewOfficeService(officeRepo)
	AddOfficeControllerRoutes(&api, officeService)
}
