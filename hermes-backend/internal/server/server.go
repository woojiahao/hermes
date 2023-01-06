package server

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"woojiahao.com/hermes/internal"
	"woojiahao.com/hermes/internal/database"
)

type (
	httpAction string

	route struct {
		action                httpAction
		route                 string
		body                  func(*gin.Context, *database.Database)
		authorizationRequired bool
	}

	Server struct {
		configuration  *Configuration
		db             *database.Database
		router         *gin.Engine
		routes         []route
		authMiddleware *jwt.GinJWTMiddleware
	}
)

const (
	GET         httpAction = "GET"
	POST        httpAction = "POST"
	DELETE      httpAction = "DELETE"
	PUT         httpAction = "PUT"
	IdentityKey            = "id"
)

func Start(c *Configuration, db *database.Database) {
	router := gin.Default()
	server := &Server{c, db, router, make([]route, 0), nil}
	server.setupCORS()
	server.setupAuth()
	server.loadRoutes()
	server.addRoutes()
	server.router.Run()
}

func (s *Server) setupCORS() {
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
	}))
}

func (s *Server) setupAuth() {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "hermes",
		Key:         []byte(s.configuration.JWTKey),
		Timeout:     time.Hour,
		MaxRefresh:  168 * time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{IdentityKey: v.Username, "role": v.Role}
			}

			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx *gin.Context) interface{} {
			claims := jwt.ExtractClaims(ctx)
			return &User{Username: claims[IdentityKey].(string), Role: claims["role"].(string)}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login Login
			if err := c.ShouldBind(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := login.Username
			password := login.Password

			user, err := s.db.GetUser(username)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &User{
				user.Id,
				user.Username,
				string(user.Role),
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok {
				if _, err := s.db.GetUser(v.Username); err != nil {
					return false
				}
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			ginError(c, code, message)
		},
	})

	if err != nil {
		log.Fatal("Failed to setup JWT auth middleware")
	}

	err = authMiddleware.MiddlewareInit()
	if err != nil {
		log.Fatal("Failed to initialize JWT auth middleware")
	}

	s.authMiddleware = authMiddleware
}

// Loading all routes into the server instance
func (s *Server) loadRoutes() {
	s.routes = internal.Flatten(
		[][]route{
			healthRoutes,
			authRoutes,
			userRoutes,
			threadRoutes,
			commentRoutes,
			voteRoutes,
		})
}

func (s *Server) addRoutes() {
	s.router.POST("/login", s.authMiddleware.LoginHandler)
	s.router.POST("/logout", s.authMiddleware.LogoutHandler)
	s.router.GET("/refresh", s.authMiddleware.RefreshHandler)

	// Must load routes that do not use the JWT authentication middleware
	internal.ForEach(
		internal.Filter(s.routes, func(r route) bool { return !r.authorizationRequired }),
		s.addRoute,
	)

	// Load routes that use JWT authentication
	internal.ForEach(
		internal.Filter(s.routes, func(r route) bool { return r.authorizationRequired }),
		func(r route) {
			s.router.Use(s.authMiddleware.MiddlewareFunc())
			s.addRoute(r)
		},
	)
}

func (s *Server) addRoute(r route) {
	switch r.action {
	case GET:
		s.router.GET(r.route, ginBody(r, s.db))
	case POST:
		s.router.POST(r.route, ginBody(r, s.db))
	case DELETE:
		s.router.DELETE(r.route, ginBody(r, s.db))
	case PUT:
		s.router.PUT(r.route, ginBody(r, s.db))
	default:
		log.Fatal("Invalid HTTP action loaded")
	}
}

func ginBody(route route, db *database.Database) func(*gin.Context) {
	return func(ctx *gin.Context) {
		route.body(ctx, db)
	}
}
