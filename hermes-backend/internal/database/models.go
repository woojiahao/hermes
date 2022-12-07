package database

import "time"

type Tag struct {
	Content string
	HexCode string
}

type Role struct {
	Title       string
	Permissions string
}

type User struct {
	Username     string
	Email        string
	PasswordHash string
}

type Thread struct {
	IsPublished bool
	IsOpen      bool
	Content     string
	CreatedAt   time.Time
	CreatedBy   string
	DeletedAt   time.Time
	DeletedBy   string
}

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
