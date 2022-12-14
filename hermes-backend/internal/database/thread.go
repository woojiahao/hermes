package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	i "woojiahao.com/hermes/internal"
	. "woojiahao.com/hermes/internal/database/q"
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
	Upvoters    []string
	Downvoters  []string
	Upvotes     int
	Downvotes   int
}

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
		log.Printf("Given tags is not unique")
		return Thread{}, DataError
	}

	return transaction(d, func(tx *sql.Tx) (Thread, error) {
		thread, err := transactionSingleQuery(
			tx,
			Insert("thread").
				Columns("title", `"content"`, "created_by").
				Values(P1, P2, P3).
				Returning(ALL),
			generateParams(title, content, userId),
			parseThreadRows,
		)
		if err != nil {
			return Thread{}, err
		}

		user, err := d.GetUserById(userId)
		if err != nil {
			return Thread{}, err
		}

		thread, err = attachTags(tx, thread, tags)
		if err != nil {
			return Thread{}, err
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
		// Retrieve the thread
		thread, err := transactionSingleQuery(
			tx,
			From("thread").
				Select("thread.*", `"user".username`).
				InnerJoin(`"user"`, "thread.created_by", `"user".id`).
				Where(And(Eq("thread.id", P1), IsNull("deleted_at"))),
			generateParams(threadId),
			parseThreadRowsWithCreator,
		)

		if err != nil {
			return Thread{}, err
		}

		// Retrieve all of its tags
		threadTags, err := transactionQuery(
			tx,
			From("tag").
				Select("tag.*").
				InnerJoin("thread_tag", "tag.id", "thread_tag.tag_id").
				InnerJoin("thread", "thread.id", "thread_tag.thread_id").
				Where(Eq("thread.id", P1)),
			generateParams(thread.Id),
			parseTagRows,
		)
		if err != nil {
			return Thread{}, err
		}

		thread.Tags = threadTags

		// Retrieve all the votes
		transform := func(v Vote) string {
			return v.UserId
		}
		votes, err := transactionQuery(
			tx,
			From("vote").
				Select("user_id", "is_upvote").
				Where(Eq("thread_id", P1)),
			generateParams(thread.Id),
			parseVoteRows,
		)
		thread.Upvoters = i.FilterMap(
			votes,
			func(v Vote) bool {
				return v.IsUpvote
			},
			transform,
		)
		thread.Downvoters = i.FilterMap(
			votes,
			func(v Vote) bool {
				return !v.IsUpvote
			},
			transform,
		)
		thread.Upvotes = len(thread.Upvoters)
		thread.Downvotes = len(thread.Downvoters)

		return thread, nil
	})
}

func (d *Database) GetThreads() ([]Thread, error) {
	return getThreadsWithFilter(d, nil)
}

func (d *Database) DeleteThread(userId, threadId string) (Thread, error) {
	isAdmin := From(`"user"`).Select(ALL).Where(And(Eq(`"user".id`, P1), Eq(`"user".role`, "'ADMIN'"))).Generate()
	isValid := Or(Eq("thread.created_by", P1), Exists(isAdmin))
	where := And(Eq("thread.id", P2), isValid)

	threads, err := query(
		d,
		Update("thread").
			Set("deleted_by", P1).
			Set("deleted_at", NOW).
			Where(where).
			Returning(ALL),
		generateParams(userId, threadId),
		parseThreadRows,
	)
	if err != nil {
		return Thread{}, err
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
		log.Printf("Given tags is not unique")
		return Thread{}, DataError
	}

	return transaction(d, func(tx *sql.Tx) (Thread, error) {
		thread, err := transactionSingleQuery(
			tx,
			Update("thread").
				Set("title", P1).
				Set(`"content"`, P2).
				Set("is_published", P3).
				Set("is_open", P4).
				Set("updated_at", NOW).
				Where(And(Eq("thread.id", P5), Eq("thread.created_by", P6))).
				Returning(ALL),
			generateParams(title, content, isPublished, isOpen, threadId, userId),
			parseThreadRows,
		)
		if err != nil {
			return Thread{}, err
		}

		// Delete all existing tags
		_, err = transactionSingleQuery(
			tx,
			Delete("thread_tag").Where(Eq("thread_id", P1)),
			generateParams(threadId),
			doNothing,
		)
		if err != nil {
			return Thread{}, err
		}

		thread, err = attachTags(tx, thread, tags)
		if err != nil {
			return Thread{}, err
		}

		return thread, err
	})
}

func (d *Database) GetTags() ([]Tag, error) {
	tags, err := query(
		d,
		From("tag").Select(ALL),
		generateParams(),
		parseTagRows,
	)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (d *Database) PinThread(threadId string, isPinned bool) (Thread, error) {
	thread, err := singleQuery(
		d,
		Update("thread").Set("is_pinned", P1).Where(Eq("id", P2)).Returning(ALL),
		generateParams(isPinned, threadId),
		parseThreadRows,
	)
	if err != nil {
		return Thread{}, err
	}

	return thread, nil
}

func getThreadsWithFilter(d *Database, userId *string) ([]Thread, error) {
	// We don't need to retrieve all of the specific voters so we just collect the total of each
	where := And(IsNull("deleted_at"), "is_published")
	sub := From("vote").
		Select("thread_id", "COUNT(is_upvote) FILTER (WHERE is_upvote) upvotes", "COUNT(is_upvote) FILTER (WHERE not is_upvote) downvotes").
		Group("thread_id")
	q := From("thread").
		Select(
			"thread.*",
			`"user".username`,
			Coalaesce("sub.upvotes", 0),
			Coalaesce("sub.downvotes", 0),
		).
		InnerJoin(`"user"`, "thread.created_by", `"user".id`).
		LeftJoin(Sub(sub, "sub"), "thread.id", "sub.thread_id").
		Order("is_pinned", DESC).
		Order("created_at", DESC)

	params := generateParams()
	if userId != nil {
		where = And(where, Eq(`"user".id`, P1))
		params = generateParams(userId)
	}
	q.Where(where)

	return transaction(d, func(tx *sql.Tx) ([]Thread, error) {
		// Get the threads
		threads, err := transactionQuery(tx, q, params, func(rows *sql.Rows) (Thread, error) {
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
				&thread.Upvotes,
				&thread.Downvotes,
			)

			return thread, err
		})
		if err != nil {
			return nil, err
		}

		// Attach the tags to the threads
		threadTags := make(map[string][]Tag)
		_, err = transactionQuery(
			tx,
			From("tag").
				Select("thread.id", "tag.id", `tag."content"`, "tag.hex_code").
				InnerJoin("thread_tag", "tag.id", "thread_tag.tag_id").
				InnerJoin("thread", "thread.id", "thread_tag.thread_id").
				Where(IsNull("thread.deleted_at")),
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
			return nil, err
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
	var parameters []string
	var tagValues []any
	// Dynamically generating the parameters and its values
	for k, j := 1, 0; k < 2*len(tags); k, j = k+2, j+1 {
		parameters = append(parameters, fmt.Sprintf("$%d", k), fmt.Sprintf("$%d", k+1))
		tagValues = append(tagValues, tags[j].Content, tags[j].HexCode)
	}

	insertTagsQuery := Insert("tag").
		Columns(`"content"`, "hex_code").
		OnConflict(`"content"`).
		DoUpdate(map[string]any{`"content"`: `EXCLUDED."content"`}).
		Returning(ALL)

	for k := 0; k < 2*len(tags); k += 2 {
		insertTagsQuery.Values(parameters[k], parameters[k+1])
	}

	ts, err := transactionQuery(tx, insertTagsQuery, tagValues, parseTagRows)
	if err != nil {
		return Thread{}, err
	}

	var joiningValues []any
	// Dynamically generate the joining table values
	for _, t := range ts {
		joiningValues = append(joiningValues, thread.Id, t.Id)
	}

	insertJoiningQuery := Insert("thread_tag").Returning(ALL)
	for k := 0; k < 2*len(tags); k += 2 {
		insertJoiningQuery.Values(parameters[k], parameters[k+1])
	}
	_, err = transactionQuery(tx, insertJoiningQuery, joiningValues, doNothing)
	if err != nil {
		return Thread{}, err
	}

	thread.Tags = ts

	return thread, nil
}
