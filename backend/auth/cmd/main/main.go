package main

import (
	"fmt"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/config"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/controllers"
	service "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// @title Разработка бизнес-приложений - лаба 1
// @version 1.0
// @description Это API лабораторной работы 1 по дисциплине "Разработка бизнес-приложений".
// @description Тема проекта - аэропорт.
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name access-token
func main() {
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	conn := service.NewStorage(cfg.ConfigDatabase)
	fmt.Println(conn)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://swagger-ui:8080", "http://frontend:3000", "http://localhost:5043"},
		AllowCredentials: true, // Разрешение отправки кук
	}))

	controllers.InitControllers(app, conn.Conn, &cfg.ConfigJWT)

	log.Info(app.Listen(fmt.Sprintf(":%d", cfg.ConfigServer.Port)))
}
