package service

import "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"

type OfficeRepo interface {
	GetByID(id int) (*models.Office, error)
	GetByTitle(title string) (*models.Office, error)
	GetAll() (*[]models.Office, error)
}

func NewOfficeService(repo OfficeRepo) OfficeService {
	return OfficeService{
		OfficeRepo: repo,
	}
}

type OfficeService struct {
	OfficeRepo OfficeRepo
}

func (s OfficeService) GetByID(id int) (*models.Office, error) {
	return s.OfficeRepo.GetByID(id)
}

func (s OfficeService) GetByTitle(title string) (*models.Office, error) {
	return s.OfficeRepo.GetByTitle(title)
}

func (s OfficeService) GetAll() (*[]models.Office, error) {
	return s.OfficeRepo.GetAll()
}
