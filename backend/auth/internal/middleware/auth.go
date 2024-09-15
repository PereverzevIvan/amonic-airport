package middleware

import (
	"net/http"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type jwtUseCase interface {
	GetUserIdFromToken(ctx fiber.Ctx, use_refresh_token bool) (int, error)
	// GetVersionFromToken(ctx fiber.Ctx, use_refresh_token bool) (int, error)
}

type userService interface {
	IsActive(user_id int) (bool, error)
	IsAdmin(user_id int) (bool, error)
}

type authMiddleware struct {
	jwtUseCase  jwtUseCase
	userService userService
}

func NewAuthMiddleware(jwtUseCase jwtUseCase, userService userService) *authMiddleware {
	return &authMiddleware{
		jwtUseCase:  jwtUseCase,
		userService: userService,
	}
}

func (am *authMiddleware) IsActive(ctx fiber.Ctx) error {
	user_id, err := am.jwtUseCase.GetUserIdFromToken(ctx, false)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}

	// проверка на активен ли аккаунт пользователя в БД
	is_active, err := am.userService.IsActive(user_id)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	if !is_active {
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	return ctx.Next()
}

func (am *authMiddleware) IsAdmin(ctx fiber.Ctx) error {
	user_id, err := am.jwtUseCase.GetUserIdFromToken(ctx, false)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}

	// проверка на активен ли аккаунт пользователя в БД
	is_active, err := am.userService.IsActive(user_id)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	if !is_active {
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	// проверяем в БД, т.к. администратор мог изменить роль
	is_admin, err := am.userService.IsAdmin(user_id)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if !is_admin {
		return ctx.SendStatus(http.StatusForbidden)
	}

	return ctx.Next()
}
