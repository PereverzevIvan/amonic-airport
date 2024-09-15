package controllers

import (
	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	userService userService
}

func AddUserControllerRoutes(api *fiber.Router, userService userService, authMiddleware AuthMiddleware) {
	controller := &UserController{userService: userService}

	(*api).Get("user/", controller.GetByID, authMiddleware.IsAdmin)
}

func (uc *UserController) GetByID(ctx fiber.Ctx) error {
	user_id := fiber.Query[int](ctx, "user_id")

	user, err := uc.userService.GetByID(user_id)
	if err != nil {
		return ctx.SendString(err.Error())
	}

	return ctx.JSON(user)
}

func (uc *UserController) GetByEmail(ctx fiber.Ctx) error {
	return ctx.SendString("test")
}
