package q

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

func Eq(column string, value any) string {
	return fmt.Sprintf("%s = %v", column, value)
}

func Exists(query string) string {
	return fmt.Sprintf("EXISTS(%s)", query)
}
