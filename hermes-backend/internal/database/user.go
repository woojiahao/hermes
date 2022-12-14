package database

import (
	"database/sql"
	"log"
	"woojiahao.com/hermes/internal"
	. "woojiahao.com/hermes/internal/database/q"

	"golang.org/x/crypto/bcrypt"
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
		log.Printf("Failed to generate hash for password due to %s", err)
		return User{}, InternalError
	}

	user, err := singleQuery(
		d,
		Insert(`"user"`).
			Columns("username", "password_hash", "role").
			Values(P1, P2, `$3::text::"role"`).
			Returning(ALL),
		generateParams(username, string(hash), role),
		parseUserRows,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (d *Database) GetUserById(userId string) (User, error) {
	return d.getUser(userId, false)
}

func (d *Database) GetUser(username string) (User, error) {
	return d.getUser(username, true)
}

func (d *Database) getUser(str string, isUsername bool) (User, error) {
	user, err := singleQuery(
		d,
		From(`"user"`).
			Select(ALL).
			Where(Eq(internal.ThisOrThat("username", "id", isUsername), P1)),
		generateParams(str),
		parseUserRows,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (d *Database) GetUsers() ([]User, error) {
	users, err := query(
		d,
		From(`"user"`).Select(ALL),
		generateParams(),
		parseUserRows,
	)
	if err != nil {
		return nil, err
	}

	return users, nil
}
