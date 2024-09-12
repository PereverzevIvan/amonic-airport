package main

import (
	"fmt"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/config"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/controllers"
	service "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api
func main() {
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	conn := service.NewStorage(cfg.ConfigDatabase)
	fmt.Println(conn)

	app := fiber.New()

	controllers.InitControllers(app, conn.Conn)

	log.Info(app.Listen(fmt.Sprintf(":%d", cfg.ConfigServer.Port)))
}
