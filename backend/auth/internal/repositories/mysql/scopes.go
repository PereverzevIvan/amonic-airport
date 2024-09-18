package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

func ScopePaginate(params *models.PaginateParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if params == nil {
			return db
		}

		limit := min(max(params.Limit, 3), models.KMaxLimit)
		page := max(params.Page-1, 0)

		return db.
			Offset(page * limit).
			Limit(limit)
	}
}
