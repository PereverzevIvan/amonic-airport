package controllers

import (
	"fmt"
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
// @Summary      Get user sessions
// @Description  Получить список сессий пользователя по его ID
// @Tags         User sessions
// @Accept       json
// @Produce      json
// @Param        id path  int  true  "User ID"
// @Param        only_invalid_sessions query  bool  false  "Get only invalid sessions of user"
// @Param        page query  int  false  "Page number"
// @Param        limit query  int  false  "Limit of records in page"
// @Success      200  {object}  []models.UserSession
// @Failure      400  {object}  error
// @Failure      403  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /user-sessions/ [get]
func (controller *UserSessionController) UserSessions(ctx fiber.Ctx) error {
	// Получаем user_id и затем пользователя
	token_user_id, err := controller.jwtUseCase.GetUserIdFromToken(ctx, false)
	if err != nil {
		log.Error(err)
		return nil
	}

	var params models.UserSessionParams

	if err := ctx.Bind().Query(&params); err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}
	if params.UserID == 0 {
		params.UserID = token_user_id
	}

	// Если user_id не равен token_user_id,
	// то проверяем, является ли пользователь админом
	if token_user_id != params.UserID {
		ctx.SendString("Вы не можете получить доступ к чужим сессиям")
		return ctx.SendStatus(http.StatusForbidden)

		// is_admin, err := controller.userService.IsAdmin(user_id)
		// if err != nil {
		// 	log.Error(err)
		// 	return ctx.SendStatus(http.StatusInternalServerError)
		// }

		// if !is_admin {
		// 	return ctx.SendStatus(http.StatusForbidden)
		// }
	}

	user_sessions, err := controller.userSessionUseCase.GetByUserID(ctx, params.UserID, &params)
	if err != nil {
		return utils.LogErrorIfNotEmpty(err)
	}

	ctx.SendStatus(http.StatusOK)
	return ctx.JSON(user_sessions)
}

// Установка причины неудачного выхода из системы
// @Summary      Set Unsuccessfull Logout Reason
// @Description  Установка причины неудачного выхода из системы
// @Tags         User sessions
// @Accept       json
// @Produce      json
// @Param        session_data body models.UserSession false "Информация о сессии"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /user-sessions/ [patch]
func (controller *UserSessionController) SetUnsuccessfullLogoutReason(ctx fiber.Ctx) error {
	// Получаем user_id и затем пользователя
	user_id, err := controller.jwtUseCase.GetUserIdFromToken(ctx, false)
	if err != nil {
		log.Error(err)
		return nil
	}

	body := map[string]string{}
	err = ctx.Bind().Body(&body)
	fmt.Println(body)
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
