package models

type Amenity struct {
	ID      int     `json:"id" gorm:"column:ID;primary_key"`
	Service string  `json:"service" gorm:"column:Service"`
	Price   float64 `json:"price" gorm:"column:Price"`
}

func (*Amenity) TableName() string {
	return "amenities"
}

type AmenityTicket struct {
	AmenityID int      `json:"amenity_id" gorm:"column:AmenityID"`
	Amenity   *Amenity `json:"amenity" gorm:"foreignKey:AmenityID"`
	TicketID  int      `json:"ticket_id" gorm:"column:TicketID"`
	Price     float64  `json:"price" gorm:"column:Price"`
}

func (*AmenityTicket) TableName() string {
	return "amenitiestickets"
}
