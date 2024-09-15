package models

type PaginateParams struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}
