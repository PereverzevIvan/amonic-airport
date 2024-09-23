package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	Conn *gorm.DB
}

func NewUserRepo(conn *gorm.DB) UserRepo {
	return UserRepo{
		Conn: conn,
	}
}

func (r UserRepo) Create(user *models.User) error {
	err := r.Conn.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (u UserRepo) GetByID(user_id int) (*models.User, error) {
	var user models.User
	err := u.Conn.First(&user, user_id).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (u UserRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := u.Conn.First(&user, "Email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (u UserRepo) GetAll(params map[string]string) (*[]models.User, error) {
	var users []models.User
	// var totalCount int

	query := u.Conn.Model(&models.User{})

	if params["office_id"] != "" {
		query = query.Where("OfficeID = ?", params["office_id"])
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (u UserRepo) IsAdmin(user_id int) (bool, error) {
	var role_id models.ERole = models.KRoleNone
	res := u.Conn.Model(&models.User{}).
		Select("RoleID").
		Where("id = ?", user_id).
		Scan(&role_id)

	is_admin := role_id == models.KRoleAdmin
	return is_admin, res.Error
}

func (u UserRepo) IsActive(user_id int) (bool, error) {
	is_active := false
	res := u.Conn.Model(&models.User{}).
		Select("Active").
		Where("id = ?", user_id).
		Scan(&is_active)

	return is_active, res.Error
}
