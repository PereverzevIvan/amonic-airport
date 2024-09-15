package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type UserSessionRepo struct {
	Conn *gorm.DB
}

func NewUserSessionRepo(conn *gorm.DB) UserSessionRepo {
	return UserSessionRepo{
		Conn: conn,
	}
}

func (repo UserSessionRepo) GetByID(id int) (*models.UserSession, error) {
	user_session := models.UserSession{}
	res := repo.Conn.Model(&models.UserSession{}).Where("id = ?", id).First(&user_session)
	return &user_session, res.Error
}

func (repo UserSessionRepo) GetByUserId(user_id int, params *models.UserSessionParams) (*[]models.UserSession, error) {
	var user_sessions []models.UserSession

	res := repo.Conn.
		Scopes(
			ScopePaginate(params),
			ScopeUserSessionParams(params),
		).
		Model(&models.UserSession{}).
		Where("user_id = ?", user_id).
		Order("login_at desc").
		Find(&user_sessions)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user_sessions, res.Error
}

func ScopeUserSessionParams(params *models.UserSessionParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if params == nil {
			return db
		}

		res := db.Where("user_id = ?", params.UserID)

		if params.OnlyInvalidSessions {
			res = res.Where("invalid_logout_reason = ?", models.KInvalidLogoutReasonUndefined)
		}

		return res
	}
}

func (repo UserSessionRepo) GetLastByUserId(user_id int) (*models.UserSession, error) {
	var user_session models.UserSession

	res := repo.Conn.
		Model(&models.UserSession{}).
		Order("login_at desc").
		Find(&user_session)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &user_session, res.Error
}

func (repo UserSessionRepo) Create(userSession *models.UserSession) error {
	res := repo.Conn.Create(userSession)
	return res.Error
}

func (repo UserSessionRepo) Update(userSession *models.UserSession) error {
	res := repo.Conn.Save(&userSession)
	return res.Error
}
