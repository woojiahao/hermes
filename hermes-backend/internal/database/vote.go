package database

import (
	"database/sql"
	"errors"
	. "woojiahao.com/hermes/internal"
	. "woojiahao.com/hermes/internal/database/q"
)

type (
	VoteType string
	Vote     struct {
		UserId   string
		ThreadId string
		IsUpvote bool
	}
)

const (
	UPVOTE   VoteType = "UPVOTE"
	DOWNVOTE VoteType = "DOWNVOTE"
	NOVOTE   VoteType = "NOVOTE"
)

func parseVoteRows(rows *sql.Rows) (Vote, error) {
	var vote Vote
	err := rows.Scan(&vote.UserId, &vote.IsUpvote)
	return vote, err
}

func (d *Database) MakeVote(userId, threadId string, isUpvote bool) error {
	_, err := transaction(d, func(tx *sql.Tx) (error, error) {
		_, err := transactionSingleQuery(
			tx,
			From("vote").
				Select("user_id", "is_upvote").
				Where(And(Eq("user_id", P1), Eq("thread_id", P2))),
			generateParams(userId, threadId),
			parseVoteRows,
		)

		if err == nil {
			// Remove the previous vote first
			_, err := transactionQuery(
				tx,
				Delete("vote").Where(And(Eq("user_id", P1), Eq("thread_id", P2))),
				generateParams(userId, threadId),
				doNothing,
			)
			if err != nil {
				return nil, err
			}
		}

		_, err = transactionQuery(
			tx,
			Insert("vote").
				Columns("user_id", "thread_id", "is_upvote").
				Values(P1, P2, P3).
				OnConflict().
				DoNothing(),
			generateParams(userId, threadId, isUpvote),
			doNothing,
		)
		return nil, ThisOrThat(nil, errors.New("failed to vote"), err == nil)
	})

	return err
}

func (d *Database) ClearVote(userId, threadId string) error {
	_, err := query(
		d,
		Delete("vote").Where(And(Eq("user_id", P1), Eq("thread_id", P2))),
		generateParams(userId, threadId),
		doNothing,
	)

	return ThisOrThat(nil, errors.New("failed to clear vote"), err == nil)
}
