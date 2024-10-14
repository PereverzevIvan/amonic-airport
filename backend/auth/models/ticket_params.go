package models

import (
	"fmt"
	"time"
)

type TicketsCountRemainingSeatsParams struct {
	ScheduleIDs []int `json:"schedule_ids"`
}

func (params *TicketsCountRemainingSeatsParams) Validate() error {
	if len(params.ScheduleIDs) == 0 {
		return fmt.Errorf("len of schedule_ids is 0")
	}

	return nil
}

type TicketsRemainingSeatsCount struct {
	EconomySeats    int `json:"economy_seats"`
	BusinessSeats   int `json:"business_seats"`
	FirstClassSeats int `json:"first_class_seats"`
}

type TicketPassengerInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	// Email string `json:"email"`
	Phone string `json:"phone"`

	PassportNumber    string `json:"passport_number"`
	PassportCountryID int    `json:"passport_country_id"`

	Birthday string `json:"birthday"`
}

func (params *TicketPassengerInfo) Validate() error {
	if params.FirstName == "" {
		return fmt.Errorf("first_name is empty")
	}
	if params.LastName == "" {
		return fmt.Errorf("last_name is empty")
	}
	// if params.Email == "" {
	// 	return fmt.Errorf("email is empty")
	// }
	if params.Phone == "" {
		return fmt.Errorf("phone is empty")
	}
	if len(params.Phone) > 14 {
		return fmt.Errorf("phone %v is too long, max size: 14", params.Phone)
	}
	if params.PassportNumber == "" {
		return fmt.Errorf("passport_number is empty")
	}
	if params.PassportCountryID == 0 {
		return fmt.Errorf("passport_country_id is empty")
	}

	_, err := time.Parse("2006-01-02", params.Birthday)
	if err != nil {
		return fmt.Errorf("birthday is invalid %v", err.Error())
	}

	return nil
}

type TicketsBookingParams struct {
	OutboundScheduleIDs []int                 `json:"outbound_schedules"`
	InboundScheduleIDs  []int                 `json:"inbound_schedules"`
	Passengers          []TicketPassengerInfo `json:"passengers"`
	CabinType           int                   `json:"cabin_type"`
	UserID              int                   `json:"-"`
}

func (params *TicketsBookingParams) Validate() error {
	if len(params.OutboundScheduleIDs) == 0 {
		return fmt.Errorf("outbound_schedules is requied")
	}
	if len(params.Passengers) == 0 {
		return fmt.Errorf("passengers are required")
	}
	for _, passenger := range params.Passengers {
		if err := passenger.Validate(); err != nil {
			return err
		}
	}
	if params.CabinType != int(KCabinTypeEconomy) &&
		params.CabinType != int(KCabinTypeBusiness) &&
		params.CabinType != int(KCabinTypeFirstClass) {
		return fmt.Errorf("cabin_type is invalid")
	}

	return nil
}

type TicketIDsParams struct {
	TicketeIDs []int `json:"tickets"`
}

func (params *TicketIDsParams) Validate() error {
	if len(params.TicketeIDs) == 0 {
		return fmt.Errorf("tickets are required")
	}
	return nil
}

type TicketsBookResult struct {
	Tickets   []*Ticket `json:"tickets"`
	TotalCost float64   `json:"total_cost"`
}
