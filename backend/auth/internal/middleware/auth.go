package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/golang-jwt/jwt"
)

type jwtUseCase interface {
	GetAccessToken(ctx fiber.Ctx) (*jwt.Token, error)
	// ClearJWTCookies(ctx fiber.Ctx) error
	GetTokenUserId(token *jwt.Token) (int, error)
	GetTokenVersion(token *jwt.Token) (int, error)
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
	// Получаем и проверям access_token из запроса
	access_token, err := am.jwtUseCase.GetAccessToken(ctx)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if access_token == nil {
		ctx.SendStatus(http.StatusUnauthorized)
		return ctx.SendString("Invalid access token")
	}

	user_id, err := am.jwtUseCase.GetTokenUserId(access_token)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
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
	access_token, err := am.jwtUseCase.GetAccessToken(ctx)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if access_token == nil {
		ctx.SendStatus(http.StatusUnauthorized)
		return ctx.SendString("Invalid access token")
	}

	// получаем user_id
	user_id, err := am.jwtUseCase.GetTokenUserId(access_token)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
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
