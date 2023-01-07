package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"woojiahao.com/hermes/internal"
	"woojiahao.com/hermes/internal/database"
)

var threadRoutes = []route{
	{"GET", "/threads", getThreads, false},
	{"GET", "/threads/:id", getThreadById, false},
	{"POST", "/threads", createThread, true},
	{"GET", "/threads/tags", getTags, false},
	{"GET", "/threads/current", getCurrentUserThreads, true},
	{"DELETE", "/threads/:id", deleteThread, true},
	{"PUT", "/threads/:id", editThread, true},
	{"PUT", "/threads/:id/pin", pinThread, true},
}

func pinThread(ctx *gin.Context, db *database.Database) {
	id := ctx.Param("id")
	user, err := getPayloadUser(ctx, db)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	if user.Role != string(database.ADMIN) {
		badRequest(ctx, "Invalid user action")
		return
	}

	var req PinThread
	if err = ctx.ShouldBindJSON(&req); err != nil {
		badRequestValidation(ctx, err)
		return
	}

	thread, err := db.PinThread(id, *req.IsPinned)
	if err != nil {
		badRequest(ctx, "Cannot pin thread")
		return
	}

	ctx.JSON(http.StatusOK, threadToDTO(thread))
}

func editThread(ctx *gin.Context, db *database.Database) {
	id := ctx.Param("id")
	user, err := getPayloadUser(ctx, db)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	var req EditThread
	if err = ctx.ShouldBindJSON(&req); err != nil {
		badRequestValidation(ctx, err)
		return
	}

	thread, err := db.EditThread(
		user.Id,
		id,
		req.Title,
		req.Content,
		*req.IsPublished,
		*req.IsOpen,
		internal.Map(req.Tags, tagToDatabaseObj),
	)
	if err != nil {
		if err == database.DataError {
			badRequest(ctx, "Tags must be unique")
		} else {
			internalSeverError(ctx)
		}
		return
	}

	ctx.JSON(http.StatusOK, thread)
}

func deleteThread(ctx *gin.Context, db *database.Database) {
	id := ctx.Param("id")
	user, err := getPayloadUser(ctx, db)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	_, err = db.DeleteThread(user.Id, id)
	if err != nil {
		internalSeverError(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func getCurrentUserThreads(ctx *gin.Context, db *database.Database) {
	u, err := getPayloadUser(ctx, db)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	threads, err := db.GetUserThreads(u.Id)
	if err != nil {
		notFound(ctx, "Unable to find user threads")
	}

	ctx.JSON(http.StatusOK, internal.Map(threads, threadToDTO))
}

func getThreads(ctx *gin.Context, db *database.Database) {
	threads, err := db.GetThreads()
	if err != nil {
		internalSeverError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, internal.Map(threads, threadToDTO))
}

func getThreadById(ctx *gin.Context, db *database.Database) {
	id := ctx.Param("id")
	thread, err := db.GetThreadById(id)
	if err != nil {
		if err == database.NotFoundError {
			notFound(ctx, fmt.Sprintf("Unable to find thread given id: %s", id))
		} else {
			internalSeverError(ctx)
		}
		return
	}

	ctx.JSON(http.StatusOK, threadToDTO(thread))
}

func createThread(ctx *gin.Context, db *database.Database) {
	user, err := getPayloadUser(ctx, db)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	var req CreateThread
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestValidation(ctx, err)
		return
	}

	thread, err := db.CreateThread(
		user.Id,
		req.Title,
		req.Content,
		internal.Map(req.Tags, tagToDatabaseObj),
	)
	if err != nil {
		internalSeverError(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, threadToDTO(thread))
}

func getTags(ctx *gin.Context, db *database.Database) {
	tags, err := db.GetTags()
	if err != nil {
		notFound(ctx, "Cannot find tags")
		return
	}

	ctx.JSON(http.StatusOK, internal.Map(tags, tagToDTO))
}

func tagToDatabaseObj(tag Tag) database.Tag {
	return database.NewTag(tag.Content, tag.HexCode)
}

func tagToDTO(tag database.Tag) Tag {
	return Tag{tag.Content, tag.HexCode}
}

func threadToDTO(thread database.Thread) Thread {
	return Thread{
		thread.Id,
		thread.IsPublished,
		thread.IsOpen,
		thread.IsPinned,
		thread.Title,
		thread.Content,
		internal.Map(thread.Tags, tagToDTO),
		thread.CreatedAt,
		thread.CreatedBy,
		thread.Creator,
		thread.Upvoters,
		thread.Downvoters,
		thread.Upvotes,
		thread.Downvotes,
	}
}
