package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type UserController struct {
	userService userService
}

func AddUserControllerRoutes(api *fiber.Router, userService userService, authMiddleware AuthMiddleware) {
	controller := &UserController{userService: userService}

	(*api).Get("/user/:id", controller.GetByID)
	(*api).Get("/users", controller.GetAll, authMiddleware.IsAdmin)
	(*api).Post("/user", controller.Create, authMiddleware.IsAdmin)
	(*api).Patch("/user", controller.Update, authMiddleware.IsAdmin)
	(*api).Put("/user/active", controller.UpdateIsActive, authMiddleware.IsAdmin)
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
		// log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return err
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

	users, err := uc.userService.GetAll(filterParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(users)
}

func (uc *UserController) Create(ctx fiber.Ctx) error {
	var new_user_params models.NewUserParams
	if err := ctx.Bind().Body(&new_user_params); err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		return nil
	}

	new_user, err := new_user_params.ToUser()
	if err != nil {
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	err = uc.userService.Create(new_user)
	if err != nil {
		if errors.Is(err, models.ErrDuplicatedEmail) {
			ctx.Status(http.StatusConflict)
			return ctx.SendString("user with this email already exists")
		}

		if errors.Is(err, models.ErrFKOfficeIDNotFound) {
			ctx.Status(http.StatusConflict)
			return ctx.SendString(fmt.Sprintf("office with id: %v not found", new_user.OfficeID))
		}

		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return nil
	}

	return ctx.Status(http.StatusCreated).JSON(new_user)
}

func (uc *UserController) Update(ctx fiber.Ctx) error {
	var update_user_params models.UpdateUserParams
	if err := ctx.Bind().Body(&update_user_params); err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		return nil
	}

	updated_user, err := update_user_params.ToUser()
	if err != nil {
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	err = uc.userService.Update(updated_user)
	if err != nil {
		if errors.Is(err, models.ErrDuplicatedEmail) {
			ctx.Status(http.StatusConflict)
			return ctx.SendString("user with this email already exists")
		}

		if errors.Is(err, models.ErrFKOfficeIDNotFound) {
			ctx.Status(http.StatusConflict)
			return ctx.SendString(fmt.Sprintf("office with id: %v not found", update_user_params.OfficeID))
		}

		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return nil
	}

	updated_user, err = uc.userService.GetByID(updated_user.ID)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return nil
	}

	return ctx.Status(http.StatusOK).JSON(updated_user)
}

func (uc *UserController) UpdateIsActive(ctx fiber.Ctx) error {
	var params models.UpdateIsActiveUserParams
	if err := ctx.Bind().Body(&params); err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		return nil
	}

	err := uc.userService.UpdateActive(params.ID, params.IsActive)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return nil
	}

	return ctx.SendStatus(http.StatusOK)
}
