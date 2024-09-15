package usecases

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v3"
)

type userSessionService interface {
	GetByUserID(user_id int, params *models.UserSessionParams) (*[]models.UserSession, error)
	GetByID(id int) (*models.UserSession, error)
	GetLastByUserId(user_id int) (*models.UserSession, error)
	Create(user_id int) (*models.UserSession, error)
	Update(user_session *models.UserSession) error
}

type UserSessionUseCase struct {
	userSessionService userSessionService
}

func NewUserSessionUseCase(UserSessionService userSessionService) UserSessionUseCase {
	return UserSessionUseCase{
		userSessionService: UserSessionService,
	}
}

func (usecase UserSessionUseCase) GetByUserID(ctx fiber.Ctx, user_id int, params *models.UserSessionParams) (*[]models.UserSession, error) {
	user_sessions, err := usecase.userSessionService.GetByUserID(user_id, params)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return nil, err
	}
	return user_sessions, nil
}

func (usecase UserSessionUseCase) GetByID(ctx fiber.Ctx, session_id int) (*models.UserSession, error) {
	user_session, err := usecase.userSessionService.GetByID(session_id)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return nil, fmt.Errorf("")
	}
	return user_session, nil
}

func (usecase UserSessionUseCase) GetLastByUserId(ctx fiber.Ctx, user_id int) (*models.UserSession, error) {
	user_session, err := usecase.userSessionService.GetLastByUserId(user_id)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return nil, fmt.Errorf("")
	}
	return user_session, nil
}

func (usecase UserSessionUseCase) Create(ctx fiber.Ctx, user_id int) (*models.UserSession, error) {
	user_session, err := usecase.userSessionService.Create(user_id)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return nil, fmt.Errorf("")
	}
	return user_session, nil
}

func (usecase UserSessionUseCase) Update(ctx fiber.Ctx, user_session *models.UserSession) error {
	err := usecase.userSessionService.Update(user_session)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return fmt.Errorf("")
	}
	return nil
}

// Bissness logic

func (usecase UserSessionUseCase) ObtainUserSessionFromBody(ctx fiber.Ctx, user_id int) (*models.UserSession, error) {
	body := map[string]string{}
	err := ctx.Bind().Body(&body)
	if err != nil {
		ctx.SendStatus(http.StatusBadRequest)
		return nil, fmt.Errorf("body is required")
	}

	// получение session_id
	session_id_str, hasUserSessionId := body["id"]
	if !hasUserSessionId || session_id_str == "" {
		ctx.SendStatus(http.StatusBadRequest)
		ctx.SendString("id for session is required")
		return nil, fmt.Errorf("")
	}

	session_id, err := strconv.Atoi(session_id_str)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		ctx.SendString("Can't parse session_id")
		return nil, fmt.Errorf("")
	}

	user_session, err := usecase.userSessionService.GetByID(session_id)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return nil, fmt.Errorf("")
	}

	if user_session == nil {
		ctx.SendStatus(http.StatusNotFound)
		ctx.SendString("Session not found")
		return nil, fmt.Errorf("")
	}

	if user_session.UserID != user_id {
		log.Error(user_session, user_id)
		ctx.SendStatus(http.StatusForbidden)
		ctx.SendString("You don't have access to this session")
		return nil, fmt.Errorf("")
	}

	if user_session.CrashReasonType == nil {
		ctx.SendStatus(http.StatusBadRequest)
		ctx.SendString("This session is correct, you can't change it")
		return nil, fmt.Errorf("")
	}

	// получение reason
	reason, hasReason := body["reason"]
	if !hasReason || reason == "" {
		ctx.SendStatus(http.StatusBadRequest)
		ctx.SendString("reason is required")
		return nil, fmt.Errorf("")
	}
	if reason == models.KInvalidLogoutReasonUndefined {
		ctx.SendStatus(http.StatusBadRequest)
		ctx.SendString("Invalid logout reason. Reason can't be '" + models.KInvalidLogoutReasonUndefined + "'")
		return nil, fmt.Errorf("")
	}
	user_session.InvalidLogoutReason = reason

	// получение session_id
	crash_reason_type_str, hasCrashReasonType := body["crash_reason_type"]
	if !hasCrashReasonType {
		ctx.SendStatus(http.StatusBadRequest)
		ctx.SendString("crash_reason_type is required")
		return nil, fmt.Errorf("")
	}

	crash_reason_type, err := strconv.Atoi(crash_reason_type_str)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		ctx.SendString("Can't parse crash_reason_type")
		return nil, fmt.Errorf("")
	}
	user_session.CrashReasonType = &crash_reason_type

	return user_session, nil
}

func (usecase UserSessionUseCase) UpdateNoLogoutSession(ctx fiber.Ctx, user_id int) error {
	last_session, err := usecase.userSessionService.GetLastByUserId(user_id)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return fmt.Errorf("")
	}

	// Если нет последнего сеанса или он завершен, то ничего не делаем
	if last_session == nil || last_session.LogoutAt != nil {
		return nil
	}

	// set undefined crash type and reason
	last_session.InvalidLogoutReason = models.KInvalidLogoutReasonUndefined

	undefined_crash_reason_type := int(models.KCrashReasonUndefined)
	last_session.CrashReasonType = &undefined_crash_reason_type

	err = usecase.userSessionService.Update(last_session)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return fmt.Errorf("")
	}
	return nil
}

func (usecase UserSessionUseCase) CreateNewLoginSession(ctx fiber.Ctx, user_id int) error {
	_, err := usecase.userSessionService.Create(user_id)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return fmt.Errorf("")
	}
	return nil
}

func (usecase UserSessionUseCase) LogoutLastSession(ctx fiber.Ctx, user_id int) error {
	last_session, err := usecase.userSessionService.GetLastByUserId(user_id)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		return fmt.Errorf("")
	}

	if last_session == nil {
		log.Error("Last session not found, user_id: ", user_id)
		ctx.SendStatus(http.StatusInternalServerError)
		ctx.SendString("Last session not found")
		return fmt.Errorf("")
	}

	if last_session.LogoutAt != nil {
		log.Error("Last session already logged out, user_id: ", user_id, ", session_id: ", last_session.ID)
		ctx.SendStatus(http.StatusInternalServerError)
		ctx.SendString("Last session already logged out")
		return fmt.Errorf("")
	}

	var now time.Time = time.Now()
	last_session.LogoutAt = &now

	err = usecase.userSessionService.Update(last_session)
	if err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusInternalServerError)
		ctx.SendString("Can't update last session")
		return fmt.Errorf("")
	}

	return nil
}
