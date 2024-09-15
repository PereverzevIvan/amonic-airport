package usecase

import (
	"time"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
)

type jwtService interface {
	GenerateAccessAndRefreshTokens(user *models.User) (string, string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GetAccessTokenExpiration() time.Duration
	GetRefreshTokenExpiration() time.Duration
}

type JWTUseCase struct {
	jwtService jwtService
}

func NewJWTUseCase(jwtService jwtService) JWTUseCase {
	return JWTUseCase{
		jwtService: jwtService,
	}
}

// Обновление токенов:
// 1. Создание новых токенов
// 2. Запись новых токенов в куки
func (j JWTUseCase) RefreshTokens(ctx fiber.Ctx, user *models.User) error {

	access_token, refresh_token, err := j.jwtService.GenerateAccessAndRefreshTokens(user)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:  "access-token",
		Value: access_token,
		// Secure:   true,
		SameSite: "Lax",
		Expires:  time.Now().Add(j.jwtService.GetAccessTokenExpiration()),
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh-token",
		Value:    refresh_token,
		HTTPOnly: true,
		// Secure:   true,
		SameSite: "Lax",
		Expires:  time.Now().Add(j.jwtService.GetRefreshTokenExpiration()),
	})
	return nil
}

func (j JWTUseCase) ClearJWTCookies(ctx fiber.Ctx) error {
	ctx.ClearCookie("access-token")
	ctx.ClearCookie("refresh-token")

	ctx.SendString("logout")
	return nil
}

func (j JWTUseCase) GetAccessToken(ctx fiber.Ctx) (*jwt.Token, error) {
	access_token := string(ctx.Request().Header.Cookie("access-token"))
	return j.jwtService.ValidateToken(access_token)
}

func (j JWTUseCase) GetRefreshToken(ctx fiber.Ctx) (*jwt.Token, error) {
	refresh_token := string(ctx.Request().Header.Cookie("refresh-token"))
	return j.jwtService.ValidateToken(refresh_token)
}

func (j JWTUseCase) GetTokenUserId(token *jwt.Token) (int, error) {

	user_id := int(token.Claims.(jwt.MapClaims)["id"].(float64))
	return user_id, nil
}

func (j JWTUseCase) GetTokenVersion(token *jwt.Token) (int, error) {
	version := int(token.Claims.(jwt.MapClaims)["ver"].(float64))
	return version, nil
}
