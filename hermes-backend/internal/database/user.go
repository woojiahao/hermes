package database

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
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
		return dummyUser, err
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
		return dummyUser, err
	}

	if len(users) != 1 {
		return dummyUser, errors.New("invalid INSERT on user")
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
		return dummyUser, err
	}

	switch len(users) {
	case 0:
		return dummyUser, errors.New("unable to find user")
	case 1:
		return users[0], nil
	default:
		return dummyUser, errors.New("user should be unique by id")
	}
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
		return dummyUser, err
	}

	switch len(users) {
	case 0:
		return dummyUser, errors.New("unable to find user")
	case 1:
		return users[0], nil
	default:
		return dummyUser, errors.New("user should be unique by username")
	}
}

func (d *Database) GetUsers() ([]User, error) {
	return query(
		d,
		"SELECT * FROM user;",
		generate_params(),
		parseUserRows,
	)
}
