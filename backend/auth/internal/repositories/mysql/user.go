package mysql_repo

import "gorm.io/gorm"

type UserRepo struct {
	Conn *gorm.DB
}

func NewUserRepo(conn *gorm.DB) UserRepo {
	return UserRepo{
		Conn: conn,
	}
}
