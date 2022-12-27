package database

import (
	"fmt"
	"testing"
	. "woojiahao.com/hermes/internal/database/query"
)

func TestShouldGenerateSelectQuery(t *testing.T) {
	q := From("customer").
		Select("id", "name").
		Join("order", "order_id", "id").
		Where("quantity > 15").
		Order("is_pinned", DESC).
		Order("created_at", ASC).
		Generate()
	fmt.Println(q)
}

func TestShouldGenerateDeleteQuery(t *testing.T) {
	q := Delete("customer").Where("quantity < 4").Generate()
	fmt.Println(q)
}

func TestShouldGenerateUpdateQuery(t *testing.T) {
	q := Update("thread").
		Where(And("thread.id = $5", "thread.created_by = $6")).
		Set("title", "$1").
		Set(`"content"`, "$2").
		Set("is_published", "$3").
		Set("is_open", "$4").
		Set("updated_at", "NOW()").
		Returning("*").
		Generate()
	fmt.Println(q)
}

func TestShouldGenerateInsertQuery(t *testing.T) {
	q := Insert("thread").Columns("title", `"content"`, "created_by").Values("$1", "$2", "$3").Returning("*").Generate()
	fmt.Println(q)
}
