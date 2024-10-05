package models

import (
	"fmt"
	"time"
)

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

type NewUserParams struct {
	OfficeID  int       `json:"office_id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  time.Time `json:"birthday"`
}

func (params NewUserParams) ToUser() (*User, error) {
	if params.OfficeID == 0 {
		return nil, fmt.Errorf("office_id can't be 0")
	}

	return &User{
		ID:        0,
		OfficeID:  params.OfficeID,
		RoleID:    int(KRoleUser),
		Email:     params.Email,
		Password:  params.Password,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Birthday:  params.Birthday,
		Active:    true,
	}, nil
}

type UpdateUserParams struct {
	ID        int    `json:"id"`
	OfficeID  int    `json:"office_id"`
	RoleID    int    `json:"role_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (params UpdateUserParams) ToUser() (*User, error) {
	if params.ID == 0 {
		return nil, fmt.Errorf("id can't be 0")
	}
	if params.RoleID == 0 {
		return nil, fmt.Errorf("role_id can't be 0")
	}
	if params.OfficeID == 0 {
		return nil, fmt.Errorf("office_id can't be 0")
	}
	if params.Email == "" {
		return nil, fmt.Errorf("email can't be empty")
	}

	return &User{
		ID:        params.ID,
		OfficeID:  params.OfficeID,
		RoleID:    params.RoleID,
		Email:     params.Email,
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}, nil
}

type UpdateIsActiveUserParams struct {
	ID       int  `json:"id"`
	IsActive bool `json:"is_active"`
}
