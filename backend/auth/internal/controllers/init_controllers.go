package controllers

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/config"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/middleware"
	mysql_repo "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/repositories/mysql"
	service "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/services"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/usecases"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func InitControllers(app *fiber.App, conn *gorm.DB, jwtConfig *config.ConfigJWT) {
	api := app.Group("/api")

	userRepo := mysql_repo.NewUserRepo(conn)
	userService := service.NewUserService(userRepo)

	tokensVersionRepo := mysql_repo.NewTokensVersionRepo(conn)
	jwtService := service.NewJWTService(jwtConfig, tokensVersionRepo)
	jwtUseCase := usecases.NewJWTUseCase(jwtService)

	authMiddleware := middleware.NewAuthMiddleware(jwtUseCase, userService)

	userSessionRepo := mysql_repo.NewUserSessionRepo(conn)
	userSessionService := service.NewUserSessionService(userSessionRepo)
	userSessionUseCase := usecases.NewUserSessionUseCase(userSessionService)

	AddUserControllerRoutes(&api, userService, authMiddleware)
	AddAuthControllerRoutes(&api, jwtUseCase, userService, userSessionUseCase, authMiddleware)
	AddUserSessionControllerRoutes(&api, jwtUseCase, userService, userSessionUseCase, authMiddleware)

	countryRepo := mysql_repo.NewCountryRepo(conn)
	countryService := service.NewCountryService(countryRepo)
	AddCountryControllerRoutes(&api, countryService)

	officeRepo := mysql_repo.NewOfficeRepo(conn)
	officeService := service.NewOfficeService(officeRepo)
	AddOfficeControllerRoutes(&api, officeService)

	scheduleRepo := mysql_repo.NewScheduleRepo(conn)
	routeRepo := mysql_repo.NewRouteRepo(conn)
	scheduleService := service.NewScheduleService(scheduleRepo, routeRepo)
	AddScheduleControllerRoutes(&api, jwtUseCase, scheduleService, authMiddleware)
}
