package database

import (
	"context"
	"database/sql"
	"log"
	"woojiahao.com/hermes/internal/database/q"
)

type perRow[T any] func(*sql.Rows) (T, error)

func doNothing(_ *sql.Rows) (string, error) {
	return "", nil
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

func singleQuery[T any](db *Database, q q.QueryBuilder, params []any, fn perRow[T]) (T, error) {
	results, err := query(db, q, params, fn)
	if err != nil {
		// The error is internal error, so just propagate it
		return *new(T), err
	}

	if len(results) < 1 {
		return *new(T), NotFoundError
	}

	return results[0], nil
}

func query[T any](db *Database, query q.QueryBuilder, params []any, fn perRow[T]) ([]T, error) {
	rows, err := db.db.QueryContext(context.TODO(), query.Generate(), params...)
	if err != nil {
		log.Printf("Internal database error occurred when querying for data due to %s", err)
		return nil, InternalError
	}
	defer rows.Close()

	return parseResults(rows, fn)
}

func transaction[T any](db *Database, fn func(*sql.Tx) (T, error)) (T, error) {
	tx, err := db.db.BeginTx(context.TODO(), nil)
	if err != nil {
		log.Printf("Internal database error occurred when setting up transaction due to %s", err)
		return *new(T), InternalError
	}

	defer tx.Rollback()

	result, err := fn(tx)
	if err != nil {
		// These errors are either InternalError or NotFoundError
		return *new(T), err
	}

	tx.Commit()

	return result, nil
}

func transactionSingleQuery[T any](tx *sql.Tx, q q.QueryBuilder, params []any, fn perRow[T]) (T, error) {
	results, err := transactionQuery(tx, q, params, fn)
	if err != nil {
		// Error will be InternalError
		return *new(T), err
	}

	if len(results) < 1 {
		return *new(T), NotFoundError
	}

	return results[0], nil
}

func transactionQuery[T any](tx *sql.Tx, query q.QueryBuilder, params []any, fn perRow[T]) ([]T, error) {
	rows, err := tx.QueryContext(context.TODO(), query.Generate(), params...)
	if err != nil {
		log.Printf("Internal database error occurred when querying for data in transaction due to %s", err)
		return nil, InternalError
	}
	defer rows.Close()

	return parseResults(rows, fn)
}

func generateParams(values ...any) []any {
	return values
}
