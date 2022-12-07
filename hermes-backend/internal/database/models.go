package database

import "time"

type Tag struct {
	Id      string
	Content string
	HexCode string
}

type Role struct {
	Id          string
	Title       string
	Permissions string
}

type User struct {
	Id           string
	Username     string
	Email        string
	PasswordHash string
}

type Thread struct {
	Id          string
	IsPublished bool
	IsOpen      bool
	Content     string
	CreatedAt   time.Time
	CreatedBy   string
	DeletedAt   time.Time
	DeletedBy   string
}

type Vote struct {
	Id       string
	UserId   string
	ThreadId string
	IsUpvote bool
}

type Comment struct {
	Id        string
	Content   string
	CreatedAt time.Time
	CreatedBy string
	ThreadId  string
}
