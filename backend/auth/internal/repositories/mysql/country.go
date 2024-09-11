package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type CountryRepo struct {
	Conn *gorm.DB
}

func NewCountryRepo(conn *gorm.DB) CountryRepo {
	return CountryRepo{
		Conn: conn,
	}
}

func (r CountryRepo) GetByID(id int) (*models.Country, error) {
	var country models.Country

	err := r.Conn.First(&country, id).Error
	if err != nil {
		return nil, err
	}

	return &country, nil
}

func (r CountryRepo) GetByName(name string) (*models.Country, error) {
	var country models.Country

	err := r.Conn.First(&country, "name = ?", name).Error
	if err != nil {
		return nil, err
	}

	return &country, nil
}
