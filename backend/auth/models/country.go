package models

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (Country) TableName() string {
	return "countries"
}
