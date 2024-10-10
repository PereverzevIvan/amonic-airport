package models

type Aircraft struct {
	ID            int     `json:"id" gorm:"Column:ID"`
	Name          string  `json:"name" gorm:"Column:Name"`
	MakeModel     *string `json:"make_model" gorm:"Column:MakeModel"`
	TotalSeats    int     `json:"total_seats" gorm:"Column:TotalSeats"`
	EconomySeats  int     `json:"economy_seats" gorm:"Column:EconomySeats"`
	BusinessSeats int     `json:"business_seats" gorm:"Column:BusinessSeats"`
}

func (Aircraft) TableName() string {
	return "aircrafts"
}
