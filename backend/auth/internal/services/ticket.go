package service

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
)

type TicketRepo interface {
	CountTakenSeatsForShedule(schedule_id int, cabin_type int) (int, error)
	BookTickets(params *models.TicketsBookingParams) ([]*models.Ticket, error)
	ChangeTicketsStatus(ticket_ids []int, set_confirmed bool) error
	GetAll(params *models.TicketsGetAllParams) ([]*models.Ticket, error)
	GetByID(ticket_id int) (*models.Ticket, error)
}

type ticketService struct {
	ticketRepo   TicketRepo
	scheduleRepo ScheduleRepo
}

func NewTicketService(
	ticketRepo TicketRepo,
	scheduleRepo ScheduleRepo,
) ticketService {
	return ticketService{
		ticketRepo:   ticketRepo,
		scheduleRepo: scheduleRepo,
	}
}

func (s ticketService) CountRemainingSeats(params *models.TicketsCountRemainingSeatsParams) (*models.TicketsRemainingSeatsCount, error) {
	result := models.TicketsRemainingSeatsCount{}

	for _, schedule_id := range params.ScheduleIDs {
		schedule, err := s.scheduleRepo.GetByID(schedule_id)
		if err != nil {
			return nil, err
		}

		taken_economy_seats, err := s.ticketRepo.CountTakenSeatsForShedule(schedule.ID, int(models.KCabinTypeEconomy))
		if err != nil {
			return nil, err
		}
		taken_business_seats, err := s.ticketRepo.CountTakenSeatsForShedule(schedule.ID, int(models.KCabinTypeBusiness))
		if err != nil {
			return nil, err
		}
		taken_first_class_seats, err := s.ticketRepo.CountTakenSeatsForShedule(schedule.ID, int(models.KCabinTypeFirstClass))
		if err != nil {
			return nil, err
		}

		aircraft := schedule.Aircraft

		remaining_economy_seats := aircraft.EconomySeats - taken_economy_seats
		remaining_business_seats := aircraft.EconomySeats - taken_business_seats
		remaining_first_class_seats := aircraft.EconomySeats - taken_first_class_seats

		if result.EconomySeats == 0 || result.EconomySeats > remaining_economy_seats {
			result.EconomySeats = remaining_economy_seats
		}
		if result.BusinessSeats == 0 || result.BusinessSeats > remaining_business_seats {
			result.BusinessSeats = remaining_business_seats
		}
		if result.FirstClassSeats == 0 || result.FirstClassSeats > remaining_first_class_seats {
			result.FirstClassSeats = remaining_first_class_seats
		}
	}

	return &result, nil
}

func (s ticketService) BookTickets(params *models.TicketsBookingParams) ([]*models.Ticket, error) {
	count_remaining_seats, err := s.CountRemainingSeats(&models.TicketsCountRemainingSeatsParams{
		ScheduleIDs: params.OutboundScheduleIDs,
	})
	if err != nil {
		return nil, err
	}

	switch params.CabinType {
	case int(models.KCabinTypeEconomy):
		if count_remaining_seats.EconomySeats < len(params.Passengers) {
			return nil, models.ErrNoTicketsAvailable
		}
	case int(models.KCabinTypeBusiness):
		if count_remaining_seats.BusinessSeats < len(params.Passengers) {
			return nil, models.ErrNoTicketsAvailable
		}
	case int(models.KCabinTypeFirstClass):
		if count_remaining_seats.FirstClassSeats < len(params.Passengers) {
			return nil, models.ErrNoTicketsAvailable
		}
	}

	tickets, err := s.ticketRepo.BookTickets(params)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (s ticketService) ChangeTicketsStatus(ticket_ids []int, set_confirmed bool) error {
	return s.ticketRepo.ChangeTicketsStatus(ticket_ids, set_confirmed)
}

// FindSchedulesByBookingReference returns list of schedules associated with given booking reference.
//
// This is used to obtain information about the flight (departure/arrival dates, route, etc.)
// associated with the tickets booked by the user with given booking reference.
func (s ticketService) GetAll(params *models.TicketsGetAllParams) ([]*models.Ticket, error) {
	return s.ticketRepo.GetAll(params)
}

// func (s ticketService) GetAll(params *models.TicketsParams) (*[]models.Ticket, error) {
// 	return s.ticketRepo.GetAll(params)
// }

// func (s ticketService) GetByID(ticket_id int) (*models.Ticket, error) {
// 	return s.ticketRepo.GetByID(ticket_id)
// }
