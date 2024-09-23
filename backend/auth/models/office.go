package models

type Office struct {
	ID        int    `json:"id"`
	CountryID int    `json:"country_id"`
	Title     string `json:"title"`
	Phone     string `json:"phone"`
	Contact   string `json:"contact"`
}

func (Office) TableName() string {
	return "offices"
}
