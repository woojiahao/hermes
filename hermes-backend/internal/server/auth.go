package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"woojiahao.com/hermes/internal/database"
)

var authRoutes = []route{
	{"GET", "/auth/check", check, true},
	{"POST", "/register", register, false},
}

func check(ctx *gin.Context, db *database.Database) {
	user, err := getPayloadUser(ctx, db)
	if err != nil {
		badRequest(ctx, "Failed to check")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Hello %s %s", user.Role, user.Username),
	})
}

func register(ctx *gin.Context, db *database.Database) {
	var req CreateUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestValidation(ctx, err)
		return
	}

	user, err := db.CreateUser(req.Username, req.Password, req.Role)
	if err != nil {
		internalSeverError(ctx)
	}

	ctx.JSON(http.StatusCreated, userToDTO(user))
}
