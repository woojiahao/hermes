package database

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
	i "woojiahao.com/hermes/internal"
)

type (
	Role string

	User struct {
		Id           string
		Username     string
		PasswordHash string
		Role
	}
)

const (
	ADMIN Role = "ADMIN"
	USER  Role = "USER"
)

var dummyUser User

func parseUserRows(rows *sql.Rows) (User, error) {
	var user User
	err := rows.Scan(
		&user.Id,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
	)
	return user, err
}

func (d *Database) CreateUser(username string, password string, role Role) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return dummyUser, &i.ServerError{Custom: "failed to generate hash for password", Base: err}
	}

	users, err := query(
		d,
		`
			INSERT INTO "user"(username, password_hash, role)
			VALUES ($1, $2, $3::text::"role")
			RETURNING *;
		`,
		generateParams(username, string(hash), role),
		parseUserRows,
	)

	if err != nil {
		return dummyUser, &i.DatabaseError{
			Custom: "failed to insert new user, reason: conflicting username or email",
			Base:   err,
			Short:  "Username is already taken",
		}
	}

	return users[0], nil
}

// TODO: Add promote/demote user function

// TODO: Refactor this
func (d *Database) GetUserById(userId string) (User, error) {
	users, err := query(
		d,
		`
			SELECT *
			FROM "user"
			WHERE id = $1;
		`,
		generateParams(userId),
		parseUserRows,
	)

	if err != nil {
		return dummyUser, &i.DatabaseError{Custom: "failed to retrieve user by id", Base: err}
	}

	err = i.ExactlyOneResultError(users)
	if err != nil {
		return dummyUser, err
	}

	return users[0], nil
}

func (d *Database) GetUser(username string) (User, error) {
	users, err := query(
		d,
		`
			SELECT *
			FROM "user"
			WHERE username = $1;
		`,
		generateParams(username),
		parseUserRows,
	)
	if err != nil || len(users) == 0 {
		return dummyUser, &i.DatabaseError{Custom: "failed to retrieve user by username", Base: err}
	}

	return users[0], nil
}

func (d *Database) GetUsers() ([]User, error) {
	users, err := query(
		d,
		`SELECT * FROM "user";`,
		generateParams(),
		parseUserRows,
	)

	if err != nil {
		return make([]User, 0), &i.DatabaseError{Custom: "failed to retrieve all users", Base: err}
	}

	return users, nil
}
