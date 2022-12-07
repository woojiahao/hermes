package database

import (
	"database/sql"
	"errors"
	"time"
)

type Tag struct {
	Id      string
	Content string
	HexCode string
}

// TODO: Add list of tags
type Thread struct {
	Id          string
	IsPublished bool
	IsOpen      bool
	Title       string
	Content     string
	CreatedAt   *time.Time
	CreatedBy   string
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	DeletedBy   sql.NullString
}

var dummyThread = Thread{"", false, false, "", "", nil, "", nil, nil, sql.NullString{}}

func parseThreadRows(rows *sql.Rows) (Thread, error) {
	var thread Thread
	err := rows.Scan(
		&thread.Id,
		&thread.IsPublished,
		&thread.IsOpen,
		&thread.Title,
		&thread.Content,
		&thread.CreatedAt,
		&thread.CreatedBy,
		&thread.UpdatedAt,
		&thread.DeletedAt,
		&thread.DeletedBy,
	)
	return thread, err
}

func (d *Database) CreateThread(userId, title, content string) (Thread, error) {
	threads, err := query(
		d,
		`
			INSERT INTO thread (title, "content", created_by)
			VALUES ($1, $2, $3)
			RETURNING *;
		`,
		generate_params(title, content, userId),
		parseThreadRows,
	)

	if err != nil {
		return dummyThread, err
	}

	return threads[0], err
}

func (d *Database) GetUserThreads(userId string) ([]Thread, error) {
	return query(
		d,
		"SELECT * FROM thread INNER JOIN \"user\" ON thread.created_by = \"user\".id WHERE \"user\".id = $1",
		generate_params(userId),
		parseThreadRows,
	)
}

func (d *Database) GetThreadById(threadId string) (Thread, error) {
	threads, err := query(
		d,
		"SELECT * FROM thread WHERE thread.id = $1",
		generate_params(threadId),
		parseThreadRows,
	)

	if err != nil {
		return dummyThread, err
	}

	switch len(threads) {
	case 0:
		return dummyThread, errors.New("unable to find thread")
	case 1:
		return threads[0], nil
	default:
		return dummyThread, errors.New("thread should be unique by id")
	}
}

func (d *Database) GetThreads() ([]Thread, error) {
	return query(
		d,
		"SELECT * FROM thread",
		generate_params(),
		parseThreadRows,
	)
}

func (d *Database) DeleteThread(userId, threadId string) (Thread, error) {
	threads, err := query(
		d,
		`
			UPDATE thread
			SET deleted_by = $1, deleted_at = NOW()
			WHERE thread.id = $2
				AND (
					thread.created_by = $1 OR
					EXISTS (SELECT * FROM "user" WHERE "user".id = $1 AND "user".role = 'ADMIN')
				)
			RETURNING *
		`,
		generate_params(userId, threadId),
		parseThreadRows,
	)

	if err != nil {
		return dummyThread, err
	}

	switch len(threads) {
	case 0:
		return dummyThread, errors.New("unable to update thread")
	case 1:
		return threads[0], nil
	default:
		return dummyThread, errors.New("only one thread should have been updated")
	}
}

func (d *Database) EditThread(
	userId,
	threadId,
	title,
	content string,
	isPublished,
	isOpen bool,
) (Thread, error) {
	threads, err := query(
		d,
		`
			UPDATE thread
			SET title = $1, "content" = $2, is_published = $3, is_open = $4, updated_at = NOW()
			WHERE thread.id = $5 AND thread.created_by = $6
			RETURNING *
		`,
		generate_params(title, content, isPublished, isOpen, threadId, userId),
		parseThreadRows,
	)

	if err != nil {
		return dummyThread, err
	}

	switch len(threads) {
	case 0:
		return dummyThread, errors.New("unable to update thread")
	case 1:
		return threads[0], nil
	default:
		return dummyThread, errors.New("only one thread should have been updated")
	}
}
