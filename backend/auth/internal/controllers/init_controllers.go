package controllers

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/config"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/middleware"
	mysql_repo "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/repositories/mysql"
	service "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/services"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func InitControllers(app *fiber.App, conn *gorm.DB, jwtConfig *config.ConfigJWT) {
	api := app.Group("/api")

	userRepo := mysql_repo.NewUserRepo(conn)
	userService := service.NewUserService(userRepo)

	tokensVersionRepo := mysql_repo.NewTokensVersionRepo(conn)
	jwtService := service.NewJWTService(jwtConfig, tokensVersionRepo)
	jwtUseCase := usecase.NewJWTUseCase(jwtService)

	authMiddleware := middleware.NewAuthMiddleware(jwtUseCase, userService)

	AddUserControllerRoutes(&api, userService, authMiddleware)
	AddAuthControllerRoutes(&api, jwtUseCase, userService, authMiddleware)
}
