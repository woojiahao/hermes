package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"woojiahao.com/hermes/internal"
	"woojiahao.com/hermes/internal/database"
)

var commentRoutes = []route{
	{"GET", "/threads/:id/comments", getThreadComments, false},
	{"POST", "/threads/:id/comments", createComment, true},
	{"DELETE", "/threads/:id/comments/:commentId", deleteComment, true},
}

func deleteComment(ctx *gin.Context, db *database.Database) {
	commentId := ctx.Param("commentId")
	user, err := getPayloadUser(ctx, db)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	_, err = db.DeleteComment(user.Id, commentId)
	if err != nil {
		internalSeverError(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func getThreadComments(ctx *gin.Context, db *database.Database) {
	id := ctx.Param("id")
	comments, err := db.GetThreadComments(id)
	if err != nil {
		notFound(ctx, "Unable to find thread comments")
		return
	}

	ctx.JSON(http.StatusOK, internal.Map(comments, commentToDTO))
}

func createComment(ctx *gin.Context, db *database.Database) {
	id := ctx.Param("id")
	var req CreateComment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestValidation(ctx, err)
		return
	}

	comment, err := db.CreateComment(req.UserId, id, req.Content)
	if err != nil {
		internalSeverError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, commentToDTO(comment))
}

func commentToDTO(comment database.Comment) Comment {
	return Comment{
		comment.Id,
		comment.Content,
		comment.CreatedAt,
		comment.CreatedBy,
		comment.Creator,
	}
}
