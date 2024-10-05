package controllers

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
)

type userService interface {
	GetByID(user_id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	IsPasswordCorrect(user *models.User, password string) bool
	GetAll(params map[string]string) (*[]models.User, error)
	IsActive(user_id int) (bool, error)
	IsAdmin(user_id int) (bool, error)
	Create(user *models.User) error
	Update(user *models.User) error
	UpdateActive(user_id int, is_active bool) error
}

type jwtUseCase interface {
	RefreshTokens(ctx fiber.Ctx, user *models.User) error
	// GetAccessToken(ctx fiber.Ctx) (*jwt.Token, error)
	// GetRefreshToken(ctx fiber.Ctx) (*jwt.Token, error)
	ClearJWTCookies(ctx fiber.Ctx) error
	GetTokenUserId(token *jwt.Token) (int, error)
	GetUserIdFromToken(ctx fiber.Ctx, use_refresh_token bool) (int, error)
}

type userSessionUseCase interface {
	ObtainUserSessionFromBody(ctx fiber.Ctx, user_id int) (*models.UserSession, error)
	CreateNewLoginSession(ctx fiber.Ctx, user_id int) error
	LogoutLastSession(ctx fiber.Ctx, user_id int) error
	UpdateNoLogoutSession(ctx fiber.Ctx, user_id int) error

	GetByUserID(ctx fiber.Ctx, user_id int, params *models.UserSessionParams) (*[]models.UserSession, error)
	GetByID(ctx fiber.Ctx, id int) (*models.UserSession, error)
	GetLastByUserId(ctx fiber.Ctx, user_id int) (*models.UserSession, error)
	Create(ctx fiber.Ctx, user_id int) (*models.UserSession, error)
	Update(ctx fiber.Ctx, user_session *models.UserSession) error
}

type AuthMiddleware interface {
	IsActive(ctx fiber.Ctx) error
	IsAdmin(ctx fiber.Ctx) error
}
