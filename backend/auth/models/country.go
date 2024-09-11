package models

type Country struct {
	ID   int
	Name string
}

func (Country) TableName() string {
	return "countries"
}
