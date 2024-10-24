package models

type Group struct {
	ID   int    `json:"id" gorm:"Column:ID"`
	Name string `json:"name" gorm:"Column:Name"`
}

func (*Group) TableName() string {
	return "groups"
}
