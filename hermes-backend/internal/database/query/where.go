package database

import "fmt"

func And(left, right string) string {
	return fmt.Sprintf("(%s AND %s)", left, right)
}

func Or(left, right string) string {
	return fmt.Sprintf("(%s OR %s)", left, right)
}

func IsNull(column string) string {
	return fmt.Sprintf("%s IS NULL", column)
}
