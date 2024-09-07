package main

import (
	"fmt"

	"gitflic.ru/project/pereverzevivan/jwt-auth-golang/config"
	"gitflic.ru/project/pereverzevivan/jwt-auth-golang/internal/services/storage"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func main() {
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	conn := storage.NewStorage(cfg.ConfigDatabase)
	fmt.Println(conn)

	app := fiber.New()

	log.Info(app.Listen(fmt.Sprintf(":%d", cfg.ConfigServer.Port)))
}
