package database

import (
	"database/sql"
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

func tagsAreUnique(tags []Tag) bool {
	var tagContents []string
	for _, tag := range tags {
		tagContents = append(tagContents, tag.Content)
	}

	return !i.HasDuplicates(tagContents)
}

// Thread is the database model for Thread objects. Contains references to the related Tags that have been attached with
// it along with the Creator's username (instead of just the user id)
type Thread struct {
	Id          string
	IsPublished bool
	IsOpen      bool
	IsPinned    bool
	Title       string
	Content     string
	CreatedAt   *time.Time
	CreatedBy   string
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	DeletedBy   sql.NullString
	Tags        []Tag
	Creator     string
}

var dummyThread Thread

// parseThreadRowsWithCreator will parse the results of a sql.Rows from a SELECT thread.*, "user".username query
func parseThreadRowsWithCreator(rows *sql.Rows) (Thread, error) {
	var thread Thread
	err := rows.Scan(
		&thread.Id,
		&thread.IsPublished,
		&thread.IsOpen,
		&thread.IsPinned,
		&thread.Title,
		&thread.Content,
		&thread.CreatedAt,
		&thread.CreatedBy,
		&thread.UpdatedAt,
		&thread.DeletedAt,
		&thread.DeletedBy,
		&thread.Creator,
	)
	return thread, err
}

// parseThreadRows will parse the results of a sql.Rows from a SELECT * query
func parseThreadRows(rows *sql.Rows) (Thread, error) {
	var thread Thread
	err := rows.Scan(
		&thread.Id,
		&thread.IsPublished,
		&thread.IsOpen,
		&thread.IsPinned,
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

// CreateThread inserts a new thread by given user (identified by userId) along with its attached tags.
// The tags provided must have unique contents, otherwise a internal.ServerError will be produced.
// Performs three different SQL queries: 1) INSERT INTO thread, 2) SELECT * "user", 3) INSERT INTO tag
func (d *Database) CreateThread(userId, title, content string, tags []Tag) (Thread, error) {
	if !tagsAreUnique(tags) {
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

		user, err := d.GetUserById(userId)
		if err != nil {
			return dummyThread, &i.DatabaseError{Custom: "failed to retrieve user", Short: "Invalid user_id", Base: err}
		}

		thread, err := attachTags(tx, threads[0], tags)
		if err != nil {
			return dummyThread, err
		}

		thread.Creator = user.Username

		return thread, nil
	})
}

func (d *Database) GetUserThreads(userId string) ([]Thread, error) {
	return getThreadsWithFilter(d, &userId)
}

func (d *Database) GetThreadById(threadId string) (Thread, error) {
	return transaction(d, func(tx *sql.Tx) (Thread, error) {
		threads, err := transactionQuery(
			tx,
			`
				SELECT thread.*, "user".username
				FROM thread
					INNER JOIN "user" ON thread.created_by = "user".id
				WHERE thread.id = $1 AND deleted_at IS NULL
			`,
			generateParams(threadId),
			parseThreadRowsWithCreator,
		)

		if err != nil {
			return dummyThread, &i.DatabaseError{Custom: "failed to retrieve thread by id", Base: err}
		}

		thread := threads[0]

		threadTags, err := transactionQuery(
			tx,
			`
				SELECT tag.*
				FROM tag
					INNER JOIN thread_tag tt on tag.id = tt.tag_id
					INNER JOIN thread t on t.id = tt.thread_id
				WHERE t.id = $1;
			`,
			generateParams(thread.Id),
			parseTagRows,
		)
		if err != nil {
			return dummyThread, &i.DatabaseError{Custom: "failed to retrieve thread tags", Base: err}
		}

		thread.Tags = threadTags

		return thread, nil
	})
}

func (d *Database) GetThreads() ([]Thread, error) {
	return getThreadsWithFilter(d, nil)
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
			RETURNING *;
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

	return threads[0], nil
}

func (d *Database) EditThread(
	userId,
	threadId,
	title,
	content string,
	isPublished,
	isOpen bool,
	tags []Tag,
) (Thread, error) {
	if !tagsAreUnique(tags) {
		return dummyThread, &i.ServerError{Custom: "tags cannot have same name", Base: nil}
	}

	return transaction(d, func(tx *sql.Tx) (Thread, error) {
		threads, err := transactionQuery(
			tx,
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

		// Delete all existing tags
		_, err = transactionQuery(
			tx,
			`DELETE FROM thread_tag WHERE thread_id = $1`,
			generateParams(threadId),
			func(r *sql.Rows) (string, error) {
				return "", nil
			},
		)
		if err != nil {
			return dummyThread, &i.DatabaseError{Custom: "failed to delete all associated tags with given thread", Base: err}
		}

		thread, err := attachTags(tx, threads[0], tags)
		if err != nil {
			return dummyThread, err
		}

		return thread, err
	})
}

func (d *Database) GetTags() ([]Tag, error) {
	tags, err := query(
		d,
		`SELECT * FROM tag;`,
		generateParams(),
		parseTagRows,
	)
	if err != nil {
		return nil, &i.DatabaseError{Custom: "failed to retrieve all tags", Short: "Cannot retrieve tags", Base: err}
	}

	return tags, nil
}

func (d *Database) PinThread(threadId string, isPinned bool) (Thread, error) {
	threads, err := query(
		d,
		`
			UPDATE thread
			SET is_pinned = $1
			WHERE thread.id = $2
			RETURNING *;
		`,
		generateParams(isPinned, threadId),
		parseThreadRows,
	)
	if err != nil {
		return dummyThread, &i.DatabaseError{Custom: "failed to pin thread", Short: "Cannot pin thread", Base: err}
	}

	return threads[0], nil
}

func getThreadsWithFilter(d *Database, userId *string) ([]Thread, error) {
	query := `
		SELECT thread.*, "user".username
		FROM thread
			INNER JOIN "user" ON thread.created_by = "user".id
		WHERE deleted_at IS NULL AND is_published
	`
	params := generateParams()
	if userId != nil {
		query += ` AND "user".id = $1`
		params = generateParams(userId)
	}
	query += "\nORDER BY is_pinned DESC, created_at DESC;"

	return transaction(d, func(tx *sql.Tx) ([]Thread, error) {
		threads, err := transactionQuery(
			tx,
			query,
			params,
			parseThreadRowsWithCreator,
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
					INNER JOIN thread t on t.id = tt.thread_id
				WHERE t.deleted_at IS NULL;
			`,
			generateParams(),
			func(r *sql.Rows) (string, error) {
				var threadId string
				var tag Tag
				err := r.Scan(&threadId, &tag.Id, &tag.Content, &tag.HexCode)

				threadTags[threadId] = append(threadTags[threadId], tag)

				return "", err
			},
		)

		if err != nil {
			return nil, &i.DatabaseError{Custom: "failed to retrieve all tags related to all threads", Base: err}
		}

		var finalThreads []Thread
		for _, thread := range threads {
			copyThread := thread
			copyThread.Tags = threadTags[thread.Id]
			finalThreads = append(finalThreads, copyThread)
		}

		return finalThreads, nil
	})
}

func attachTags(tx *sql.Tx, thread Thread, tags []Tag) (Thread, error) {
	var dbTags []Tag

	// Create or retrieve all the tags from the database
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

		dbTags = append(dbTags, ts[0])

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

	thread.Tags = dbTags

	return thread, nil
}
