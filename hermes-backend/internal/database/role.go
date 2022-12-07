package database

import (
	"database/sql"
)

type Role struct {
	Title       string
	Permissions string
}

func (d *Database) GetRoles() ([]Role, error) {
	return query(
		d,
		"SELECT title, permissions FROM role;",
		generate_params(),
		func(rows *sql.Rows) (Role, error) {
			var role Role
			err := rows.Scan(&role.Title, &role.Permissions)
			return role, err
		},
	)
}
