package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"woojiahao.com/hermes/internal"
	"woojiahao.com/hermes/internal/database"
)

var authRoutes = []route{
	{"GET", "/auth/check", check, true},
	{"POST", "/register", register, false},
}

func check(ctx *gin.Context, db *database.Database) {
	if user, ok := ctx.Get(IdentityKey); ok {
		u := user.(*User)
		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Hello %s %s", u.Role, u.Username),
		})
		return
	}

	badRequest(ctx, "Failed to check")
}

func register(ctx *gin.Context, db *database.Database) {
	var req CreateUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestValidation(ctx, err)
		return
	}

	user, err := db.CreateUser(req.Username, req.Password, database.Role(req.Role))
	if err != nil {
		if dbe, ok := err.(*internal.DatabaseError); ok {
			log.Println(dbe.Error())
			badRequest(ctx, dbe.Short)
			return
		} else {
			// Hash function failed
			internalSeverError(ctx)
			return
		}
	}

	ctx.JSON(http.StatusCreated, userToDTO(user))
}
