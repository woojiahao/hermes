// DTOs to be used for the server
package server

import "woojiahao.com/hermes/internal/database"

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
	CreateUser struct {
		Username string        `form:"username" json:"username" binding:"required,min=3"`
		Email    string        `form:"email" json:"email" binding:"required,email"`
		Password string        `json:"password" binding:"required,min=3"`
		Role     database.Role `json:"role" binding:"required"`
	}

	User struct {
		Id           string `json:"id"`
		Username     string `json:"username"`
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash"`
		Role         string `json:"role"`
	}
)
