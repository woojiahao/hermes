package q

import (
	"fmt"
	"strings"
	"woojiahao.com/hermes/internal"
)

type InsertQuery struct {
	table              string
	columns            []string
	values             [][]any
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
	i.values = append(i.values, values)
	return i
}

func (i *InsertQuery) Returning(columns ...string) *InsertQuery {
	i.returns = columns
	return i
}

func (i *InsertQuery) OnConflict(columns ...string) *InsertQuery {
	i.onConflict = true
	i.onConflictColumns = columns
	return i
}

func (i *InsertQuery) DoNothing() *InsertQuery {
	i.onConflictBehavior = "DO NOTHING"
	return i
}

func (i *InsertQuery) DoUpdate(set map[string]any) *InsertQuery {
	var setList []string
	for k, v := range set {
		setList = append(setList, fmt.Sprintf("%s = %v", k, v))
	}
	i.onConflictBehavior = fmt.Sprintf("DO UPDATE SET %s", strings.Join(setList, ", "))
	return i
}

func (i *InsertQuery) Generate() string {
	topLine := fmt.Sprintf("INSERT INTO %s", i.table)
	if len(i.columns) > 0 {
		topLine = fmt.Sprintf("INSERT INTO %s (%s)", i.table, strings.Join(i.columns, ", "))
	}
	insertLines := []string{
		topLine,
		"VALUES",
	}

	var values []string
	for _, v := range i.values {
		line := fmt.Sprintf("\t(%s)", strings.Join(internal.Map(v, func(a any) string {
			return fmt.Sprintf("%v", a)
		}), ", "))
		values = append(values, line)
	}
	insertLines = append(insertLines, strings.Join(values, ",\n"))

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
