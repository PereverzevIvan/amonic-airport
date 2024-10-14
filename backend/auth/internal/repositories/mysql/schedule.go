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

// Вспомогательная стуктура для восстановления ответа в SearchFlights
type FlightNode struct {
	FlightID      int
	PrevFlightIdx int
}

const (
	MAX_TRANSFERS_COUNT           int = 3
	MINIMUM_TRANSFER_TIME_MINUTES int = 60
)

// SearchFlights - поиск полётов между двумя аэропортами с пересадками
func (r ScheduleRepo) SearchFlights(params *models.SearchFlightsParams) ([][]*models.Schedule, error) {
	prev_visited_airports_set := map[int]bool{} // посещенные на прошлых итерациях аэропорты
	cur_visited_airports_set := map[int]bool{}  // посещенные на текущей итерации аэропорты

	flights_map := map[int]models.Schedule{} // Полёты по id

	prev_flights := []FlightNode{} // массив для восстановления пути
	cur_flights := []FlightNode{}  // текущие полёты, которые нужно обработать
	next_flights := []FlightNode{} // следующие полёты, которые нужно обработать. Когда cur_flights заканчивается, то меняем их, а next_flights очищаем,

	var results [][]*models.Schedule // сами полёты от A до B, с возможными перелётами

	// Получаем все полёты из пункта A, и добавляем их в очередь
	schedules, err := r.findInitialFlights(params)
	if err != nil {
		return nil, err
	}

	// Заполняем изначальными полётами из пункта А
	for _, schedule := range schedules {
		// Если аэропорт прибытия совпадает с аэропортом вылета, то добавляем полёт в результат
		if schedule.Route.ArrivalAirportID == params.ArrivalAirportID {
			results = append(results, []*models.Schedule{&schedule})
			continue
		}
		// Добавляем полёт для обработки
		cur_flights = append(cur_flights, FlightNode{
			FlightID:      schedule.ID,
			PrevFlightIdx: -1,
		})
		// Запоминаем полёты
		flights_map[schedule.ID] = schedule
		// Помечаем аэропорты, как посещенные
		prev_visited_airports_set[schedule.Route.DepartureAirportID] = true
		prev_visited_airports_set[schedule.Route.ArrivalAirportID] = true

	}

	// Основной цикл
	for transfers_count := 0; transfers_count < MAX_TRANSFERS_COUNT; transfers_count++ {
		if len(cur_flights) == 0 {
			break
		}

		for _, cur_flight_node := range cur_flights {
			cur_schedule := flights_map[cur_flight_node.FlightID]

			// Добавляем текущий перелёт в пути, для восстановления пути
			prev_flights = append(prev_flights, cur_flight_node)

			// Находим маршруты из пересадочного пункта
			schedules, err = r.findTransferFlights(&cur_schedule)
			if err != nil {
				return nil, err
			}

			// Обрабатываем новые перелёты
			for _, schedule := range schedules {
				// Если мы уже посещали этот аэропорт, на прошлых итерациях, то пропускаем его
				if _, ok := prev_visited_airports_set[schedule.Route.ArrivalAirportID]; ok {
					continue
				}

				// Найден конечный маршрут
				if schedule.Route.ArrivalAirportID == params.ArrivalAirportID {
					// Добавляем полёт в результат
					results = append(
						results,
						restoreFlightPath(flights_map, prev_flights, &schedule))
					continue
				}

				// Добавляем полёт в мапу
				if _, ok := flights_map[schedule.ID]; !ok {
					flights_map[schedule.ID] = schedule
				}

				// Помечаем посещенный аэропорт
				cur_visited_airports_set[schedule.Route.ArrivalAirportID] = true

				// Добавляем новый перелёт в очередь
				next_flights = append(next_flights, FlightNode{
					FlightID:      schedule.ID,
					PrevFlightIdx: len(prev_flights) - 1,
				})
			}
		}

		cur_flights = next_flights
		next_flights = nil
		for airport_id := range cur_visited_airports_set {
			prev_visited_airports_set[airport_id] = true
			// Очищаем посещенные аэропорты на текущей итерации
			delete(cur_visited_airports_set, airport_id)
		}

	}

	return results, nil
}

func restoreFlightPath(
	flights_map map[int]models.Schedule,
	prev_flights []FlightNode,
	cur_schedule *models.Schedule,
) []*models.Schedule {
	path := []*models.Schedule{cur_schedule}

	for i := len(prev_flights) - 1; i >= 0; i = prev_flights[i].PrevFlightIdx {
		cur_flight := flights_map[prev_flights[i].FlightID]
		path = append(path, &cur_flight)
	}

	// Переворачиваем массив
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func (r *ScheduleRepo) findInitialFlights(params *models.SearchFlightsParams) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.Conn.
		// Debug().
		Scopes(
			scopeSearchInitialFlights(params),
		).
		Find(&schedules).
		Error

	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func scopeSearchInitialFlights(params *models.SearchFlightsParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db.
			InnerJoins("Route").
			Where("Confirmed = ?", true).
			Where("Route.DepartureAirportID = ?", params.DepartureAirportID)

		if params.IncreaseSearchInterval {
			query = query.
				Where(`schedules.Date BETWEEN
						DATE_SUB(?, INTERVAL ? DAY)
						AND DATE_ADD(?, INTERVAL ? DAY)`,
					params.OutboundDate, 3,
					params.OutboundDate, 3,
				)
		} else {
			query = query.Where(`schedules.Date = ?`, params.OutboundDate)
		}

		return query
	}
}

func (r *ScheduleRepo) findTransferFlights(cur_schedule *models.Schedule) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.Conn.Model(&models.Schedule{}).
		// Debug().
		Scopes(
			scopeSearchTransferFlights(cur_schedule),
		).
		Find(&schedules).
		Error

	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func scopeSearchTransferFlights(schedule *models.Schedule) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		query := db.
			InnerJoins("Route").
			Where("Confirmed = ?", true).
			Where("Route.DepartureAirportID = ?", schedule.Route.ArrivalAirportID).
			Where(`schedules.Outbound BETWEEN DATE_ADD(?, INTERVAL ? MINUTE) AND DATE_ADD(?, INTERVAL ? DAY)`,
				schedule.Date[:10]+" "+schedule.Time, schedule.Route.FlightTime+MINIMUM_TRANSFER_TIME_MINUTES,
				schedule.Date[:10]+" "+schedule.Time, 1,
			)

		return query
	}
}
