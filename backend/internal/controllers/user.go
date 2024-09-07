package controllers

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
)

type UserService interface {
	GetByID(user_id int) (models.User, error)
	// GetByEmail(user_email string) (models.User, error)
	// Create(*models.User) error
	// Update(*models.User) error
	// Delete(user_id int) error
}

type UserController struct {
	Service UserService
}

func NewUserController(api *fiber.Router, userService *UserService) {
	controller := &UserController{Service: *userService}

	(*api).Get("user/", controller.GetByID)
}

func (uc *UserController) GetByID(ctx fiber.Ctx) error {
	return ctx.SendString("Пошел нахуй")
}
