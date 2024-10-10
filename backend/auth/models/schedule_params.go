package models

import (
	"fmt"
	"strconv"
	"time"
)

type SchedulesParams struct {
	Outbound     string `query:"outbound"`
	FlightNumber string `query:"flight_number"`

	DepartureAirportID *int `query:"from"`
	ArrivalAirportID   *int `query:"to"`

	SortBy *string `query:"sort_by"`
}

type KSchedulesSortBy string

const (
	KSchedulesSortByDateTime    KSchedulesSortBy = "date_time"
	KSchedulesSortByTicketPrice KSchedulesSortBy = "ticket_price"
	KSchedulesSortByConfirmed   KSchedulesSortBy = "confirmed"

	KSchedulesSortByDefault = KSchedulesSortByDateTime
)

var StringToKSchedulesSortByMap map[string]KSchedulesSortBy = map[string]KSchedulesSortBy{
	string(KSchedulesSortByDateTime):    KSchedulesSortByDateTime,
	string(KSchedulesSortByTicketPrice): KSchedulesSortByTicketPrice,
	string(KSchedulesSortByConfirmed):   KSchedulesSortByConfirmed,
}

func (s *SchedulesParams) Validate() error {
	if s == nil {
		return nil
	}

	if s.Outbound != "" {
		_, err := time.Parse("2006-01-02", s.Outbound)
		if err != nil {
			return fmt.Errorf("invalid outbound: %v", s.Outbound)
		}
	}

	if s.DepartureAirportID != nil && s.ArrivalAirportID != nil &&
		*s.DepartureAirportID == *s.ArrivalAirportID {
		return fmt.Errorf("departure_airport_id and arrival_airport_id must be different")
	}

	if s.SortBy != nil {
		if _, ok := StringToKSchedulesSortByMap[*s.SortBy]; !ok {
			return fmt.Errorf("invalid sort_by: %v", *s.SortBy)
		}
	}

	return nil
}

func (s *SchedulesParams) GetSortBy() KSchedulesSortBy {
	if s == nil || s.SortBy == nil {
		return KSchedulesSortByDefault
	}

	if sort_by, ok := StringToKSchedulesSortByMap[*s.SortBy]; ok {
		return sort_by
	}

	return KSchedulesSortByDefault
}

type ScheduleUpdateConfirmedParams struct {
	Confirmed *bool `json:"confirmed"`
}

func (s *ScheduleUpdateConfirmedParams) Validate() error {
	if s.Confirmed == nil {
		return fmt.Errorf("confirmed is required")
	}

	return nil
}

type ScheduleUpdateParams struct {
	Date         *string  `json:"date"`
	Time         *string  `json:"time"`
	EconomyPrice *float64 `json:"economy_price"`
}

func (s *ScheduleUpdateParams) Validate() error {
	if s.Date != nil {
		_, err := time.Parse("2006-01-02", *s.Date)
		if err != nil {
			return fmt.Errorf("invalid date: %v", *s.Date)
		}
	}

	if s.Time != nil {
		_, err := time.Parse("15:04", *s.Time)
		if err != nil {
			return fmt.Errorf("invalid time: %v", *s.Time)
		}
	}

	if s.EconomyPrice != nil && *s.EconomyPrice <= 0 {
		return fmt.Errorf("invalid economy_price: %v", *s.EconomyPrice)
	}

	return nil
}

type SchedulesUploadResult struct {
	TotalRowsCnt         int `json:"total_rows_cnt"`
	SuccessfulRowsCnt    int `json:"successful_rows_cnt"`
	FailedRowsCnt        int `json:"failed_rows_cnt"`
	MissingFieldsRowsCnt int `json:"missing_fields_rows_cnt"`
	DuplicatedRowsCnt    int `json:"duplicated_rows_cnt"`
}

type ScheduleAddEditCommand struct {
	IsAddCommand bool `json:"is_add_command"`

	RouteID int `json:"route_id"`

	DepartureAirportCode string `json:"from"`
	ArrivalAirportCode   string `json:"to"`

	AircraftID   int     `json:"aircraft_id"`
	FlightNumber string  `json:"flight_number"`
	EconomyPrice float64 `json:"cost"`
	Confirmed    bool    `json:"confirmed"`

	Outbound time.Time `json:"outbound"`
}

func (s *ScheduleAddEditCommand) ToStrHash() string {
	return fmt.Sprintf("%v;%v",
		s.Outbound.Format("2006-01-02"),
		s.FlightNumber,
	)
}

func (s *ScheduleAddEditCommand) ToSchedule() *Schedule {
	schedule := &Schedule{
		AircraftID:   s.AircraftID,
		RouteID:      s.RouteID,
		FlightNumber: s.FlightNumber,
		EconomyPrice: s.EconomyPrice,
		Confirmed:    s.Confirmed,
		Outbound:     s.Outbound.Add(-time.Hour * 3),
		Date:         s.Outbound.Format("2006-01-02"),
		Time:         s.Outbound.Format("15:04"),
	}

	return schedule
}

func (s *ScheduleAddEditCommand) Validate() error {
	// if s.DepartureAirportID == s.ArrivalAirportID {
	// 	return fmt.Errorf("departure_airport_id and arrival_airport_id must be different")
	// }

	if s.EconomyPrice <= 0 {
		return fmt.Errorf("cost must be greater than 0")
	}

	return nil
}

func ParseScheduleAddEditCommandFromCSVRecord(record []string) (*ScheduleAddEditCommand, error) {
	if len(record) != 9 {
		return nil, ErrCSVMissingFields
	}

	// Parse the record into ScheduleAddEditParams
	params := ScheduleAddEditCommand{}
	switch condition := record[0]; condition {
	case "ADD":
		params.IsAddCommand = true
	case "EDIT":
		params.IsAddCommand = false
	default:
		return nil, fmt.Errorf("invalid condition: %v", condition)
	}

	var err error
	// Parse datetime
	params.Outbound, err = time.Parse(
		"2006-01-02 15:04",
		record[1]+" "+record[2])
	if err != nil {
		return nil, fmt.Errorf("invalid datetime: %v", record[1])
	}

	params.FlightNumber = record[3]
	if params.FlightNumber == "" {
		return nil, fmt.Errorf("flight_number is required")
	}

	params.DepartureAirportCode = record[4]
	params.ArrivalAirportCode = record[5]

	params.AircraftID, err = strconv.Atoi(record[6])
	if err != nil || params.AircraftID <= 0 {
		return nil, fmt.Errorf("invalid aircraft_id: %v", record[6])
	}

	params.EconomyPrice, err = strconv.ParseFloat(record[7], 64)
	if err != nil || params.EconomyPrice <= 0 {
		return nil, fmt.Errorf("invalid cost: %v", record[7])
	}

	switch record[8] {
	case "OK":
		params.Confirmed = true
	case "CANCELLED":
		params.Confirmed = false
	default:
		return nil, fmt.Errorf("invalid confirmed: %v", record[8])
	}

	// fmt.Printf("%v\n", []string)
	return &params, nil
}

// func (s *ScheduleAddEditParams) Validate() error {
// 	if s.DepartureAirportID == s.ArrivalAirportID {
// 		return fmt.Errorf("departure_airport_id and arrival_airport_id must be different")
// 	}

// 	return nil
// }
