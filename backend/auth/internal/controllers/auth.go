package controllers

import (
	"net/http"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/golang-jwt/jwt"
)

type jwtUseCase interface {
	RefreshTokens(ctx fiber.Ctx, user *models.User) error
	GetAccessToken(ctx fiber.Ctx) (*jwt.Token, error)
	GetRefreshToken(ctx fiber.Ctx) (*jwt.Token, error)
	ClearJWTCookies(ctx fiber.Ctx) error
	GetTokenUserId(token *jwt.Token) (int, error)
}

type userService interface {
	GetByID(user_id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	IsPasswordCorrect(user *models.User, password string) bool
}

type AuthMiddleware interface {
	IsActive(ctx fiber.Ctx) error
	IsAdmin(ctx fiber.Ctx) error
}

type AuthController struct {
	jwtUseCase  jwtUseCase
	userService userService
}

func AddAuthControllerRoutes(api *fiber.Router, jwtUseCase jwtUseCase, userService userService, authMiddleware AuthMiddleware) {
	controller := &AuthController{
		jwtUseCase:  jwtUseCase,
		userService: userService,
	}

	(*api).Post("login/", controller.Login)
	// (*api).Get("register/", controller.Register, authMiddleware.IsAdmin)
	(*api).Get("logout/", controller.Logout)
	(*api).Get("refresh/", controller.Refresh, authMiddleware.IsActive)
}

// Вход в систему
// 1. Поиск пользователя по email
// 2. Проверка пароля
// 3. Получить версию токенов и увеличить на 1
// 4. Запись новых токенов в куки
// 5. Проверить корректность предыдущего выхода из системы
//  1. Если выхода не было -> изменить запись добавив причину 'Не указано'
//
// 6. Записать новую сессию
func (ac *AuthController) Login(ctx fiber.Ctx) error {
	body := map[string]string{}
	err := ctx.Bind().Body(&body)
	if err != nil {
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString("body is required")
	}

	// получение email и пароля
	email, hasEmail := body["email"]
	if !hasEmail || (email == "") {
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString("email is required")
	}

	password, hasPassword := body["password"]
	if !hasPassword || (password == "") {
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString("password is required")
	}

	// Получение пользователя и проверка пароля
	user, err := ac.userService.GetByEmail(email)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if user == nil {
		ctx.SendStatus(http.StatusUnauthorized)
		return ctx.SendString("Wrong email or password")
	}

	if !ac.userService.IsPasswordCorrect(user, password) {
		ctx.SendStatus(http.StatusUnauthorized)
		return ctx.SendString("Wrong email or password")
	}

	// Обновление токенов
	err = ac.jwtUseCase.RefreshTokens(ctx, user)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	ctx.SendStatus(http.StatusOK)
	return ctx.SendString("login success")
}

func (ac *AuthController) Refresh(ctx fiber.Ctx) error {
	// Получаем resfresh токен
	refresh_token, err := ac.jwtUseCase.GetRefreshToken(ctx)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	// Получаем user_id и затем пользователя
	user_id, err := ac.jwtUseCase.GetTokenUserId(refresh_token)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	user, err := ac.userService.GetByID(user_id)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	// Обновляем токены
	err = ac.jwtUseCase.RefreshTokens(ctx, user)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}

func (ac *AuthController) Logout(ctx fiber.Ctx) error {
	// Получаем resfresh токен
	refresh_token, err := ac.jwtUseCase.GetRefreshToken(ctx)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	// Получаем user_id и затем пользователя
	user_id, err := ac.jwtUseCase.GetTokenUserId(refresh_token)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	user, err := ac.userService.GetByID(user_id)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	// Обновляем токены чтобы увеличить версию и не подключать другие сервисы
	ac.jwtUseCase.RefreshTokens(ctx, user)

	ctx.ClearCookie("access-token")
	ctx.ClearCookie("refresh-token")

	return ctx.SendString("logout")
}
