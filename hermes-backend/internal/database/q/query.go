package q

import "fmt"

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

const (
	ALL = "*"
	NOW = "NOW()"
	P1  = "$1"
	P2  = "$2"
	P3  = "$3"
	P4  = "$4"
	P5  = "$5"
	P6  = "$6"
	P7  = "$7"
)

func Sub(subQuery QueryBuilder, name string) string {
	return fmt.Sprintf("(%s) %s", subQuery.Generate(), name)
}

func Coalaesce(intended string, other any) string {
	return fmt.Sprintf("COALESCE(%s, %v)", intended, other)
}
