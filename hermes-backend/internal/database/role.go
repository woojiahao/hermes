package database

func (d *Database) GetRoles() ([]Role, error) {
	rows, err := d.db.Query("SELECT title, permissions FROM role;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role

	for rows.Next() {
		var role Role

		err := rows.Scan(&role.Title, &role.Permissions)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
