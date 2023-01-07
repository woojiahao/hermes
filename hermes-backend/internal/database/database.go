package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

//go:embed sql/create.sql
var createSqlScript string

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

func (d *Database) CreateTablesIfNotExists() {
	log.Println("Creating database tables")
	_, err := d.db.ExecContext(context.TODO(), createSqlScript)
	if err != nil {
		log.Fatalf("Error created tables for hermes %s", err)
	}
}
