// DTOs to be used for the server
package server

import (
	"time"
	"woojiahao.com/hermes/internal/database"
)

type (
	errorBody struct {
		HttpCode int `json:"http_code"`
		Message  any `json:"message"`
	}

	errorField struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

type (
	Login struct {
		Username string `form:"username" json:"username" binding:"required,min=3"`
		Password string `json:"password" binding:"required,min=3"`
	}
)

type (
	CreateUser struct {
		Username string        `form:"username" json:"username" binding:"required,min=3"`
		Password string        `json:"password" binding:"required,min=3"`
		Role     database.Role `json:"role" binding:"required"`
	}

	User struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}
)

type (
	Tag struct {
		Content string `json:"content" binding:"required"`
		HexCode string `json:"hex_code" binding:"required"`
	}

	CreateThread struct {
		Title   string `json:"title" binding:"required,min=5"`
		Content string `json:"content" binding:"required,min=30"`
		Tags    []Tag  `json:"tags"`
	}

	EditThread struct {
		Title       string `json:"title" binding:"required,min=5"`
		Content     string `json:"content" binding:"required,min=30"`
		IsPublished *bool  `json:"is_published" binding:"required"`
		IsOpen      *bool  `json:"is_open" binding:"required"`
		Tags        []Tag  `json:"tags"`
	}

	PinThread struct {
		IsPinned *bool `json:"is_pinned" binding:"required"`
	}

	Thread struct {
		Id          string     `json:"id"`
		IsPublished bool       `json:"is_published"`
		IsOpen      bool       `json:"is_open"`
		IsPinned    bool       `json:"is_pinned"`
		Title       string     `json:"title"`
		Content     string     `json:"content"`
		Tags        []Tag      `json:"tags"`
		CreatedAt   *time.Time `json:"created_at"`
		CreatedBy   string     `json:"created_by"`
		Creator     string     `json:"creator"`
		Upvoters    []string   `json:"upvoters"`
		Downvoters  []string   `json:"downvoters"`
		Upvotes     int        `json:"upvotes"`
		Downvotes   int        `json:"downvotes"`
	}

	CreateComment struct {
		UserId  string `json:"user_id" binding:"required,uuid"`
		Content string `json:"content" binding:"required"`
	}

	Comment struct {
		Id        string     `json:"id"`
		Content   string     `json:"content"`
		CreatedAt *time.Time `json:"created_at"`
		CreatedBy string     `json:"created_by"`
		Creator   string     `json:"creator"`
	}

	MakeVote struct {
		IsUpvote bool `json:"is_upvote"`
	}
)
