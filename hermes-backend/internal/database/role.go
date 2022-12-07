package database

import (
	"database/sql"
)

func (d *Database) GetRoles() ([]Role, error) {
	return query(
		d,
		"SELECT title, permissions FROM role;",
		make([]any, 0),
		func(rows *sql.Rows) (Role, error) {
			var role Role
			err := rows.Scan(&role.Title, &role.Permissions)
			return role, err
		},
	)
}
