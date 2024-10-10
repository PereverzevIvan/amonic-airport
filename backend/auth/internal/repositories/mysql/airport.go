package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type AirportRepo struct {
	Conn *gorm.DB
}

func NewAirportRepo(conn *gorm.DB) AirportRepo {
	return AirportRepo{
		Conn: conn,
	}
}

func (r AirportRepo) GetAll() (*[]models.Airport, error) {
	var airports []models.Airport

	err := r.Conn.
		Model(&models.Airport{}).
		Find(&airports).
		Error

	if err != nil {
		return nil, err
	}

	return &airports, nil
}
