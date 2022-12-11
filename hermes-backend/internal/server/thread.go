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
		// TODO: Have more specific control over internal server error or not found
		notFound(ctx, fmt.Sprintf("Unable to find thread given id: %s", id))
		return
	}

	ctx.JSON(http.StatusOK, threadToDTO(thread))
}

func createThread(ctx *gin.Context, db *database.Database) {
	var req CreateThread
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestValidation(ctx, err)
		return
	}
  fmt.Println(req)

	thread, err := db.CreateThread(
		req.UserId,
		req.Title,
		req.Content,
		internal.Map(req.Tags, func(tag Tag) database.Tag {
			return database.NewTag(tag.Content, tag.HexCode)
		}),
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

	ctx.JSON(http.StatusOK, tags)
}

func tagToDTO(tag database.Tag) Tag {
	return Tag{tag.Content, tag.HexCode}
}

func threadToDTO(thread database.Thread) Thread {
	return Thread{
		thread.Id,
		thread.IsPublished,
		thread.IsOpen,
		thread.Title,
		thread.Content,
		internal.Map(thread.Tags, tagToDTO),
	}
}
