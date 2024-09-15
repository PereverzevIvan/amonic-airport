package service

import (
	"time"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	// "github.com/gofiber/fiber/v3/log"
)

type UserSessionRepo interface {
	GetByID(id int) (*models.UserSession, error)
	GetByUserId(user_id int, params *models.UserSessionParams) (*[]models.UserSession, error)
	GetLastByUserId(user_id int) (*models.UserSession, error)
	Create(user_session *models.UserSession) error
	Update(user_session *models.UserSession) error
	// Delete(user_id int) error
}

type UserSessionService struct {
	userSessionRepo UserSessionRepo
}

func NewUserSessionService(ur UserSessionRepo) UserSessionService {
	return UserSessionService{ur}
}

func (us UserSessionService) GetByUserID(user_id int, params *models.UserSessionParams) (*[]models.UserSession, error) {
	user_sessions, err := us.userSessionRepo.GetByUserId(user_id, params)
	return user_sessions, err
}
func (us UserSessionService) GetByID(id int) (*models.UserSession, error) {
	user_session, err := us.userSessionRepo.GetByID(id)
	return user_session, err
}

func (us UserSessionService) GetLastByUserId(user_id int) (*models.UserSession, error) {
	user_session, err := us.userSessionRepo.GetLastByUserId(user_id)
	if err != nil {
		return nil, err
	}

	return user_session, err
}
func (us UserSessionService) Create(user_id int) (*models.UserSession, error) {
	user_session := models.UserSession{
		UserID:          user_id,
		LoginAt:         time.Now().UTC(),
		CrashReasonType: nil,
	}

	err := us.userSessionRepo.Create(&user_session)
	if err != nil {
		return nil, err
	}
	return &user_session, err
}
func (us UserSessionService) Update(user_session *models.UserSession) error {
	err := us.userSessionRepo.Update(user_session)
	return err
}
