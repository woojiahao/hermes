package server

import (
	"fmt"
	"log"
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
	log.Println("hi")
	fmt.Println("retrieving current user")
	user, err := getPayloadUser(ctx, db)
	fmt.Println(user)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
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
