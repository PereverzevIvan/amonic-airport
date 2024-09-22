package usecases

import (
	"fmt"
	"net/http"
	"time"

	service "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/services"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/utils"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v2/log"
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
func (usecase JWTUseCase) RefreshTokens(ctx fiber.Ctx, user *models.User) error {

	access_token, refresh_token, err := usecase.jwtService.GenerateAccessAndRefreshTokens(user)
	if err != nil {
		ctx.SendStatus(http.StatusInternalServerError)
		ctx.SendString(err.Error())
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access-token",
		Value:    access_token,
		Secure:   false, // В продакшене должно быть true
		SameSite: "Lax",
		Expires:  time.Now().Add(usecase.jwtService.GetAccessTokenExpiration()),
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh-token",
		Value:    refresh_token,
		HTTPOnly: true,
		Secure:   false, // В продакшене должно быть true
		SameSite: "Lax",
		Expires:  time.Now().Add(usecase.jwtService.GetRefreshTokenExpiration()),
	})
	return nil
}

func (usecase JWTUseCase) ClearJWTCookies(ctx fiber.Ctx) error {
	ctx.ClearCookie("access-token")
	ctx.ClearCookie("refresh-token")

	ctx.SendString("logout")
	return nil
}

// func (usecase JWTUseCase) getRefreshToken(ctx fiber.Ctx) (*jwt.Token, error) {
// 	refresh_token_str := string(ctx.Request().Header.Cookie("refresh-token"))
// 	log.Info(refresh_token_str)
// 	if refresh_token_str == "" {
// 		ctx.SendStatus(http.StatusUnauthorized)
// 		return nil, fmt.Errorf("")
// 	}

// 	refresh_token, err := usecase.jwtService.ValidateToken(refresh_token_str)
// 	if err != nil {
// 		log.Error(err)
// 		ctx.SendStatus(http.StatusInternalServerError)
// 		return nil, fmt.Errorf("")
// 	}

// 	return refresh_token, nil
// }

func (usecase JWTUseCase) GetTokenUserId(token *jwt.Token) (int, error) {
	user_id := int(token.Claims.(jwt.MapClaims)["id"].(float64))
	return user_id, nil
}

// func (usecase JWTUseCase) GetTokenVersion(token *jwt.Token) (int, error) {
// 	version := int(token.Claims.(jwt.MapClaims)["ver"].(float64))
// 	return version, nil
// }

func (usecase JWTUseCase) getToken(ctx fiber.Ctx, use_refresh_token bool) (*jwt.Token, error) {
	var jwt_token_str string

	if use_refresh_token {
		jwt_token_str = string(ctx.Request().Header.Cookie("refresh-token"))
	} else {
		jwt_token_str = string(ctx.Request().Header.Cookie("access-token"))
	}

	if jwt_token_str == "" {
		ctx.SendStatus(http.StatusUnauthorized)
		return nil, fmt.Errorf("")
	}

	jwt_token, err := usecase.jwtService.ValidateToken(jwt_token_str)
	if err != nil {
		if err.Error() == service.KInvalidTokenVersion {
			ctx.SendStatus(http.StatusUnauthorized)
			ctx.SendString(service.KInvalidTokenVersion)
			return nil, fmt.Errorf("")
		}

		utils.LogErrorIfNotEmpty(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return nil, fmt.Errorf("")
	}

	return jwt_token, nil
}

func (usecase JWTUseCase) GetUserIdFromToken(ctx fiber.Ctx, use_refresh_token bool) (int, error) {
	token, err := usecase.getToken(ctx, use_refresh_token)
	if err != nil || token == nil {
		return -1, fmt.Errorf("")
	}

	user_id, err := usecase.GetTokenUserId(token)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return -1, fmt.Errorf("")
	}
	return user_id, err
}

// func (usecase JWTUseCase) GetVersionFromToken(ctx fiber.Ctx, use_refresh_token bool) (int, error) {
// 	token, err := usecase.getToken(ctx, use_refresh_token)
// 	if err != nil || token == nil {
// 		return -1, fmt.Errorf("")
// 	}

// 	version := int(token.Claims.(jwt.MapClaims)["ver"].(float64))
// 	return version, nil
// }
