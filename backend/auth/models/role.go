package models

type Role struct {
	ID    int
	Title string
}

func (Role) TableName() string {
	return "role"
}

type ERole int

const (
	KRoleNone  ERole = iota // 0
	KRoleAdmin              // 1
	KRoleUser               // 2
)
