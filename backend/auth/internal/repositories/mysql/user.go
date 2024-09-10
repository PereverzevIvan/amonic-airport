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

func (u UserRepo) GetByID(user_id int) (models.User, error) {
	var user models.User
	err := u.Conn.First(&user, user_id).Error
	if err != nil {
		return models.User{}, err
	}
	return user, err
}
