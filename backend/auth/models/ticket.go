package models

type Ticket struct {
	ID     int `json:"id" gorm:"Column:ID"`
	UserID int `json:"user_id" gorm:"Column:UserID"`

	ScheduleID int       `json:"schedule_id" gorm:"Column:ScheduleID"`
	Schedule   *Schedule `json:"schedule" gorm:"foreignKey:ScheduleID"`

	CabinTypeID int `json:"cabin_type_id" gorm:"Column:CabinTypeID"`
	CabinType   *CabinType

	FirstName         string  `json:"first_name" gorm:"Column:Firstname"`
	LastName          string  `json:"last_name" gorm:"Column:Lastname"`
	Email             *string `json:"email" gorm:"Column:Email"`
	Phone             string  `json:"phone" gorm:"Column:Phone"`
	PassportNumber    string  `json:"passport_number" gorm:"Column:PassportNumber"`
	PassportCountryID int     `json:"passport_country_id" gorm:"Column:PassportCountryID"`
	BookingReference  string  `json:"booking_reference" gorm:"Column:BookingReference"`
	Confirmed         bool    `json:"confirmed" gorm:"Column:Confirmed"`
}

func (t *Ticket) TableName() string {
	return "tickets"
}
