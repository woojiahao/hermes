package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"woojiahao.com/hermes/internal"
	"woojiahao.com/hermes/internal/database"
)

var userRoutes = []route{
	{"GET", "/users/:id", getUser, true},
	{"GET", "/users", getUsers, true},
	{"GET", "/users/current", getCurrentUser, true},
}

func getCurrentUser(ctx *gin.Context, db *database.Database) {
	if u, ok := ctx.Get(IdentityKey); ok {
		username := u.(*User).Username
		user, err := db.GetUser(username)
		if err != nil {
			notFound(ctx, "Unable to find user")
			return
		}
		ctx.JSON(http.StatusOK, userToDTO(user))
		return
	}

	badRequest(ctx, "Failed to retrieve current user")
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

	ctx.JSON(http.StatusOK, userToDTO(user))
}

func getUsers(ctx *gin.Context, db *database.Database) {
	users, err := db.GetUsers()
	if err != nil {
		internalSeverError(ctx)
		return
	}

	userDTOs := internal.Map(users, userToDTO)
	ctx.JSON(http.StatusOK, userDTOs)
}

func userToDTO(user database.User) User {
	return User{user.Id, user.Username, string(user.Role)}
}
