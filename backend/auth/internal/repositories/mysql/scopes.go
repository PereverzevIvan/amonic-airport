package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

func ScopePaginate(params *models.UserSessionParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if params == nil {
			return db
		}

		return db.
			Offset((params.Page - 1) * params.Limit).
			Limit(params.Limit)
	}
}
