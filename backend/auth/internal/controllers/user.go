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

// Get User By ID
// @Summary      Get user by ID
// @Description  Получение информации о пользователе по его идентификатору
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user_id query  int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Router       /user/ [get]
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
