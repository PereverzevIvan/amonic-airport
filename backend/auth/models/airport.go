package models

type Airport struct {
	ID        int    `json:"id" gorm:"Column:ID"`
	CountryID int    `json:"country_id" gorm:"Column:CountryID"`
	IATACode  string `json:"iata_code" gorm:"Column:IATACode"`
	Name      string `json:"name" gorm:"Column:Name"`
}

func (Airport) TableName() string {
	return "airports"
}
