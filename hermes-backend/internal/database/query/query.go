package database

type (
	Join struct {
		other   string
		thisId  string
		otherId string
	}
	// TODO: Introduce errors for invalid queries
	QueryBuilder interface {
		Generate() string
	}
)
