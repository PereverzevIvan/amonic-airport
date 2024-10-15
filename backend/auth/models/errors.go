package models

// Define a custom error type based on string
type Error string

// Implement the Error() method for the custom error type
func (e Error) Error() string {
	return string(e)
}

// Declare constant errors
const (
	ErrNotFound                     = Error("not found error")
	ErrUnique                       = Error("unique constraint error")
	ErrFK                           = Error("foreign key error")
	ErrDuplicatedEmail              = Error("duplicated email error")
	ErrFKOfficeIDNotFound           = Error("office with this id not found")
	ErrCSVMissingFields             = Error("missing fields in csv file")
	ErrNoTicketsAvailable           = Error("no tickets available")
	ErrCantEditAmenitiesTimeExpired = Error("can't edit amenities time expired")
)
