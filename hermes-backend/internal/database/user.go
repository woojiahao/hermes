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
		"INSERT INTO \"user\"(username, email, password_hash) VALUES ($1, $2, $3) RETURNING *",
		generate_params(username, email, string(hash)),
		func(r *sql.Rows) (User, error) {
			var user User
			err := r.Scan(&user.Username, &user.Email, &user.PasswordHash)
			return user, err
		},
	)

	if err != nil {
		return dummyUser, err
	}

	if len(users) != 1 {
		return dummyUser, errors.New("Invalid INSERT on user")
	}

	return users[0], nil
}

func (d *Database) GetUser(username string) (User, error) {
	return User{"", "", ""}, nil
}

func (d *Database) GetUsers() ([]User, error) {
	return make([]User, 0), nil
}
