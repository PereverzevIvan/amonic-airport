package mysql_repo

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v3/log"
)

// Unique constraint violation error handling
func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) { // Check if the error is of type *mysql.MySQLError
		return mysqlErr.Number == 1062 // MySQL error code for duplicate entry is 1062
	}

	return false
}

func IsForeignKeyConstraintError(err error) bool {
	if err == nil {
		return false
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) { // Check if the error is of type *mysql.MySQLError
		return mysqlErr.Number == 1452
	}

	return false
}

func IsRecordNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) { // Check if the error is of type *mysql.MySQLError
		log.Info(mysqlErr.Number)
	}

	return err.Error() == "record not found"
}
