package models

import (
	"fmt"
	"time"
)

type GetTicketAmenitiesParams struct {
	TicketID int `query:"ticket_id"`
}

func (params *GetTicketAmenitiesParams) Validate() error {
	if params.TicketID <= 0 {
		return fmt.Errorf("ticket_id is required and must be greater than 0")
	}

	return nil
}

type GetCabinTypeDefaultAmenitiesParams struct {
	CabinTypeID int `query:"cabin_type_id"`
}

func (params *GetCabinTypeDefaultAmenitiesParams) Validate() error {
	if params.CabinTypeID <= 0 {
		return fmt.Errorf("cabin_type_id is required and must be greater than 0")
	}

	return nil
}

type EditTicketAmenitiesParams struct {
	TicketID   int   `json:"ticket_id"`
	AmenityIDs []int `json:"amenity_ids"`
}

func (params *EditTicketAmenitiesParams) Validate() error {
	if params.TicketID <= 0 {
		return fmt.Errorf("ticket_id is required and must be greater than 0")
	}
	if params.AmenityIDs == nil {
		return fmt.Errorf("amenity_ids are required")
	}

	return nil
}

type AmenityCountAllParams struct {
	FlightNumber string `query:"flight_number"`
	FromDate     string `query:"from_date"`
	ToDate       string `query:"to_date"`
}

func (params *AmenityCountAllParams) Validate() error {
	if params.FlightNumber == "" {
		_, err := time.Parse("2006-01-02", params.FromDate)
		if err != nil {
			return fmt.Errorf("count in range requires 'from_date', %v", err)
		}

		_, err = time.Parse("2006-01-02", params.ToDate)
		if err != nil {
			return fmt.Errorf("count in range requires 'to_date', %v", err)
		}

		return nil
	}

	_, err := time.Parse("2006-01-02", params.FromDate)
	if err != nil {
		return fmt.Errorf("count for flight requires 'from_date', %v", err)
	}

	return nil
}

type AmenityCountAllResult struct {
	CabinTypes     map[int]string   `json:"cabin_types"`
	AmenitiesCount map[string][]int `json:"amenities_count"`
}
