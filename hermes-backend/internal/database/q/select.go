package q

import (
	"fmt"
	"strings"
	"woojiahao.com/hermes/internal"
)

type (
	SelectOrder string
	SelectQuery struct {
		table     string
		columns   []string
		condition string
		joins     []Join
		order     []string
	}
)

const (
	ASC  SelectOrder = "ASC"
	DESC SelectOrder = "DESC"
)

func From(name string) *SelectQuery {
	return &SelectQuery{table: name}
}

func (q *SelectQuery) Select(columns ...string) *SelectQuery {
	q.columns = columns
	return q
}

func (q *SelectQuery) Join(other, thisId, otherId string) *SelectQuery {
	join := Join{other, thisId, otherId}
	q.joins = append(q.joins, join)
	return q
}

func (q *SelectQuery) Where(condition string) *SelectQuery {
	q.condition = condition
	return q
}

func (q *SelectQuery) Order(column string, order SelectOrder) *SelectQuery {
	q.order = append(q.order, fmt.Sprintf("%s %s", column, order))
	return q
}

func (q *SelectQuery) Generate() string {
	queryLines := []string{
		fmt.Sprintf("SELECT %s", strings.Join(q.columns, ", ")),
		fmt.Sprintf("FROM %s", q.table),
	}
	if len(q.joins) > 0 {
		joins := internal.Map(q.joins, func(join Join) string {
			return fmt.Sprintf("\tINNER JOIN %s ON %s = %s", join.other, join.otherId, join.thisId)
		})
		queryLines = append(queryLines, joins...)
	}

	if q.condition != "" {
		queryLines = append(queryLines, fmt.Sprintf("WHERE %s", q.condition))
	}

	if len(q.order) > 0 {
		queryLines = append(queryLines, fmt.Sprintf("ORDER BY %s", strings.Join(q.order, ", ")))
	}

	return strings.Join(queryLines, "\n")
}

/*
From("threads").Delete().Where("")
*/
