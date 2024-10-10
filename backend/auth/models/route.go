package models

type Route struct {
	ID int `json:"id" gorm:"Column:ID"`

	DepartureAirportID int      `json:"departure_airport_id" gorm:"Column:DepartureAirportID"`
	DepartureAirport   *Airport `json:"departure_airport" gorm:"foreignKey:DepartureAirportID"`

	ArrivalAirportID int      `json:"arrival_airport_id" gorm:"Column:ArrivalAirportID"`
	ArrivalAirport   *Airport `json:"arrival_airport" gorm:"foreignKey:ArrivalAirportID"`

	Distance   int `json:"distance" gorm:"Column:Distance"`
	FlightTime int `json:"flight_time" gorm:"Column:FlightTime"`
}

func (Route) TableName() string {
	return "routes"
}
