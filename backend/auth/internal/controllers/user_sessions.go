package controllers

import (
	"net/http"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/utils"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type UserSessionController struct {
	jwtUseCase         jwtUseCase
	userService        userService
	userSessionUseCase userSessionUseCase
}

func AddUserSessionControllerRoutes(api *fiber.Router,
	jwtUseCase jwtUseCase,
	userService userService,
	userSessionUseCase userSessionUseCase,
	authMiddleware AuthMiddleware,
) {
	controller := &UserSessionController{
		jwtUseCase:         jwtUseCase,
		userSessionUseCase: userSessionUseCase,
		userService:        userService,
	}

	(*api).Get("user-sessions/", controller.UserSessions, authMiddleware.IsActive)
	(*api).Patch("user-sessions/", controller.SetUnsuccessfullLogoutReason, authMiddleware.IsActive)
}

// Список сессий пользователя
func (controller *UserSessionController) UserSessions(ctx fiber.Ctx) error {
	// Получаем user_id и затем пользователя
	token_user_id, err := controller.jwtUseCase.GetUserIdFromToken(ctx, false)
	if err != nil {
		log.Error(err)
		return nil
	}

	user_id := fiber.Query[int](ctx, "user_id", token_user_id)
	// Если user_id не равен token_user_id,
	// то проверяем, является ли пользователь админом
	if token_user_id != user_id {
		is_admin, err := controller.userService.IsAdmin(user_id)
		if err != nil {
			log.Error(err)
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		if !is_admin {
			return ctx.SendStatus(http.StatusForbidden)
		}
	}

	only_invalid_sessions := fiber.Query[bool](ctx, "only_invalid_sessions", false)
	page := fiber.Query[int](ctx, "page", 1)
	limit := fiber.Query[int](ctx, "limit", 10)

	params := models.UserSessionParams{
		UserID:              user_id,
		OnlyInvalidSessions: only_invalid_sessions,
		Page:                page,
		Limit:               limit,
	}

	log.Info("params: ", params)

	user_sessions, err := controller.userSessionUseCase.GetByUserID(ctx, user_id, &params)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}

	ctx.SendStatus(http.StatusOK)
	return ctx.JSON(user_sessions)
}

// Установка причины неудачного выхода из системы
func (controller *UserSessionController) SetUnsuccessfullLogoutReason(ctx fiber.Ctx) error {
	// Получаем user_id и затем пользователя
	user_id, err := controller.jwtUseCase.GetUserIdFromToken(ctx, false)
	if err != nil {
		log.Error(err)
		return nil
	}

	body := map[string]string{}
	err = ctx.Bind().Body(&body)
	if err != nil {
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString("body is required")
	}

	user_session, err := controller.userSessionUseCase.ObtainUserSessionFromBody(ctx, user_id)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}

	err = controller.userSessionUseCase.Update(ctx, user_session)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}

	return ctx.SendStatus(http.StatusOK)
}
