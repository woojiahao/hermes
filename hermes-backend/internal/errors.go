package internal

import "fmt"

// Database error that occurs during query
type DatabaseError struct {
	Custom string
	Base   error
}

func (dbe *DatabaseError) Error() string {
	return fmt.Sprintf("database error occurred: %s, base: %s", dbe.Custom, dbe.Base.Error())
}

// Generic server error
type ServerError struct {
	Custom string
	Base   error
}

func (se *ServerError) Error() string {
	base := ""
	if se.Base != nil {
		base = fmt.Sprintf(", base: %s", se.Base.Error())
	}
	return fmt.Sprintf("generic server error occurred: %s%s", se.Custom, base)
}
