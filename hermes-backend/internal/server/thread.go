package server

import (
	"github.com/gin-gonic/gin"
	"woojiahao.com/hermes/internal/database"
)

var threadRoutes = []route{
	{"GET", "/threads", getThreads},
	{"GET", "/threads/:id", getThreadById},
	{"POST", "/threads", createThread},
}

func getThreads(ctx *gin.Context, db *database.Database) {

}

func getThreadById(ctx *gin.Context, db *database.Database) {

}

func createThread(ctx *gin.Context, db *database.Database) {

}
