package models

import "time"

type Survey struct {
	ID               int       `json:"id" gorm:"Column:ID"`
	Date             time.Time `json:"date" gorm:"Column:Date"`
	RespondentsCount int       `json:"surveyed_count" gorm:"Column:RespondentsCount"`
}

func (*Survey) TableName() string {
	return "surveys"
}
