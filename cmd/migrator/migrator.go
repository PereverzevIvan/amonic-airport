package main

import (
	"database/sql"
	"fmt"

	"gitflic.ru/project/pereverzevivan/jwt-auth-golang/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoadConfig().ConfigDatabase
	var dsn string

	switch cfg.DBType {
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
	case "mysql":
		dsn = fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/%v?multiStatements=true",
			cfg.AdminName,
			cfg.DBPassword,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
		)
	default:
		panic("Неверный тип базы данных")
	}

	db, err := sql.Open(cfg.DBType, dsn)

	if err != nil {
		panic(err.Error())
	}

	var driver database.Driver

	switch cfg.DBType {
	case "mysql":
		driver, err = mysql.WithInstance(db, &mysql.Config{})
	case "postgres":
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	default:
		panic("Неверный тип базы данных")

	}

	if err != nil {
		panic(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./storage/migrations",
		"postgres", driver)

	// m.Down()
	err = m.Up()
	if err != nil {
		fmt.Println(err.Error())
	}
	// m.Steps()
}
