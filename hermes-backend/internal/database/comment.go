package database

import (
	"database/sql"
	"time"
	. "woojiahao.com/hermes/internal/database/q"
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
		comment, err := transactionSingleQuery(
			tx,
			Insert("comment").
				Values(P1, P2, P3).
				Columns(`"comment"`, "created_by", "thread_id").
				Returning(ALL),
			generateParams(content, userId, threadId),
			parseCommentRows,
		)
		if err != nil {
			return Comment{}, InternalError
		}

		username, err := transactionSingleQuery(
			tx,
			From(`"user"`).
				Select("username").
				Where(Eq("id", "$1")),
			generateParams(userId),
			func(rows *sql.Rows) (string, error) {
				var username string
				err := rows.Scan(&username)
				return username, err
			},
		)

		comment.Creator = username

		return comment, nil
	})
}

func (d *Database) DeleteComment(userId, commentId string) (Comment, error) {
	isAdmin := From(`"user"`).Select(ALL).Where(And(Eq(`"user".id`, P1), Eq(`"user".role`, "ADMIN"))).Generate()
	isValid := Or(Eq("comment.created_by", P1), Exists(isAdmin))
	where := And(Eq("comment.id", P2), isValid)

	comment, err := singleQuery(
		d,
		Update("comment").
			Set("deleted_by", P1).
			Set("deleted_at", NOW).
			Where(where).
			Returning(ALL),
		generateParams(userId, commentId),
		parseCommentRows,
	)
	if err != nil {
		return Comment{}, InternalError
	}

	return comment, nil
}

func (d *Database) GetThreadComments(threadId string) ([]Comment, error) {
	comments, err := query(
		d,
		From("comment").
			Select("comment.*", `"user".username`).
			InnerJoin(`"user"`, "comment.created_by", `"user".id`).
			Where(And(Eq("thread_id", P1), IsNull("deleted_at"))).
			Order("created_at", DESC),
		generateParams(threadId),
		parseCommentRowsWithCreator,
	)
	if err != nil {
		return nil, InternalError
	}

	return comments, nil
}
