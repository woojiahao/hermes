package q

import (
	"fmt"
	"strings"
	"woojiahao.com/hermes/internal"
)

type (
	SelectOrder string
	SelectQuery struct {
		table      string
		columns    []string
		condition  string
		innerJoins []Join
		leftJoins  []Join
		rightJoins []Join
		order      []string
		group      []string
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

func (q *SelectQuery) InnerJoin(other, thisId, otherId string) *SelectQuery {
	join := Join{other, thisId, otherId}
	q.innerJoins = append(q.innerJoins, join)
	return q
}

func (q *SelectQuery) LeftJoin(other, thisId, otherId string) *SelectQuery {
	join := Join{other, thisId, otherId}
	q.leftJoins = append(q.leftJoins, join)
	return q
}

func (q *SelectQuery) RightJoin(other, thisId, otherId string) *SelectQuery {
	join := Join{other, thisId, otherId}
	q.rightJoins = append(q.rightJoins, join)
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

func (q *SelectQuery) Group(columns ...string) *SelectQuery {
	q.group = columns
	return q
}

func (q *SelectQuery) Generate() string {
	queryLines := []string{
		fmt.Sprintf("SELECT %s", strings.Join(q.columns, ", ")),
		fmt.Sprintf("FROM %s", q.table),
	}
	if len(q.innerJoins) > 0 {
		joins := internal.Map(q.innerJoins, func(join Join) string {
			return fmt.Sprintf("\tINNER JOIN %s ON %s = %s", join.other, join.otherId, join.thisId)
		})
		queryLines = append(queryLines, joins...)
	}

	if len(q.leftJoins) > 0 {
		joins := internal.Map(q.leftJoins, func(join Join) string {
			return fmt.Sprintf("\tLEFT JOIN %s ON %s = %s", join.other, join.otherId, join.thisId)
		})
		queryLines = append(queryLines, joins...)
	}

	if len(q.rightJoins) > 0 {
		joins := internal.Map(q.rightJoins, func(join Join) string {
			return fmt.Sprintf("\tRIGHT JOIN %s ON %s = %s", join.other, join.otherId, join.thisId)
		})
		queryLines = append(queryLines, joins...)
	}

	if q.condition != "" {
		queryLines = append(queryLines, fmt.Sprintf("WHERE %s", q.condition))
	}

	if len(q.group) > 0 {
		queryLines = append(queryLines, fmt.Sprintf("GROUP BY %s", strings.Join(q.group, ", ")))
	}

	if len(q.order) > 0 {
		queryLines = append(queryLines, fmt.Sprintf("ORDER BY %s", strings.Join(q.order, ", ")))
	}

	return strings.Join(queryLines, "\n")
}

/*
From("threads").Delete().Where("")
*/
