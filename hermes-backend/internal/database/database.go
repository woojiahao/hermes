package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type perRow[T any] func(*sql.Rows) (T, error)

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

func parseResults[T any](rows *sql.Rows, fn perRow[T]) ([]T, error) {
	var items []T

	for rows.Next() {
		item, err := fn(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func query[T any](db *Database, query string, params []any, fn perRow[T]) ([]T, error) {
	rows, err := db.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseResults(rows, fn)
}

func transaction[T any](db *Database, fn func(*sql.Tx) (T, error)) (T, error) {
	tx, err := db.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return *new(T), err
	}

	defer tx.Rollback()

	return fn(tx)
}

func transactionQuery[T any](tx *sql.Tx, query string, params []any, fn perRow[T]) ([]T, error) {
	rows, err := tx.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseResults[T](rows, fn)
}

func generate_params(values ...any) []any {
	return values
}