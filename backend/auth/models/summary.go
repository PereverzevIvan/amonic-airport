package models

import (
	"fmt"
	"time"
)

type Summary struct {
	NumberConfirmedFlights        int `json:"number_confirmed_flights"`
	NumberCancelledFlights        int `json:"number_cancelled_flights"`
	AverageDailyFlightTimeMinutes int `json:"average_daily_flight_time_minutes"`

	TopCustomersByPurchasedTickets []string `json:"top_customer_by_purchased_tickets"`

	BusiestDay                     string `json:"busiest_day"`
	BusiestDayNumberOfPassengers   int    `json:"busiest_day_number_of_passengers"`
	MostQuietDay                   string `json:"most_quiet_day"`
	MostQuietDayNumberOfPassengers int    `json:"most_quiet_day_number_of_passengers"`

	TopOffices []string `json:"top_offices"`

	RevenueFromTicketSales []float64 `json:"revenue_from_ticket_sales"`

	WeeklyReportOfPercentageOfEmptySeats []float64 `json:"weekly_report_of_percentage_of_empty_seats"`

	TimeTakenToGenerateReport int `json:"time_taken_to_generate_report"`
}

type SummaryParams struct {
	StartDate string `json:"start_date" query:"start_date"`
	EndDate   string `json:"end_date" query:"end_date"`
}

func (params SummaryParams) Validate() error {
	_, err := time.Parse("2006-01-02", params.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date: %v", err.Error())
	}
	return nil
}
