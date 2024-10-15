package service

import "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"

type SummaryRepo interface {
	GetFlightsInfo(params *models.SummaryParams, summary *models.Summary) error
	GetTopCustomersInfo(params *models.SummaryParams, summary *models.Summary) error
	GetTopFlightsInfo(params *models.SummaryParams, summary *models.Summary) error
	GetTopOfficesInfo(params *models.SummaryParams, summary *models.Summary) error
	GetRevenueFromTicketSales(params *models.SummaryParams, summary *models.Summary) error
	GetWeeklyReportOfPercentageOfEmptySeats(params *models.SummaryParams, summary *models.Summary) error
}

type summaryService struct {
	summaryRepo SummaryRepo
}

func NewSummaryService(
	summaryRepo SummaryRepo,
) summaryService {
	return summaryService{
		summaryRepo: summaryRepo,
	}
}

func (s summaryService) GetFlightsInfo(params *models.SummaryParams, summary *models.Summary) error {
	return s.summaryRepo.GetFlightsInfo(params, summary)
}

func (s summaryService) GetTopCustomersInfo(params *models.SummaryParams, summary *models.Summary) error {
	return s.summaryRepo.GetTopCustomersInfo(params, summary)
}

func (s summaryService) GetTopFlightsInfo(params *models.SummaryParams, summary *models.Summary) error {
	return s.summaryRepo.GetTopFlightsInfo(params, summary)
}

func (s summaryService) GetTopOfficesInfo(params *models.SummaryParams, summary *models.Summary) error {
	return s.summaryRepo.GetTopOfficesInfo(params, summary)
}

func (s summaryService) GetRevenueFromTicketSales(params *models.SummaryParams, summary *models.Summary) error {
	return s.summaryRepo.GetRevenueFromTicketSales(params, summary)
}

func (s summaryService) GetWeeklyReportOfPercentageOfEmptySeats(params *models.SummaryParams, summary *models.Summary) error {
	return s.summaryRepo.GetWeeklyReportOfPercentageOfEmptySeats(params, summary)
}
