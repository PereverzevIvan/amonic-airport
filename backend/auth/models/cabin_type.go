package models

type CabinType struct {
	ID   uint   `json:"id" gorm:"Column:ID"`
	Name string `json:"name" gorm:"Column:Name"`
}

func (CabinType) TableName() string {
	return "cabintypes"
}

type ECabinType int

const (
	KCabinTypeNone       ECabinType = iota // 0
	KCabinTypeEconomy                      // 1
	KCabinTypeBusiness                     // 2
	KCabinTypeFirstClass                   // 3
)
