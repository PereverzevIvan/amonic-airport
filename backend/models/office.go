package models

type Office struct {
	ID        int
	CountryID int
	Title     string
	Phone     string
	Contact   string
}

func (Office) TableName() string {
	return "office"
}
