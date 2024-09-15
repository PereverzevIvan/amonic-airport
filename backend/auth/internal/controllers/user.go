package controllers

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
)

type UserUseCase interface {
	GetByID(user_id int) (*models.User, error)
	// GetByEmail(user_email string) (models.User, error)
	// Create(*models.User) error
	// Update(*models.User) error
	// Delete(user_id int) error
}

type AuthMiddleware interface {
	IsAdmin(ctx fiber.Ctx) error
}

type UserController struct {
	userUseCase UserUseCase
}

func AddUserControllerRoutes(api *fiber.Router, userUseCase UserUseCase, authMiddleware AuthMiddleware) {
	controller := &UserController{userUseCase: userUseCase}

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

	user, err := uc.userUseCase.GetByID(user_id)
	if err != nil {
		return ctx.SendString(err.Error())
	}

	return ctx.JSON(user)
}

func (uc *UserController) GetByEmail(ctx fiber.Ctx) error {
	return ctx.SendString("test")
}
