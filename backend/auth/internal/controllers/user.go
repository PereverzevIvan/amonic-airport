package controllers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	userService userService
}

func AddUserControllerRoutes(api *fiber.Router, userService userService, authMiddleware AuthMiddleware) {
	controller := &UserController{userService: userService}

	(*api).Get("/user/:id", controller.GetByID)
	(*api).Get("/users", controller.GetAll)
}

// Get User By ID
// @Summary      Get user by ID
// @Description  Получение информации о пользователе по его идентификатору
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path  int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400
// @Failure      404
// @Router       /user/{id} [get]
func (uc *UserController) GetByID(ctx fiber.Ctx) error {
	user_id := fiber.Params[int](ctx, "id")

	if user_id < 1 {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	user, err := uc.userService.GetByID(user_id)
	if err != nil {
		return ctx.SendString(err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

// @Summary      Get all users
// @Description  Получение информации о всех пользователях
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        office_id query  int  false  "Фильтр по id офиса"
// @Success      200  {object}  []models.User
// @Failure      500
// @Router       /users [get]
func (uc *UserController) GetAll(ctx fiber.Ctx) error {
	filterParams := ctx.Queries()
	fmt.Println(filterParams)

	users, err := uc.userService.GetAll(filterParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(users)
}
