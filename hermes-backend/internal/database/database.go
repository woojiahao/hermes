package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Database struct {
	configuration *Configuration
	db            *sql.DB
}

func Initialize(c *Configuration) *Database {
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
		log.Fatalf("Error loading database from given configuration %s", err)
	}

	return &Database{c, db}
}
