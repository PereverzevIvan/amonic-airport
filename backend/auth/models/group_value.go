package models

type GroupValue struct {
	ID      int    `json:"id" gorm:"Column:ID"`
	GroupID int    `json:"group_id" gorm:"Column:GroupID"`
	Group   *Group `json:"group" gorm:"foreignKey:GroupID"`
	Name    string `json:"name" gorm:"Column:Name"`
}

func (*GroupValue) TableName() string {
	return "group_values"
}

type GroupWithValues struct {
	Group
	Values []GroupValue `json:"values"`
}
