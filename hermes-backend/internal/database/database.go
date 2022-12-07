package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DatabaseConfiguration struct {
	Username string
	Password string
	Host     string
	Name     string
	Port     int
}

// Loads the database configuration from .env
// TODO: Support development vs production .env
func LoadConfiguration() *DatabaseConfiguration {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	name := os.Getenv("DATABASE_NAME")
	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		log.Fatal("Invalid DATABASE_PORT in .env")
	}

	return &DatabaseConfiguration{username, password, host, name, port}
}

type Database struct {
	configuration *DatabaseConfiguration
	db            *sql.DB
}

func Initialize(c *DatabaseConfiguration) *Database {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Name,
	)

	db, err := sql.Open(
		"postgres",
		connStr,
	)
	if err != nil {
		log.Fatal("Error loading database from given configuration " + err.Error())
	}

	return &Database{c, db}
}

func query[T any](
	db *Database,
	query string,
	params []any,
	perRow func(*sql.Rows) (T, error),
) ([]T, error) {
	rows, err := db.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []T

	for rows.Next() {
		item, err := perRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
