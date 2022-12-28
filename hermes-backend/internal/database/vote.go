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
	err := rows.Scan(&vote.UserId, &vote.ThreadId)
	return vote, err
}

func (d *Database) MakeVote(userId, threadId string, isUpvote bool) error {
	_, err := query(
		d,
		Insert("vote").
			Columns("user_id", "thread_id", "is_upvote").
			Values(P1, P2, P3).
			OnConflict().
			DoNothing(),
		generateParams(userId, threadId, isUpvote),
		doNothing,
	)
	return ThisOrThat(nil, errors.New("failed to vote"), err == nil)
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

func GetVotes(tx *sql.Tx, threadId string) {

}
