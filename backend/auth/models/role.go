package models

type Role struct {
	ID    int
	Title string
}

func (Role) TableName() string {
	return "role"
}
