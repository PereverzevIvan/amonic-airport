package models

import "time"

type User struct {
	ID        int       `json:"id"`
	OfficeID  int       `json:"office_id"`
	RoleID    int       `json:"role_id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  time.Time `json:"birthday"`
	Active    bool      `json:"active"`
}

func (User) TableName() string {
	return "users"
}
