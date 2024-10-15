package mysql_repo

import (
	"time"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3/log"
	"gorm.io/gorm"
)

type SummaryRepo struct {
	Conn *gorm.DB
}

func NewSummaryRepo(conn *gorm.DB) SummaryRepo {
	return SummaryRepo{
		Conn: conn,
	}
}

func (r SummaryRepo) GetFlightsInfo(params *models.SummaryParams, summary *models.Summary) error {
	rows, err := r.Conn.Raw(`
		SELECT
			schedules.Confirmed, COUNT(*)
		FROM
			schedules
		WHERE
			schedules.Date BETWEEN ? AND ?
		GROUP BY 
			schedules.Confirmed
		ORDER BY 
			schedules.Confirmed ASC`,
		params.StartDate, params.EndDate).
		Rows()
	if err != nil {
		return err
	}

	var confirmed_type, count int

	for rows.Next() {
		err = rows.Scan(&confirmed_type, &count)
		if err != nil {
			return err
		}

		if confirmed_type == 1 {
			summary.NumberConfirmedFlights = count
		} else {
			summary.NumberCancelledFlights = count
		}
	}

	rows, err = r.Conn.Raw(`
			SELECT
				SUM(routes.FlightTime)
			FROM
				schedules
			INNER JOIN routes ON routes.ID = schedules.RouteID
			WHERE
				schedules.Date BETWEEN ? AND ?
				AND schedules.Confirmed = TRUE`,
		params.StartDate, params.EndDate).
		Rows()

	if err != nil {
		return err
	}

	for rows.Next() {
		var total_time int
		err = rows.Scan(&total_time)
		if err != nil {
			return err
		}

		summary.AverageDailyFlightTimeMinutes = total_time / 30
	}

	return nil
}

func (r SummaryRepo) GetTopCustomersInfo(params *models.SummaryParams, summary *models.Summary) error {
	rows, err := r.Conn.Raw(`
		SELECT
			tickets.Phone,
			tickets.Firstname,
			tickets.Lastname,
			COUNT(*) as count_tickets
		FROM
			schedules
		INNER JOIN tickets ON tickets.ScheduleID = schedules.ID
		WHERE
			schedules.Date BETWEEN ? AND ? 
			AND tickets.Confirmed = TRUE
		GROUP BY
			tickets.Phone, tickets.Firstname, tickets.Lastname
		ORDER BY
			count_tickets DESC
		LIMIT 3`,
		params.StartDate, params.EndDate).
		Rows()
	if err != nil {
		return err
	}

	for rows.Next() {
		var phone, firstname, lastname string
		var count_tickets int
		err = rows.Scan(&phone, &firstname, &lastname, &count_tickets)
		if err != nil {
			return err
		}

		summary.TopCustomersByPurchasedTickets = append(summary.TopCustomersByPurchasedTickets, firstname+" "+lastname)
	}
	return nil
}

func (r SummaryRepo) GetTopFlightsInfo(params *models.SummaryParams, summary *models.Summary) error {

	rows, err := r.Conn.Raw(`
		SELECT
			schedules.Date,
			COUNT(*) AS count_tickets
		FROM
			schedules
		INNER JOIN tickets ON tickets.ScheduleID = schedules.ID
		WHERE
			schedules.Date BETWEEN ? AND ? 
			AND tickets.Confirmed = TRUE
		GROUP BY
			schedules.Date
		ORDER BY
			count_tickets DESC
		LIMIT 1`,
		params.StartDate, params.EndDate).
		Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		rows.Scan(
			&summary.BusiestDay,
			&summary.BusiestDayNumberOfPassengers)
		summary.BusiestDay = summary.BusiestDay[:10]
	}

	rows, err = r.Conn.Raw(`
		SELECT
			schedules.Date,
			COUNT(*) AS count_tickets
		FROM
			schedules
		INNER JOIN tickets ON tickets.ScheduleID = schedules.ID
		WHERE
			schedules.Date BETWEEN ? AND ? 
			AND tickets.Confirmed = TRUE
		GROUP BY
			schedules.Date
		ORDER BY
			count_tickets ASC
		LIMIT 1`,
		params.StartDate, params.EndDate).
		Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		rows.Scan(
			&summary.MostQuietDay,
			&summary.MostQuietDayNumberOfPassengers)
		summary.MostQuietDay = summary.MostQuietDay[:10]
	}
	return nil
}

func (r SummaryRepo) GetTopOfficesInfo(params *models.SummaryParams, summary *models.Summary) error {
	rows, err := r.Conn.Raw(`
		SELECT
			offices.Title,
			COUNT(*) AS count_tickets
		FROM
			offices
		INNER JOIN users ON users.OfficeID = offices.ID
		INNER JOIN tickets ON tickets.UserID = users.ID
		INNER JOIN schedules ON schedules.ID = tickets.ScheduleID
		WHERE
			schedules.Date BETWEEN ? AND ? 
			AND tickets.Confirmed = TRUE
		GROUP BY
			offices.Title
		ORDER BY
			count_tickets DESC
		LIMIT 3`,
		params.StartDate, params.EndDate).
		Rows()

	if err != nil {
		return err
	}
	for rows.Next() {
		var name string
		var count_tickets int
		err = rows.Scan(&name, &count_tickets)
		if err != nil {
			return err
		}

		summary.TopOffices = append(summary.TopOffices, name)
	}
	return nil
}

func (r SummaryRepo) GetRevenueFromTicketSales(params *models.SummaryParams, summary *models.Summary) error {
	end_date, err := time.Parse("2006-01-02", params.EndDate)
	if err != nil {
		return err
	}

	for days_before := 0; days_before < 3; days_before++ {
		date := end_date.
			Add(-time.Hour * 24 * time.Duration(days_before+1)).
			Format("2006-01-02")

		log.Info(date)

		rows, err := r.Conn.Raw(`
			SELECT
				tickets.CabinTypeID,
				schedules.EconomyPrice,
				COUNT(*)
				FROM
				schedules
			INNER JOIN tickets ON tickets.ScheduleID = schedules.ID
			WHERE
			schedules.Date = ?
				AND tickets.Confirmed = TRUE
			GROUP BY
			tickets.CabinTypeID,
				schedules.EconomyPrice`,
			date).
			Rows()

		if err != nil {
			return err
		}

		summary.RevenueFromTicketSales = append(summary.RevenueFromTicketSales, 0)
		for rows.Next() {
			var cabin_type_id, count int
			var economy_price float64
			err = rows.Scan(&cabin_type_id, &economy_price, &count)
			if err != nil {
				return err
			}

			var cabin_type_ratio float64 = 1
			switch cabin_type_id {
			case 2:
				cabin_type_ratio = models.KBusinessRatio
			case 3:
				cabin_type_ratio = models.KFirstClassRation
			}

			var total_price float64 = float64(count) * economy_price * cabin_type_ratio

			summary.RevenueFromTicketSales[days_before] += total_price
		}
	}

	return nil
}

func (r SummaryRepo) GetWeeklyReportOfPercentageOfEmptySeats(params *models.SummaryParams, summary *models.Summary) error {
	end_date, err := time.Parse("2006-01-02", params.EndDate)
	if err != nil {
		return err
	}

	week_end_date := end_date

	week_end_day := end_date.Weekday()
	week_start_date := end_date.Add(-time.Hour * 24 * time.Duration(week_end_day-1))

	for weeks_before := 0; weeks_before < 3; weeks_before++ {
		log.Info(week_start_date, week_start_date.Weekday())
		log.Info(week_end_date, week_end_date.Weekday())

		rows, err := r.Conn.
			// Debug().
			Raw(`
			SELECT
				SUM(aircrafts.TotalSeats) as total_seats
			FROM
				schedules
			INNER JOIN aircrafts ON aircrafts.ID = schedules.AircraftID
			WHERE
				schedules.Date BETWEEN ? AND ?
				AND schedules.Confirmed = TRUE`,
				week_start_date.Format("2006-01-02"),
				week_end_date.Format("2006-01-02"),
			).
			Rows()
		if err != nil {
			return err
		}

		total_seats := 0
		for rows.Next() {
			err = rows.Scan(&total_seats)
			if err != nil {
				return err
			}
		}

		rows, err = r.Conn.
			Raw(`
				SELECT
				COUNT(*) as taken_seats
				FROM
					schedules
				INNER JOIN tickets ON tickets.ScheduleID = schedules.ID
				WHERE
				schedules.Date BETWEEN ? AND ?
					AND schedules.Confirmed = TRUE 
					AND tickets.Confirmed = TRUE`,
				week_start_date.Format("2006-01-02"),
				week_end_date.Format("2006-01-02"),
			).
			Rows()
		if err != nil {
			return err
		}

		taken_seats := 0
		for rows.Next() {
			err = rows.Scan(&taken_seats)
			if err != nil {
				return err
			}
		}

		empty_seats := total_seats - taken_seats
		summary.WeeklyReportOfPercentageOfEmptySeats = append(
			summary.WeeklyReportOfPercentageOfEmptySeats,
			float64(empty_seats)/float64(total_seats),
		)

		// update week start date
		week_end_date = week_start_date.Add(-time.Hour * 24)
		week_start_date = week_start_date.Add(-7 * time.Hour * 24)
	}
	return nil
}
