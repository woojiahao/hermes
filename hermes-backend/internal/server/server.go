package server

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"woojiahao.com/hermes/internal/database"
)

type ServerConfiguration struct {
	Port int
}

func LoadConfiguration() *ServerConfiguration {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal("Invalid SERVER_PORT in .env")
	}

	return &ServerConfiguration{port}
}

type Server struct {
	configuration *ServerConfiguration
	db            *database.Database
	router        *gin.Engine
}

func Start(c *ServerConfiguration, db *database.Database) {
	router := gin.Default()
	server := &Server{c, db, router}
	server.loadRoutes()
	server.router.Run()
}

func (s *Server) loadRoutes() {

}
