package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"woojiahao.com/hermes/internal/database"
)

var healthRoutes = []route{
	{GET, "/ping", ping, false},
}

func ping(c *gin.Context, _ *database.Database) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
