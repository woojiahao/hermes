package database

import (
	"database/sql"
	"fmt"
	"time"

	i "woojiahao.com/hermes/internal"
)

type Tag struct {
	Id      string
	Content string
	HexCode string
}

func NewTag(content, hexCode string) Tag {
	return Tag{Content: content, HexCode: hexCode}
}

func parseTagRows(rows *sql.Rows) (Tag, error) {
	var tag Tag
	err := rows.Scan(&tag.Id, &tag.Content, &tag.HexCode)
	return tag, err
}

// Opt to use two queries in a transaction to populate the Thread and its associated tags to avoid unnecessarily complex
// queries or parsing
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
	Tags        []Tag
}

var dummyThread Thread

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

func (d *Database) CreateThread(userId, title, content string, tags []Tag) (Thread, error) {
	var tag_contents []string
	for _, tag := range tags {
		tag_contents = append(tag_contents, tag.Content)
	}

	if i.HasDuplicates(tag_contents) {
		return dummyThread, &i.ServerError{Custom: "tags cannot have same name", Base: nil}
	}

	return transaction(d, func(tx *sql.Tx) (Thread, error) {
		threads, err := transactionQuery(
			tx,
			`
				INSERT INTO thread (title, "content", created_by)
				VALUES ($1, $2, $3)
				RETURNING *;
			`,
			generateParams(title, content, userId),
			parseThreadRows,
		)
		if err != nil {
			return dummyThread, &i.DatabaseError{Custom: "failed to create new thread", Base: err}
		}

		thread := threads[0]

		var db_tags []Tag

		// Create or retrieve all of the tags from the database
		// Once retrieved, attach the tag to the thread
		for _, tag := range tags {
			ts, err := transactionQuery(
				tx,
				`
					WITH q AS(
						INSERT INTO tag ("content", hex_code)
						VALUES ($1, $2)
						ON CONFLICT("content")
						DO NOTHING
						RETURNING *
					)
					SELECT * FROM q
					UNION
					SELECT * FROM tag WHERE tag."content" = $1;
				`,
				generateParams(tag.Content, tag.HexCode),
				parseTagRows,
			)

			if err != nil {
				return dummyThread, &i.DatabaseError{Custom: "failed to create/retrieve new tag", Base: err}
			}

			db_tags = append(db_tags, ts[0])

			_, err = transactionQuery(
				tx,
				`
					INSERT INTO thread_tag
					VALUES ($1, $2)
					RETURNING *;
				`,
				generateParams(thread.Id, ts[0].Id),
				func(r *sql.Rows) (string, error) {
					return "", nil // this perRow fn does not parse the results
				},
			)

			if err != nil {
				return dummyThread, &i.DatabaseError{Custom: "failed to link thread with tag", Base: err}
			}
		}

		thread.Tags = db_tags

		tx.Commit()

		return thread, nil
	})
}

func (d *Database) GetUserThreads(userId string) ([]Thread, error) {
	return query(
		d,
		"SELECT * FROM thread INNER JOIN \"user\" ON thread.created_by = \"user\".id WHERE \"user\".id = $1",
		generateParams(userId),
		parseThreadRows,
	)
}

// TODO: Support loading tags
func (d *Database) GetThreadById(threadId string) (Thread, error) {
	threads, err := query(
		d,
		"SELECT * FROM thread WHERE thread.id = $1",
		generateParams(threadId),
		parseThreadRows,
	)

	if err != nil {
		return dummyThread, &i.DatabaseError{Custom: "failed to retrieve thread by id", Base: err}
	}

	err = i.ExactlyOneResultError(threads)
	if err != nil {
		return dummyThread, err
	}

	return threads[0], nil
}

func (d *Database) GetThreads() ([]Thread, error) {
	return transaction(d, func(tx *sql.Tx) ([]Thread, error) {
		threads, err := transactionQuery(
			tx,
			"SELECT * FROM thread",
			generateParams(),
			parseThreadRows,
		)

		if err != nil {
			return nil, &i.DatabaseError{Custom: "failed to retrieve all threads", Base: err}
		}

		threadTags := make(map[string][]Tag)

		_, err = transactionQuery(
			tx,
			`
				SELECT t.id, tag.id, tag."content", tag.hex_code
				FROM tag
								INNER JOIN thread_tag tt on tag.id = tt.tag_id
								INNER JOIN thread t on t.id = tt.thread_id;
			`,
			generateParams(),
			func(r *sql.Rows) (string, error) {
				var threadId string
				var tag Tag
				err := r.Scan(&threadId, &tag.Id, &tag.Content, &tag.HexCode)

				if _, found := threadTags[threadId]; found {
					threadTags[threadId] = append(threadTags[threadId], tag)
				} else {
					threadTags[threadId] = []Tag{tag}
				}

				return "", err
			},
		)

		if err != nil {
			return nil, &i.DatabaseError{Custom: "failed to retrieve all tags related to all threads", Base: err}
		}

		for _, thread := range threads {
			thread.Tags = threadTags[thread.Id]
		}

		return threads, nil
	})
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
		generateParams(userId, threadId),
		parseThreadRows,
	)

	if err != nil {
		return dummyThread, &i.DatabaseError{
			Custom: "failed to delete thread, reasons: user not original poster or user not admin",
			Base:   err,
		}
	}

	err = i.ExactlyOneResultError(threads)
	if err != nil {
		return dummyThread, err
	}

	return threads[0], nil
}

// TODO: Support editing tags
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
		generateParams(title, content, isPublished, isOpen, threadId, userId),
		parseThreadRows,
	)

	if err != nil {
		return dummyThread, &i.DatabaseError{Custom: "failed to edit a thread", Base: err}
	}

	err = i.ExactlyOneResultError(threads)
	if err != nil {
		return dummyThread, err
	}

	return threads[0], err
}
