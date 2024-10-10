package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type RouteRepo struct {
	Conn *gorm.DB
}

func NewRouteRepo(conn *gorm.DB) RouteRepo {
	return RouteRepo{Conn: conn}
}

func (r RouteRepo) GetIDByAirportCodes(departure_airport_code, arrival_airport_code string) (int, error) {

	var route_id int
	err := r.Conn.
		Model(&models.Route{}).
		Select("routes.ID").
		Where(`
			routes.DepartureAirportID IN (
				SELECT airports.ID FROM airports 
				WHERE airports.IATACode = ?
			)`, departure_airport_code).
		Where(` 
			routes.ArrivalAirportID IN (
				SELECT airports.ID FROM airports 
				WHERE airports.IATACode = ?
			)`, arrival_airport_code).
		First(&route_id).
		Error

	if err == gorm.ErrRecordNotFound {
		return 0, models.ErrNotFound
	}

	return route_id, err
}
