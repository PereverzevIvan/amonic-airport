package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type OfficeRepo struct {
	Conn *gorm.DB
}

func NewOfficeRepo(conn *gorm.DB) OfficeRepo {
	return OfficeRepo{Conn: conn}
}

func (r OfficeRepo) GetByID(id int) (*models.Office, error) {
	var office models.Office
	err := r.Conn.First(&office, id).Error
	if err != nil {
		return nil, err
	}

	return &office, nil
}

func (r OfficeRepo) GetByTitle(title string) (*models.Office, error) {
	var office models.Office
	err := r.Conn.First(&office, "title = ?", title).Error
	if err != nil {
		return nil, err
	}

	return &office, nil
}
