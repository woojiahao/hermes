package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"woojiahao.com/hermes/internal/database"
)

var healthRoutes = []route{
	{GET, "/ping", ping},
}

func ping(c *gin.Context, db *database.Database) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
