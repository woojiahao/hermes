package database

import "time"

type Vote struct {
	UserId   string
	ThreadId string
	IsUpvote bool
}

type Comment struct {
	Content   string
	CreatedAt time.Time
	CreatedBy string
	ThreadId  string
}
