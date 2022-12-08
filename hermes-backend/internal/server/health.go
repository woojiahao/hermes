package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"woojiahao.com/hermes/internal/database"
)

func ping(c *gin.Context, db *database.Database) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
