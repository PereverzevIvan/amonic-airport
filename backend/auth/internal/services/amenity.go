package service

import (
	"time"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
)

type AmenityRepo interface {
	GetAll() ([]*models.Amenity, error)
	CountAll(params *models.AmenityCountAllParams) (*models.AmenityCountAllResult, error)
	GetTicketAmenities(params *models.GetTicketAmenitiesParams) ([]*models.AmenityTicket, error)
	GetCabinTypeDefaultAmenities(params *models.GetCabinTypeDefaultAmenitiesParams) ([]*models.Amenity, error)
	EditTicketAmenities(prev_amenity_ids []int, default_amenity_ids []int, params *models.EditTicketAmenitiesParams) error
}

type amenityService struct {
	amenityRepo  AmenityRepo
	ticketRepo   TicketRepo
	scheduleRepo ScheduleRepo
}

func NewAmenityService(
	amenityRepo AmenityRepo,
	ticketRepo TicketRepo,
	scheduleRepo ScheduleRepo,
) amenityService {
	return amenityService{
		amenityRepo:  amenityRepo,
		ticketRepo:   ticketRepo,
		scheduleRepo: scheduleRepo,
	}
}

func (s amenityService) GetAll() ([]*models.Amenity, error) {
	return s.amenityRepo.GetAll()
}

func (s amenityService) GetTicketAmenities(params *models.GetTicketAmenitiesParams) ([]*models.AmenityTicket, error) {
	return s.amenityRepo.GetTicketAmenities(params)
}

func (s amenityService) GetCabinTypeDefaultAmenities(params *models.GetCabinTypeDefaultAmenitiesParams) ([]int, error) {
	cabin_type_default_amenities, err := s.amenityRepo.GetCabinTypeDefaultAmenities(params)
	if err != nil {
		return nil, err
	}

	amenities_ids := make([]int, len(cabin_type_default_amenities))
	for i, amenity := range cabin_type_default_amenities {
		amenities_ids[i] = amenity.ID
	}

	return amenities_ids, nil
}

func (s amenityService) EditTicketAmenities(params *models.EditTicketAmenitiesParams) error {
	ticket, err := s.ticketRepo.GetByID(params.TicketID)
	if err != nil {
		return err
	}

	schedule, err := s.scheduleRepo.GetByID(ticket.ScheduleID)
	if err != nil {
		return err
	}

	// Проверяем, что время рейса не истекло
	if schedule.Outbound.Before(time.Now().Add(time.Hour * 24)) {
		return models.ErrCantEditAmenitiesTimeExpired
	}

	// получаем предыдущие купленные услуги
	prev_ticket_amenities, err := s.GetTicketAmenities(&models.GetTicketAmenitiesParams{
		TicketID: params.TicketID,
	})
	if err != nil {
		return err
	}

	// формируем массив ID предыдущих купленных услуг
	prev_amenity_ids := make([]int, len(prev_ticket_amenities))
	for i, amenity := range prev_ticket_amenities {
		prev_amenity_ids[i] = amenity.AmenityID
	}

	// получаем услуги по умолчанию по типу кабины
	default_amenity_ids, err := s.GetCabinTypeDefaultAmenities(&models.GetCabinTypeDefaultAmenitiesParams{
		CabinTypeID: ticket.CabinTypeID,
	})
	if err != nil {
		return err
	}

	// Сохраняем изменения
	err = s.amenityRepo.EditTicketAmenities(prev_amenity_ids, default_amenity_ids, params)
	if err != nil {
		return err
	}

	return nil
}

func (s amenityService) CountAll(params *models.AmenityCountAllParams) (*models.AmenityCountAllResult, error) {
	return s.amenityRepo.CountAll(params)
}
