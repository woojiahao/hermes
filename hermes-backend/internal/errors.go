package internal

import "fmt"

// Database error that occurs during query
type DatabaseError struct {
	base   error
	custom string
}

func (dbe *DatabaseError) Error() string {
	return fmt.Sprintf("database error occurred: %s, base: %s", dbe.custom, dbe.base.Error())
}

// Generic server error
type ServerError struct {
	base   error
	custom string
}

func (se *ServerError) Error() string {
	return fmt.Sprintf("generic server error occurred: %s, base: %s", se.custom, se.base.Error())
}
