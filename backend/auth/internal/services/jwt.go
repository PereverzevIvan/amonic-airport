package service

import (
	"fmt"
	"time"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/config"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3/log"
	"github.com/golang-jwt/jwt"
)

const KInvalidTokenVersion string = "invalid token version"

type TokensVersionRepo interface {
	GetUserTokensVersion(user_id int) (int, error)
	IncreaseUserTokensVersion(user_id int) error
}

type jwtService struct {
	config            *config.ConfigJWT
	tokensVersionRepo TokensVersionRepo
}

func NewJWTService(jwtConfig *config.ConfigJWT, tokensVersionRepo TokensVersionRepo) jwtService {
	return jwtService{
		config:            jwtConfig,
		tokensVersionRepo: tokensVersionRepo,
	}
}

func (j jwtService) generateToken(user *models.User, expiration time.Duration, version int) (string, error) {
	role := "user"
	if user.RoleID == int(models.KRoleAdmin) {
		role = "admin"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"role":  role,
		"id":    user.ID,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(expiration).Unix(),
		"ver":   version,
		// "nbf":   user.NotBefore,
		// "jti":   user.JTI,
	})

	return token.SignedString([]byte(j.config.SecretKey))
}

func (j jwtService) GenerateAccessAndRefreshTokens(user *models.User) (string, string, error) {
	// Увеличить и получить версию токенов
	err := j.tokensVersionRepo.IncreaseUserTokensVersion(user.ID)
	if err != nil {
		return "", "", err
	}

	userTokensVersion, err := j.tokensVersionRepo.GetUserTokensVersion(user.ID)
	if err != nil {
		return "", "", err
	}

	// Генерация токенов
	access_token, err := j.generateToken(user, j.config.AccessTokenExpiration, userTokensVersion)
	if err != nil {
		return "", "", err
	}

	refresh_token, err := j.generateToken(user, j.config.RefreshTokenExpiration, userTokensVersion)
	if err != nil {
		return "", "", err
	}

	return access_token, refresh_token, nil
}

func (j jwtService) ValidateToken(token string) (*jwt.Token, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	jwt_token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Проверка токена
	user_id := int(jwt_token.Claims.(jwt.MapClaims)["id"].(float64))
	db_token_version, err := j.tokensVersionRepo.GetUserTokensVersion(user_id)
	if err != nil {
		return nil, err
	}

	jwt_token_version := int(jwt_token.Claims.(jwt.MapClaims)["ver"].(float64))

	// Проверка версии токена
	if jwt_token_version != db_token_version {
		log.Error(fmt.Sprintf("invalid token version: token %v != db %v", jwt_token_version, db_token_version))
		return nil, fmt.Errorf(KInvalidTokenVersion)
	}

	if !jwt_token.Valid {
		return nil, fmt.Errorf("jwt_token.Valid is false")
	}

	return jwt_token, nil
}

func (j jwtService) GetAccessTokenExpiration() time.Duration {
	return j.config.AccessTokenExpiration
}

func (j jwtService) GetRefreshTokenExpiration() time.Duration {
	return j.config.RefreshTokenExpiration
}
