package database

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username     string
	Email        string
	PasswordHash string
}

var dummyUser = User{"", "", ""}

func parseUserRows(rows *sql.Rows) (User, error) {
	var user User
	err := rows.Scan(&user.Username, &user.Email, &user.PasswordHash)
	return user, err
}

func (d *Database) CreateUser(
	username string,
	email string,
	password string,
) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return dummyUser, err
	}

	users, err := query(
		d,
		"INSERT INTO \"user\"(username, email, password_hash) VALUES ($1, $2, $3) RETURNING *;",
		generate_params(username, email, string(hash)),
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

func (d *Database) GetUser(username string) (User, error) {
	users, err := query(
		d,
		"SELECT username, email, password_hash FROM \"user\" WHERE username = $1;",
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
		"SELECT username, email, password_hash FROM \"user\";",
		generate_params(),
		parseUserRows,
	)
}
