package internal

import "fmt"

// TODO: Add ways to differentiate errors for the server
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

// Handles when a result set should exactly one result
func ExactlyOneResultError[T any](arr []T) error {
	if len(arr) == 0 {
		return &ServerError{"unable to find intended result", nil}
	}

	if len(arr) > 1 {
		return &ServerError{"result returned should only be one", nil}
	}

	return nil
}
