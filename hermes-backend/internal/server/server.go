package server

import (
	"fmt"
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

type (
	httpAction string

	route struct {
		action httpAction
		route  string
		body   func(*gin.Context, *database.Database)
	}

	// TODO: Reassess if this struct is necessary
	Server struct {
		configuration *ServerConfiguration
		db            *database.Database
		router        *gin.Engine
		routes        []route
	}
)

const (
	GET    httpAction = "GET"
	POST   httpAction = "POST"
	DELETE httpAction = "DELETE"
	PUT    httpAction = "PUT"
)

func Start(c *ServerConfiguration, db *database.Database) {
	router := gin.Default()
	server := &Server{c, db, router, make([]route, 0)}
	server.loadRoutes()
	server.addRoutes()
	server.router.Run(fmt.Sprintf(":%d", c.Port))
}

// Loading all routes into the server instance
func (s *Server) loadRoutes() {
	s.routes = append(s.routes, healthRoutes...)
}

func (s *Server) addRoutes() {
	for _, route := range s.routes {
		switch route.action {
		case GET:
			s.router.GET(route.route, ginBody(route, s.db))
		case POST:
			s.router.POST(route.route, ginBody(route, s.db))
		case DELETE:
			s.router.DELETE(route.route, ginBody(route, s.db))
		case PUT:
			s.router.PUT(route.route, ginBody(route, s.db))
		default:
			log.Fatal("Invalid HTTP action loaded")
		}
	}
}

func ginBody(route route, db *database.Database) func(*gin.Context) {
	return func(ctx *gin.Context) {
		route.body(ctx, db)
	}
}
