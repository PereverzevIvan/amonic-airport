package service

import (
	"encoding/csv"
	"io"
	"mime/multipart"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v2/log"
)

type ScheduleRepo interface {
	GetAll(*models.SchedulesParams) (*[]models.Schedule, error)
	GetByID(schedule_id int) (*models.Schedule, error)
	UpdateConfirmed(schedule_id int, set_confirmed bool) error
	UpdateByID(schedule_id int, params *models.ScheduleUpdateParams) error
	Add(schedule *models.Schedule) error
	Edit(schedule *models.Schedule) error
}

type RouteRepo interface {
	// GetByID(route_id int) (*models.Route, error)
	GetIDByAirportCodes(departure_airport_code, arrival_airport_code string) (int, error)
}

type scheduleService struct {
	scheduleRepo ScheduleRepo
	routeRepo    RouteRepo
}

func NewScheduleService(
	scheduleRepo ScheduleRepo,
	routeRepo RouteRepo,
) scheduleService {
	return scheduleService{
		scheduleRepo: scheduleRepo,
		routeRepo:    routeRepo,
	}
}

func (s scheduleService) GetAll(params *models.SchedulesParams) (*[]models.Schedule, error) {
	return s.scheduleRepo.GetAll(params)
}

func (s scheduleService) GetByID(schedule_id int) (*models.Schedule, error) {
	return s.scheduleRepo.GetByID(schedule_id)
}

func (s scheduleService) UpdateConfirmed(schedule_id int, set_confirmed bool) error {
	return s.scheduleRepo.UpdateConfirmed(schedule_id, set_confirmed)
}

func (s scheduleService) UpdateByID(schedule_id int, params *models.ScheduleUpdateParams) error {
	return s.scheduleRepo.UpdateByID(schedule_id, params)
}

func (s scheduleService) ApplyChangesFromSCV(src *multipart.File) (models.SchedulesUploadResult, error) {
	// Create a buffer to store the file contents
	reader := csv.NewReader(*src)
	// Set FieldsPerRecord to -1 to allow variable number of fields
	reader.FieldsPerRecord = -1

	res := models.SchedulesUploadResult{}

	unique_rows_set := map[string]bool{}

	// Read the CSV line by line
	for {
		record, err := reader.Read()
		log.Info(res)
		if err == io.EOF {
			break // End of file reached
		}

		res.TotalRowsCnt++
		if err != nil {
			res.FailedRowsCnt++
			// Log the error and continue to the next record
			log.Info("Error reading record: %v\n", err)
			continue
		}

		// Process the valid record (e.g., print it)
		log.Info(record)

		params, err := models.ParseScheduleAddEditCommandFromCSVRecord(record)
		if err != nil {
			log.Error(err)
			if err.Error() == models.ErrCSVMissingFields.Error() {
				res.MissingFieldsRowsCnt++
				continue
			}

			res.FailedRowsCnt++
			continue
		}

		// Check if the record is unique
		if _, ok := unique_rows_set[params.ToStrHash()]; ok {
			res.DuplicatedRowsCnt++
			continue
		}
		unique_rows_set[params.ToStrHash()] = true

		// Process the record
		schedule := params.ToSchedule()

		// log.Info(schedule)

		route_id, err := s.routeRepo.GetIDByAirportCodes(
			params.DepartureAirportCode,
			params.ArrivalAirportCode,
		)
		if err != nil {
			log.Error(err)
			res.FailedRowsCnt++
			continue
		}
		schedule.RouteID = route_id

		if params.IsAddCommand {
			err = s.scheduleRepo.Add(schedule)
		} else {
			err = s.scheduleRepo.Edit(schedule)
		}

		if err != nil {
			log.Error(err)
			if err == models.ErrUnique {
				res.DuplicatedRowsCnt++
				continue
			}

			res.FailedRowsCnt++
			continue
		}

		res.SuccessfulRowsCnt++
	}

	return res, nil
}
