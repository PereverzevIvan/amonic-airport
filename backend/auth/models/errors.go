package models

// Define a custom error type based on string
type Error string

// Implement the Error() method for the custom error type
func (e Error) Error() string {
	return string(e)
}

// Declare constant errors
const (
	ErrDuplicatedEmail    = Error("duplicated email error")
	ErrFKOfficeIDNotFound = Error("office with this id not found")
)
