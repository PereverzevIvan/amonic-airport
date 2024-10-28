package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3/log"
	"gorm.io/gorm"
)

type AmenityRepo struct {
	Conn *gorm.DB
}

func NewAmenityRepo(conn *gorm.DB) AmenityRepo {
	return AmenityRepo{
		Conn: conn,
	}
}

func (r AmenityRepo) GetAll() ([]*models.Amenity, error) {
	amenities := []*models.Amenity{}

	err := r.Conn.
		Model(&models.Amenity{}).
		Find(&amenities).
		Error

	if err != nil {
		return nil, err
	}
	return amenities, nil
}

func (r AmenityRepo) GetTicketAmenities(params *models.GetTicketAmenitiesParams) ([]*models.AmenityTicket, error) {

	amenities := []*models.AmenityTicket{}

	err := r.Conn.
		Model(&models.AmenityTicket{}).
		Where("TicketID = ?", params.TicketID).
		Find(&amenities).
		Error

	if err != nil {
		return nil, err
	}
	return amenities, nil
}

func (r AmenityRepo) GetCabinTypeDefaultAmenities(params *models.GetCabinTypeDefaultAmenitiesParams) ([]*models.Amenity, error) {
	amenities := []*models.Amenity{}

	err := r.Conn.
		Model(&models.Amenity{}).
		Select("amenities.*, amenitiescabintype.*").
		InnerJoins("INNER JOIN amenitiescabintype ON amenitiescabintype.AmenityID = amenities.ID").
		Where("amenitiescabintype.CabinTypeID = ?", params.CabinTypeID).
		Find(&amenities).
		Error

	if err != nil {
		return nil, err
	}
	return amenities, nil
}

func (r AmenityRepo) EditTicketAmenities(prev_amenity_ids []int, default_amenity_ids []int, params *models.EditTicketAmenitiesParams) error {
	amenities, err := r.GetAll()
	if err != nil {
		return err
	}

	log.Info(prev_amenity_ids)
	log.Info(default_amenity_ids)
	log.Info(params.AmenityIDs)

	amenities_map := map[int]*models.Amenity{}
	for _, amenity := range amenities {
		amenities_map[amenity.ID] = amenity
	}

	prev_amenity_ids_map := map[int]bool{}
	for _, amenity_id := range prev_amenity_ids {
		prev_amenity_ids_map[amenity_id] = true
	}

	default_amenity_ids_map := map[int]bool{}
	for _, amenity_id := range default_amenity_ids {
		default_amenity_ids_map[amenity_id] = true
	}

	new_amenities_ids_map := map[int]bool{}
	for _, amenity_id := range params.AmenityIDs {
		new_amenities_ids_map[amenity_id] = true
	}

	tx := r.Conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, prev_amenity_id := range prev_amenity_ids {
		if _, ok := new_amenities_ids_map[prev_amenity_id]; ok {
			continue
		}

		err := tx.
			Debug().
			Model(&models.AmenityTicket{}).
			Where("TicketID = ? AND AmenityID = ?", params.TicketID, prev_amenity_id).
			Delete(&models.AmenityTicket{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, new_amenity_id := range params.AmenityIDs {
		// Добавляем новую услугу, если она не была ранее добавлена и не включена по умолчанию
		if _, ok := prev_amenity_ids_map[new_amenity_id]; ok {
			continue
		}
		if _, ok := default_amenity_ids_map[new_amenity_id]; ok {
			continue
		}

		if _, ok := amenities_map[new_amenity_id]; !ok {
			continue
		}
		// Если услуга бесплатная, то она не может быть добавлена
		if amenity, ok := amenities_map[new_amenity_id]; ok {
			if amenity.Price == 0 {
				continue
			}
		}

		var amenity_price float64
		if _, ok := amenities_map[new_amenity_id]; !ok {
			return gorm.ErrRecordNotFound
		}
		amenity_price = amenities_map[new_amenity_id].Price

		err := tx.
			Debug().
			Model(&models.AmenityTicket{}).
			Create(&models.AmenityTicket{
				TicketID:  params.TicketID,
				AmenityID: new_amenity_id,
				Price:     amenity_price,
			}).Error

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r AmenityRepo) CountAll(params *models.AmenityCountAllParams) (*models.AmenityCountAllResult, error) {
	result := models.AmenityCountAllResult{}
	result.CabinTypes = map[int]string{
		1: "Economy",
		2: "Business",
		3: "First Class",
	}
	result.AmenitiesCount = map[string][]int{}

	amenities, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	amenities_map := map[int]*models.Amenity{}
	for _, amenity := range amenities {
		amenities_map[amenity.ID] = amenity
	}

	log.Info(amenities)

	for _, amenity := range amenities {
		result.AmenitiesCount[amenity.Service] = make([]int, len(result.CabinTypes))
	}

	err = r.addPurchasedAmenities(params, &result, amenities_map)
	if err != nil {
		return nil, err
	}

	err = r.addTicketsDefaultAmenities(params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r AmenityRepo) addPurchasedAmenities(
	params *models.AmenityCountAllParams,
	result *models.AmenityCountAllResult,
	amenities_map map[int]*models.Amenity,
) error {
	rows, err := r.Conn.
		// Debug().
		Table("schedules").
		Scopes(
			scopeCountAmenities(params),
		).
		Select("amenitiestickets.AmenityID, tickets.CabinTypeID, COUNT(*)").
		InnerJoins("INNER JOIN tickets ON tickets.ScheduleID = schedules.ID").
		InnerJoins("INNER JOIN amenitiestickets ON amenitiestickets.TicketID = tickets.ID").
		Group("amenitiestickets.AmenityID").
		Group("tickets.CabinTypeID").
		Rows()
	if err != nil {
		return err
	}

	for rows.Next() {

		var amenity_id, cabin_type_id, count int
		err = rows.Scan(&amenity_id, &cabin_type_id, &count)
		if err != nil {
			return err
		}

		amenity_name := amenities_map[amenity_id].Service
		// log.Info(amenity_id, cabin_type_id, count, amenity_name)
		result.AmenitiesCount[amenity_name][cabin_type_id-1] += count
	}
	return nil
}

func (r AmenityRepo) addTicketsDefaultAmenities(
	params *models.AmenityCountAllParams,
	result *models.AmenityCountAllResult,
) error {
	cabin_type_ids := []int{1, 2, 3}

	tickets_count_by_cabin_type := map[int]int{}
	rows, err := r.Conn.
		Scopes(
			scopeCountAmenities(params),
		).
		Table("schedules").
		Select("tickets.CabinTypeID, COUNT(*)").
		InnerJoins("INNER JOIN tickets ON tickets.ScheduleID = schedules.ID").
		Group("tickets.CabinTypeID").
		Rows()
	if err != nil {
		return err
	}

	for rows.Next() {
		var cabin_type_id, count_tickets int
		err = rows.Scan(&cabin_type_id, &count_tickets)
		if err != nil {
			return err
		}

		tickets_count_by_cabin_type[cabin_type_id] = count_tickets
	}
	log.Info(tickets_count_by_cabin_type)

	for _, cabin_type_id := range cabin_type_ids {
		default_amenities, err := r.GetCabinTypeDefaultAmenities(&models.GetCabinTypeDefaultAmenitiesParams{
			CabinTypeID: cabin_type_id,
		})
		if err != nil {
			return err
		}
		for _, amenity := range default_amenities {
			result.AmenitiesCount[amenity.Service][cabin_type_id-1] += tickets_count_by_cabin_type[cabin_type_id]
		}
	}

	return nil
}

func scopeCountAmenities(params *models.AmenityCountAllParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.
			Where("tickets.Confirmed = ?", true)

		if params == nil {
			return query
		}

		if params.FlightNumber != "" {
			return query.
				Where("schedules.FlightNumber = ?", params.FlightNumber).
				Where("schedules.Date = ?", params.FromDate)
		}

		return query.Where("schedules.Date BETWEEN ? AND ?", params.FromDate, params.ToDate)
	}
}
