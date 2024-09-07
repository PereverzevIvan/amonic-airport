package user

import "time"

type User struct {
	ID        int
	OfficeID  int
	RoleID    int
	Email     string
	Password  string
	FirstName string
	LastName  string
	Birthday  time.Time
	Active    bool
}
