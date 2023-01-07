package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"woojiahao.com/hermes/internal/database"
)

var voteRoutes = []route{
	{"PUT", "/threads/:id/votes", makeVote, true},
	{"DELETE", "/threads/:id/votes", deleteVote, true},
}

func makeVote(ctx *gin.Context, db *database.Database) {
	id := ctx.Param("id")
	user, err := getPayloadUser(ctx, db)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	var req MakeVote
	if err = ctx.ShouldBindJSON(&req); err != nil {
		badRequestValidation(ctx, err)
		return
	}

	err = db.MakeVote(user.Id, id, req.IsUpvote)
	if err != nil {
		internalSeverError(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func deleteVote(ctx *gin.Context, db *database.Database) {
	id := ctx.Param("id")
	user, err := getPayloadUser(ctx, db)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	err = db.ClearVote(user.Id, id)
	if err != nil {
		internalSeverError(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}
