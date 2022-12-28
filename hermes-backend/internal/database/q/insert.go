package q

import (
	"fmt"
	"strings"
	"woojiahao.com/hermes/internal"
)

type InsertQuery struct {
	table              string
	columns            []string
	values             []any
	onConflict         bool
	onConflictColumns  []string
	onConflictBehavior string
	returns            []string
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

func (i *InsertQuery) OnConflict(columns ...string) *InsertQuery {
	i.onConflictColumns = columns
	return i
}

func (i *InsertQuery) DoNothing() *InsertQuery {
	i.onConflictBehavior = "DO NOTHING"
	return i
}

func (i *InsertQuery) Generate() string {
	topLine := fmt.Sprintf("INSERT INTO %s", i.table)
	if len(i.columns) > 0 {
		topLine = fmt.Sprintf("INSERT INTO %s (%s)", i.table, strings.Join(i.columns, ", "))
	}
	insertLines := []string{
		topLine,
		fmt.Sprintf("VALUES (%s)", strings.Join(internal.Map(i.values, func(a any) string {
			return fmt.Sprintf("%v", a)
		}), ", ")),
	}

	if i.onConflict {
		line := internal.ThisOrThat(
			"ON CONFLICT",
			fmt.Sprintf("ON CONFLICT (%s)", strings.Join(i.onConflictColumns, ", ")),
			len(i.onConflictColumns) == 0,
		)
		insertLines = append(insertLines, line)
		insertLines = append(insertLines, i.onConflictBehavior)
	}

	if len(i.returns) > 0 {
		insertLines = append(insertLines, fmt.Sprintf("RETURNING %s", strings.Join(i.returns, ", ")))
	}

	return strings.Join(insertLines, "\n")
}
