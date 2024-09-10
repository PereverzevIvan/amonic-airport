package main

import (
	"fmt"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/config"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/controllers"
	service "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func main() {
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	conn := service.NewStorage(cfg.ConfigDatabase)
	fmt.Println(conn)

	app := fiber.New()

	log.Info(app.Listen(fmt.Sprintf(":%d", cfg.ConfigServer.Port)))
	controllers.InitControllers(app, conn.Conn)
}
