package service

import "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"

type CountryRepo interface {
	GetByID(id int) (*models.Country, error)
	GetByName(name string) (*models.Country, error)
}

type CountryService struct {
	CountryRepo CountryRepo
}

func NewCountryService(repo CountryRepo) CountryService {
	return CountryService{CountryRepo: repo}
}

func (s CountryService) GetByID(id int) (*models.Country, error) {
	return s.CountryRepo.GetByID(id)
}

func (s CountryService) GetByName(name string) (*models.Country, error) {
	return s.CountryRepo.GetByName(name)
}
