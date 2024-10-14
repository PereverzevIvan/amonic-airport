package mysql_repo

import (
	"math/rand"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type TicketRepo struct {
	Conn *gorm.DB
}

func NewTicketRepo(conn *gorm.DB) TicketRepo {
	return TicketRepo{
		Conn: conn,
	}
}

func (r TicketRepo) CountTakenSeatsForShedule(schedule_id int, cabin_type int) (int, error) {
	var count int64

	err := r.Conn.
		Model(&models.Ticket{}).
		Where("Confirmed = ?", true).
		Where("ScheduleID = ? AND CabinTypeID = ?", schedule_id, cabin_type).
		Count(&count).
		Error

	if err != nil {
		return 0, err
	}

	countInt := int(count)

	return countInt, nil
}

func (r TicketRepo) BookTickets(params *models.TicketsBookingParams) ([]*models.Ticket, error) {
	tickets := make([]*models.Ticket, 0)

	bookTicketsForSchedules := func(
		passenger *models.TicketPassengerInfo,
		schedule_ids []int,
	) error {
		booking_reference := r.generateBookingReference()
		for _, schedule_id := range schedule_ids {
			ticket := &models.Ticket{
				ScheduleID:        schedule_id,
				CabinTypeID:       params.CabinType,
				UserID:            params.UserID,
				FirstName:         passenger.FirstName,
				LastName:          passenger.LastName,
				Phone:             passenger.Phone,
				PassportNumber:    passenger.PassportNumber,
				PassportCountryID: passenger.PassportCountryID,
				BookingReference:  booking_reference,
				Email:             nil,
				Confirmed:         false,
			}

			err := r.Conn.Create(&ticket).Error
			if err != nil {
				return err
			}

			err = r.Conn.
				Model(&models.Schedule{}).
				First(&ticket.Schedule, "ID = ?", schedule_id).
				Error
			if err != nil {
				return err
			}

			tickets = append(tickets, ticket)
		}
		return nil
	}

	for _, passenger := range params.Passengers {
		err := bookTicketsForSchedules(
			&passenger,
			params.OutboundScheduleIDs,
		)
		if err != nil {
			return nil, err
		}

		err = bookTicketsForSchedules(
			&passenger,
			params.InboundScheduleIDs,
		)
		if err != nil {
			return nil, err
		}
	}

	return tickets, nil
}

func (r TicketRepo) generateBookingReference() string {
	// Seed the random number generator
	booking_reference := make([]byte, 6)

	for i := range booking_reference {
		booking_reference[i] = byte('A' + rand.Intn('Z'-'A'+1))
	}
	return string(booking_reference)
}

func (r TicketRepo) ChangeTicketsStatus(ticket_ids []int, set_confirmed bool) error {
	// Note the use of tx as the database handle once you are within a transaction
	tx := r.Conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, ticket_id := range ticket_ids {
		err := tx.
			Model(&models.Ticket{ID: ticket_id}).
			Update("Confirmed", set_confirmed).
			Error

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
