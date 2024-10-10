package models

import (
	"time"
)

type Schedule struct {
	ID int `json:"id" gorm:"Column:ID"`

	AircraftID int       `json:"aircraft_id" gorm:"Column:AircraftID"`
	Aircraft   *Aircraft `json:"aircraft" gorm:"foreignKey:AircraftID"`

	RouteID int    `json:"route_id" gorm:"Column:RouteID"`
	Route   *Route `json:"route" gorm:"foreignKey:RouteID"`

	FlightNumber string  `json:"flight_number" gorm:"Column:FlightNumber"`
	EconomyPrice float64 `json:"economy_price" gorm:"Column:EconomyPrice"`
	Confirmed    bool    `json:"confirmed" gorm:"Column:Confirmed"`

	Outbound time.Time `json:"outbound" gorm:"column:Outbound;type:datetime;not null"`
	Date     string    `gorm:"column:Date;type:date;not null"`
	Time     string    `gorm:"column:Time;not null"`
}

func (Schedule) TableName() string {
	return "schedules"
}
