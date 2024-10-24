package models

type Question struct {
	ID   int    `json:"id" gorm:"Column:ID"`
	Text string `json:"text" gorm:"Column:Text"`
}

func (*Question) TableName() string {
	return "questions"
}
