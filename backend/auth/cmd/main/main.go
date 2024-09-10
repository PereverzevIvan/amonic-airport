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
	// app.Get("/asdf", func(c *fiber.Ctx) error {
	// 	return c.SendString("asdf")
	// })

	controllers.InitControllers(app, conn.Conn)
	log.Info(app.Listen(fmt.Sprintf(":%d", cfg.ConfigServer.Port)))
}
