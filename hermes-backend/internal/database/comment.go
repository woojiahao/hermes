package database

import (
	"database/sql"
	"time"

	i "woojiahao.com/hermes/internal"
)

type Comment struct {
	Id        string
	Content   string
	CreatedAt *time.Time
	CreatedBy string
	ThreadId  string
	DeletedAt *time.Time
	DeletedBy sql.NullString
	Creator   string
}

var dummyComment Comment

func parseCommentRowsWithCreator(rows *sql.Rows) (Comment, error) {
	var comment Comment
	err := rows.Scan(
		&comment.Id,
		&comment.Content,
		&comment.CreatedAt,
		&comment.CreatedBy,
		&comment.ThreadId,
		&comment.DeletedAt,
		&comment.DeletedBy,
		&comment.Creator,
	)
	return comment, err
}

func parseCommentRows(rows *sql.Rows) (Comment, error) {
	var comment Comment
	err := rows.Scan(
		&comment.Id,
		&comment.Content,
		&comment.CreatedAt,
		&comment.CreatedBy,
		&comment.ThreadId,
		&comment.DeletedAt,
		&comment.DeletedBy,
	)
	return comment, err
}

func (d *Database) CreateComment(userId, threadId, content string) (Comment, error) {
	return transaction(d, func(tx *sql.Tx) (Comment, error) {
		comments, err := transactionQuery(
			tx,
			`
			INSERT INTO comment ("content", created_by, thread_id)
			VALUES ($1, $2, $3)
			RETURNING *;
		`,
			generateParams(content, userId, threadId),
			parseCommentRows,
		)
		if err != nil {
			return dummyComment, &i.DatabaseError{Custom: "failed to insert new comment", Base: err}
		}

		usernames, err := transactionQuery(
			tx,
			`
				SELECT username FROM "user" WHERE id = $1
			`,
			generateParams(userId),
			func(rows *sql.Rows) (string, error) {
				var username string
				err := rows.Scan(&username)
				return username, err
			},
		)

		comment := comments[0]
		comment.Creator = usernames[0]

		return comment, nil
	})
}

func (d *Database) DeleteComment(userId, commentId string) (Comment, error) {
	comments, err := query(
		d,
		`
			UPDATE comment
			SET deleted_by = $1, deleted_at = NOW()
			WHERE comment.id = $2
				AND (
					comment.created_by = $1 OR
					EXISTS (SELECT * FROM "user" WHERE "user".id = $1 AND "user".role = 'ADMIN')
				)
			RETURNING *;
		`,
		generateParams(userId, commentId),
		parseCommentRows,
	)
	if err != nil {
		return dummyComment, &i.DatabaseError{
			Custom: "failed to delete comment, reasons: user not original poster or user not admin",
			Base:   err,
		}
	}

	err = i.ExactlyOneResultError(comments)
	if err != nil {
		return dummyComment, err
	}

	return comments[0], nil
}

func (d *Database) GetThreadComments(threadId string) ([]Comment, error) {
	comments, err := query(
		d,
		`
			SELECT comment.*, "user".username 
			FROM comment 
				INNER JOIN "user" ON "user".id = comment.created_by
			WHERE thread_id = $1 AND deleted_at IS NULL
			ORDER BY created_at DESC;
		`,
		generateParams(threadId),
		parseCommentRowsWithCreator,
	)
	if err != nil {
		return nil, &i.DatabaseError{Custom: "failed to retrieve comments for given thread", Base: err}
	}

	return comments, nil
}
