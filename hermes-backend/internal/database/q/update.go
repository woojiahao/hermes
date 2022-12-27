package q

import (
	"fmt"
	"strings"
)

type UpdateQuery struct {
	table     string
	sets      map[string]any
	condition string
	returns   []string
}

func Update(table string) *UpdateQuery {
	return &UpdateQuery{table: table, sets: make(map[string]any)}
}

func (u *UpdateQuery) Set(column string, value any) *UpdateQuery {
	u.sets[column] = value
	return u
}

func (u *UpdateQuery) Where(condition string) *UpdateQuery {
	u.condition = condition
	return u
}

func (u *UpdateQuery) Returning(columns ...string) *UpdateQuery {
	u.returns = columns
	return u
}

func (u *UpdateQuery) Generate() string {
	updateLines := []string{
		fmt.Sprintf("UPDATE %s", u.table),
	}

	var setLines []string
	for col, val := range u.sets {
		setLines = append(setLines, fmt.Sprintf("%s = %v", col, val))
	}
	updateLines = append(updateLines, fmt.Sprintf("SET %s", strings.Join(setLines, ", ")))
	updateLines = append(updateLines, fmt.Sprintf("WHERE %s", u.condition))
	if len(u.returns) > 0 {
		updateLines = append(updateLines, fmt.Sprintf("RETURNING %s", strings.Join(u.returns, ", ")))
	}

	return strings.Join(updateLines, "\n")
}
