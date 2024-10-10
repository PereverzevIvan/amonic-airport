package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type ScheduleRepo struct {
	Conn *gorm.DB
}

func NewScheduleRepo(conn *gorm.DB) ScheduleRepo {
	return ScheduleRepo{
		Conn: conn,
	}
}

func (r ScheduleRepo) GetAll(params *models.SchedulesParams) (*[]models.Schedule, error) {
	var schedules []models.Schedule

	err := r.Conn.
		Scopes(
			ScopeSchedulesParams(params),
		).
		Preload("Aircraft").
		Preload("Route").
		Preload("Route.DepartureAirport").
		Preload("Route.ArrivalAirport").
		Model(&models.Schedule{}).
		Find(&schedules).
		Error

	if err != nil {
		return nil, err
	}

	return &schedules, nil
}

func ScopeSchedulesParams(params *models.SchedulesParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		res := db

		if params == nil {
			defaultParams := models.SchedulesParams{}

			switch defaultParams.GetSortBy() {
			case models.KSchedulesSortByDateTime:
				res = res.Order("Outbound DESC")
			case models.KSchedulesSortByConfirmed:
				res = res.Order("Confirmed DESC")
			case models.KSchedulesSortByTicketPrice:
				res = res.Order("EconomyPrice ASC")
			}

			return res
		}

		switch params.GetSortBy() {
		case models.KSchedulesSortByDateTime:
			res = res.Order("Outbound DESC")
		case models.KSchedulesSortByConfirmed:
			res = res.Order("Confirmed DESC")
		case models.KSchedulesSortByTicketPrice:
			res = res.Order("EconomyPrice ASC")
		}

		if params.ArrivalAirportID != nil || params.DepartureAirportID != nil {
			res = res.Joins("JOIN routes ON routes.ID = schedules.RouteID")
		}
		if params.ArrivalAirportID != nil {
			res = res.Where("routes.ArrivalAirportID = ?", *params.ArrivalAirportID)
		}
		if params.DepartureAirportID != nil {
			res = res.Where("routes.DepartureAirportID = ?", *params.DepartureAirportID)
		}

		if params.FlightNumber != "" {
			res = res.Where("FlightNumber = ?", params.FlightNumber)
		}
		if params.Outbound != "" {
			res = res.Where("Outbound LIKE ?", params.Outbound+"%")
		}

		return res
	}
}

func (r ScheduleRepo) GetByID(schedule_id int) (*models.Schedule, error) {
	var schedule models.Schedule

	err := r.Conn.
		Preload("Aircraft").
		Preload("Route").
		Preload("Route.DepartureAirport").
		Preload("Route.ArrivalAirport").
		Model(&models.Schedule{}).
		Where("ID = ?", schedule_id).
		First(&schedule).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, models.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r ScheduleRepo) UpdateConfirmed(schedule_id int, set_confirmed bool) error {

	err := r.Conn.Model(&models.Schedule{ID: schedule_id}).
		Update("Confirmed", set_confirmed).Error

	if err == gorm.ErrRecordNotFound {
		return models.ErrNotFound
	}

	return err
}

func (r ScheduleRepo) UpdateByID(schedule_id int, params *models.ScheduleUpdateParams) error {
	query := r.Conn.Model(&models.Schedule{ID: schedule_id})

	update_fields := []string{}
	update_values_map := make(map[string]interface{})

	if params.EconomyPrice != nil {
		update_fields = append(update_fields, "EconomyPrice")
		update_values_map["EconomyPrice"] = *params.EconomyPrice
	}

	if (params.Date != nil) || (params.Time != nil) {
		if params.Date != nil {
			update_fields = append(update_fields, "Date")
			update_values_map["Date"] = *params.Date
		}
		if params.Time != nil {
			update_fields = append(update_fields, "Time")
			update_values_map["Time"] = *params.Time
		}

		update_fields = append(update_fields, "Outbound")
		switch {
		case params.Date != nil && params.Time != nil:
			update_values_map["Outbound"] = gorm.Expr("CONCAT(?, ' ', ?)", *params.Date, *params.Time)
		case params.Date != nil:
			update_values_map["Outbound"] = gorm.Expr("CONCAT(?, ' ', TIME(Outbound))", *params.Date)
		case params.Time != nil:
			update_values_map["Outbound"] = gorm.Expr("CONCAT(DATE(Outbound), ' ', ?)", *params.Time)
		}
	}

	// log.Info(schedule_id)
	// log.Info(update_fields)
	// log.Info(update_values_map)

	err := query.
		Select(update_fields).
		Updates(update_values_map).
		Error

	if err == gorm.ErrRecordNotFound {
		return models.ErrNotFound
	}

	return err
}

func (r ScheduleRepo) Add(schedule *models.Schedule) error {
	err := r.Conn.
		Create(schedule).
		Error

	if IsUniqueConstraintError(err) {
		return models.ErrUnique
	}

	if IsForeignKeyConstraintError(err) {
		return models.ErrFK
	}

	return err
}

func (r ScheduleRepo) Edit(schedule *models.Schedule) error {

	// log.Info("edit:", *schedule)
	res := r.Conn.Model(&schedule).
		Debug().
		Select(
			"RouteID",
			"AircraftID",
			"EconomyPrice",
			"Confirmed",
			"Outbound",
			"Time",
		).
		Where("Date = ?", schedule.Date).
		Where("FlightNumber = ?", schedule.FlightNumber).
		Updates(*schedule)

	err := res.Error

	if IsUniqueConstraintError(err) {
		return models.ErrUnique
	}

	if IsForeignKeyConstraintError(err) {
		return models.ErrFK
	}

	return err
}

// func (r ScheduleRepo) Update(user *models.User) error {

// 	err := r.Conn.Model(&user).
// 		Select("FirstName", "LastName", "Email", "RoleID", "OfficeID").
// 		Updates(*user).Error

// 	if IsUniqueConstraintError(err) {
// 		return models.ErrDuplicatedEmail
// 	}

// 	if IsForeignKeyConstraintError(err) &&
// 		strings.Contains(err.Error(), "`FK_Users_Offices` FOREIGN KEY (`OfficeID`) foreignKey `offices` (`ID`))") {
// 		return models.ErrFKOfficeIDNotFound
// 	}

// 	return err
// }

// func (r ScheduleRepo) UpdateActive(user_id int, is_active bool) error {
// 	err := r.Conn.Model(&models.User{ID: user_id}).
// 		Update("Active", is_active).Error

// 	return err
// }

// func (r ScheduleRepo) GetByID(user_id int) (*models.User, error) {
// 	var user models.User
// 	err := u.Conn.First(&user, user_id).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, err
// }
