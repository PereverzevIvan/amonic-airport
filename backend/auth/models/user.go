package models

import "time"

type User struct {
	ID        int       `json:"id" gorm:"Column:ID"`
	OfficeID  int       `json:"office_id" gorm:"Column:OfficeID"`
	RoleID    int       `json:"role_id" gorm:"Column:RoleID"`
	Email     string    `json:"email" gorm:"Column:Email"`
	Password  string    `json:"-" gorm:"Column:Password"`
	FirstName string    `json:"first_name" gorm:"Column:FirstName"`
	LastName  string    `json:"last_name" gorm:"Column:LastName"`
	Birthday  time.Time `json:"birthday" gorm:"Column:Birthdate"`
	Active    bool      `json:"active" gorm:"Column:Active"`
}

func (User) TableName() string {
	return "users"
}
