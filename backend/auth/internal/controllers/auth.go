package controllers

import (
	"net/http"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type AuthController struct {
	jwtUseCase         jwtUseCase
	userService        userService
	userSessionUseCase userSessionUseCase
}

type LoginRequest struct {
	Email    string `json:"email" example:"j.doe@amonic.com"`
	Password string `json:"password" example:"123"`
}

func AddAuthControllerRoutes(
	api *fiber.Router,
	jwtUseCase jwtUseCase,
	userService userService,
	userSessionUseCase userSessionUseCase,
	authMiddleware AuthMiddleware,
) {
	controller := &AuthController{
		jwtUseCase:         jwtUseCase,
		userService:        userService,
		userSessionUseCase: userSessionUseCase,
	}

	(*api).Post("login/", controller.Login)
	// (*api).Get("register/", controller.Register, authMiddleware.IsAdmin)
	(*api).Get("logout/", controller.Logout)
	(*api).Get("refresh/", controller.Refresh)
}

// Вход в систему
// 1. Поиск пользователя по email
// 2. Проверка пароля
// 3. Получить версию токенов и увеличить на 1
// 4. Запись новых токенов в куки
// 5. Проверить корректность предыдущего выхода из системы
//  1. Если выхода не было -> изменить запись добавив причину 'Не указано'
// 6. Записать новую сессию

// @Summary      User login
// @Description  Вход пользователя по адресу электронной почты и паролю.
// @Description  Возвращает два токена для авторизации в куках.
// @Description  Название кук: "access-token" и "refresh-token"
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login params"
// @Success      200 "login success"
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /login [post]
func (ac *AuthController) Login(ctx fiber.Ctx) error {
	body := map[string]string{}
	err := ctx.Bind().Body(&body)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString("body is required")
	}

	// получение email и пароля
	email, hasEmail := body["email"]
	if !hasEmail || (email == "") {
		return ctx.Status(http.StatusBadRequest).SendString("email is required")
	}

	password, hasPassword := body["password"]
	if !hasPassword || (password == "") {
		return ctx.Status(http.StatusBadRequest).SendString("password is required")
	}

	// Получение пользователя и проверка пароля
	user, err := ac.userService.GetByEmail(email)
	if err != nil {
		log.Error(err)
		return ctx.Status(http.StatusInternalServerError).SendString("Wrong email or password")
	}

	if user == nil || !ac.userService.IsPasswordCorrect(user, password) {
		ctx.SendStatus(http.StatusUnauthorized)
		return ctx.SendString("Wrong email or password")
	}

	// Проверить, что текущий пользователь активен
	if !user.Active {
		ctx.SendStatus(http.StatusForbidden)
		return ctx.SendString("User deactivated")
	}

	// Обновление токенов
	err = ac.jwtUseCase.RefreshTokens(ctx, user)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}

	// Проверить последнюю сессию, если не было выхода
	// и crashReasonType не установлен -> обновить эту невалидную сессию,
	// установив crashReason = 0 и KDefaultInvalidSessionReason
	err = ac.userSessionUseCase.UpdateNoLogoutSession(ctx, user.ID)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}
	// Запись сессии
	err = ac.userSessionUseCase.CreateNewLoginSession(ctx, user.ID)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}

	ctx.SendStatus(http.StatusOK)
	return ctx.SendString("login success")
}

// @Summary      Refresh token
// @Description  Получить новую пару
// @Tags         Auth
// @Success      200 "refresh success"
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /refresh [get]
func (ac *AuthController) Refresh(ctx fiber.Ctx) error {
	user_id, err := ac.jwtUseCase.GetUserIdFromToken(ctx, true)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}

	user, err := ac.userService.GetByID(user_id)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	// Обновляем токены
	err = ac.jwtUseCase.RefreshTokens(ctx, user)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}

	return ctx.SendStatus(http.StatusOK)
}

// @Summary      User Logout
// @Description  Выход пользователя из системы
// @Tags         Auth
// @Success      200 "logout success"
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /logout [get]
func (ac *AuthController) Logout(ctx fiber.Ctx) error {
	user_id, err := ac.jwtUseCase.GetUserIdFromToken(ctx, true)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
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

	// Помечаем сессию
	err = ac.userSessionUseCase.LogoutLastSession(ctx, user.ID)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}
	return ctx.SendString("logout")
}
