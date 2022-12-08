package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"woojiahao.com/hermes/internal"
	"woojiahao.com/hermes/internal/database"
)

var userRoutes = []route{
	{"GET", "/users/:id", getUser},
	{"GET", "/users", getUsers},
	{"POST", "/users", createUser},
}

func getUser(ctx *gin.Context, db *database.Database) {
	id := ctx.Param("id")
	var user database.User
	var err error
	if internal.IsUUID(id) {
		user, err = db.GetUserById(id)
	} else {
		user, err = db.GetUser(id)
	}

	if err != nil {
		notFound(ctx, fmt.Sprintf("Unable to find user by given user id: %s", id))
		return
	}

	ctx.JSON(http.StatusFound, userToDTO(user))
}

func getUsers(ctx *gin.Context, db *database.Database) {
	users, err := db.GetUsers()
	if err != nil {
		internalSeverError(ctx)
		return
	}

	user_dtos := internal.Map(users, userToDTO)
	ctx.JSON(http.StatusFound, user_dtos)
}

func createUser(ctx *gin.Context, db *database.Database) {
	var req CreateUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestValidation(ctx, err)
		return
	}

	user, err := db.CreateUser(req.Username, req.Email, req.Password, database.Role(req.Role))
	if err != nil {
		internalSeverError(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, userToDTO(user))
}

func userToDTO(user database.User) User {
	return User{user.Id, user.Username, user.Email, user.PasswordHash, string(user.Role)}
}
