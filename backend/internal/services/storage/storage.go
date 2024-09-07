package storage

import (
	"fmt"

	"gitflic.ru/project/pereverzevivan/jwt-auth-golang/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	Conn *gorm.DB
}

func NewStorage(cfg config.ConfigDatabase) Storage {
	var db *gorm.DB
	var err error
	var dsn string

	switch cfg.DBType {
	case "mysql":
		dsn = fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.AdminName,
			cfg.DBPassword,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
		)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Не удалось подключиться к базе данных")
		}
	case "postgres":
		dsn = fmt.Sprintf(
			"host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
			cfg.Host,
			cfg.Port,
			cfg.AdminName,
			cfg.DBPassword,
			cfg.DBName,
			cfg.SSL,
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Не удалось подключиться к базе данных")
		}
	default:
		panic("Неверный тип базы данных")
	}

	return Storage{
		Conn: db,
	}
}
