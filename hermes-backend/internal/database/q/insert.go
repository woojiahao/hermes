package q

import (
	"fmt"
	"strings"
	"woojiahao.com/hermes/internal"
)

type InsertQuery struct {
	table   string
	columns []string
	values  []any
	returns []string
}

func Insert(table string) *InsertQuery {
	return &InsertQuery{table: table}
}

func (i *InsertQuery) Columns(columns ...string) *InsertQuery {
	i.columns = columns
	return i
}

func (i *InsertQuery) Values(values ...any) *InsertQuery {
	i.values = values
	return i
}

func (i *InsertQuery) Returning(columns ...string) *InsertQuery {
	i.returns = columns
	return i
}

func (i *InsertQuery) Generate() string {
	insertLines := []string{
		fmt.Sprintf("INSERT INTO %s (%s)", i.table, strings.Join(i.columns, ", ")),
		fmt.Sprintf("VALUES (%s)", strings.Join(internal.Map(i.values, func(a any) string {
			return fmt.Sprintf("%v", a)
		}), ", ")),
	}

	if len(i.returns) > 0 {
		insertLines = append(insertLines, fmt.Sprintf("RETURNING %s", strings.Join(i.returns, ", ")))
	}

	return strings.Join(insertLines, "\n")
}
