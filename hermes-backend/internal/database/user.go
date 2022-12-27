package database

import (
	"database/sql"
	"log"
	"woojiahao.com/hermes/internal"
	"woojiahao.com/hermes/internal/database/q"

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
		log.Printf("Failed to generate hash for password due to %s", err)
		return dummyUser, InternalError
	}

	user, err := singleQuery(
		d,
		q.Insert(`"user"`).
			Columns("username", "password_hash", "role").
			Values("$1", "$2", `$3::text::"role"`).
			Returning(q.ALL).
			Generate(),
		generateParams(username, string(hash), role),
		parseUserRows,
	)
	if err != nil {
		return dummyUser, err
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
		q.From(`"user"`).
			Select("*").
			Where(q.Eq(internal.ThisOrThat("username", "id", !isUsername), "$1")).
			Generate(),
		generateParams(str),
		parseUserRows,
	)
	if err != nil {
		return dummyUser, err
	}

	return user, nil
}

func (d *Database) GetUsers() ([]User, error) {
	users, err := query(
		d,
		q.From(`"user"`).Select(q.ALL).Generate(),
		generateParams(),
		parseUserRows,
	)
	if err != nil {
		return nil, err
	}

	return users, nil
}
