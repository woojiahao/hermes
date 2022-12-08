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
		Email        string
		PasswordHash string
		Role
	}
)

const (
	ADMIN Role = "ADMIN"
	USER  Role = "USER"
)

var dummyUser = User{"", "", "", "", ""}

func parseUserRows(rows *sql.Rows) (User, error) {
	var user User
	err := rows.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)
	return user, err
}

func (d *Database) CreateUser(
	username string,
	email string,
	password string,
	role Role,
) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return dummyUser, &i.ServerError{"failed to generate hash for password", err}
	}

	users, err := query(
		d,
		`
			INSERT INTO "user"(username, email, password_hash, role)
			VALUES ($1, $2, $3, $4::text::"role")
			RETURNING *;
		`,
		generate_params(username, email, string(hash), role),
		parseUserRows,
	)

	if err != nil {
		return dummyUser, &i.DatabaseError{"failed to insert new user, reason: conflicting username or email", err}
	}

	err = i.ExactlyOneResultError(users)
	if err != nil {
		return dummyUser, err
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
		generate_params(userId),
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
		generate_params(username),
		parseUserRows,
	)

	if err != nil {
		return dummyUser, &i.DatabaseError{Custom: "failed to retrieve user by username", Base: err}
	}

	err = i.ExactlyOneResultError(users)
	if err != nil {
		return dummyUser, err
	}

	return users[0], nil
}

func (d *Database) GetUsers() ([]User, error) {
	users, err := query(
		d,
		"SELECT * FROM user;",
		generate_params(),
		parseUserRows,
	)

	if err != nil {
		return make([]User, 0), &i.DatabaseError{Custom: "failed to retrieve all users", Base: err}
	}

	return users, nil
}
