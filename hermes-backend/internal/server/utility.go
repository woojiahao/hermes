package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"woojiahao.com/hermes/internal/database"
)

func getPayloadUser(ctx *gin.Context, db *database.Database) (User, error) {
	if u, ok := ctx.Get(IdentityKey); ok {
		username := u.(*User).Username
		user, err := db.GetUser(username)
		if err != nil {
			return User{}, err
		}

		return userToDTO(user), nil
	}

	return User{}, errors.New("invalid payload")
}
