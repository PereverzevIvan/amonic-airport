package models

type PaginateParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

const KMaxLimit = 100
