package q

import (
	"fmt"
	"strings"
)

type DeleteQuery struct {
	table     string
	condition string
}

func Delete(table string) *DeleteQuery {
	return &DeleteQuery{table: table}
}

func (d *DeleteQuery) Where(condition string) *DeleteQuery {
	d.condition = strings.Trim(condition, " ")
	return d
}

func (d *DeleteQuery) Generate() string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s", d.table, d.condition)
}
