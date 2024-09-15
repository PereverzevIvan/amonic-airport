package mysql_repo

import (
	"errors"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserExists   = errors.New("Пользователь уже существует")
	ErrUserNotFound = errors.New("Пользователь не найден")
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
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				return ErrUserExists
			}
		}

		return err
	}

	return nil
}

func (u UserRepo) GetByID(user_id int) (*models.User, error) {
	var user models.User
	err := u.Conn.First(&user, user_id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, err
}

func (u UserRepo) IsAdmin(user_id int) (bool, error) {
	var role_id models.ERole
	res := u.Conn.Model(&models.User{}).
		Select("RoleID").
		Where("id = ?", user_id).
		Scan(&role_id)

	is_admin := role_id == models.KRoleAdmin
	return is_admin, res.Error
}
