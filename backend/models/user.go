package models

import "time"

type User struct {
	ID        int ``
	OfficeID  int
	RoleID    int
	Email     string
	Password  string
	FirstName string
	LastName  string
	Birthday  time.Time
	Active    bool
}

func (User) TableName() string {
	return "users"
}
