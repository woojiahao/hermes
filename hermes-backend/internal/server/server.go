package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"woojiahao.com/hermes/internal"
	"woojiahao.com/hermes/internal/database"
)

type ServerConfiguration struct {
	Port   int
	JWTKey string
}

func LoadConfiguration() *ServerConfiguration {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal("Invalid SERVER_PORT in .env")
	}

	return &ServerConfiguration{port, os.Getenv("JWT_KEY")}
}

type (
	httpAction string

	route struct {
		action                httpAction
		route                 string
		body                  func(*gin.Context, *database.Database)
		authorizationRequired bool
	}

	// TODO: Reassess if this struct is necessary
	Server struct {
		configuration  *ServerConfiguration
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

func Start(c *ServerConfiguration, db *database.Database) {
	router := gin.Default()
	server := &Server{c, db, router, make([]route, 0), nil}
	server.setupCORS()
	server.setupAuth()
	server.loadRoutes()
	server.addRoutes()
	server.router.Run(fmt.Sprintf(":%d", c.Port))
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
			fmt.Println(data)
			// TODO: Check roles
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
		})
}

func (s *Server) addRoutes() {
	s.router.POST("/login", s.authMiddleware.LoginHandler)
	s.router.POST("/logout", s.authMiddleware.LogoutHandler)
	s.router.GET("/refresh", s.authMiddleware.RefreshHandler)

	internal.ForEach(
		internal.Filter(s.routes, func(r route) bool { return !r.authorizationRequired }),
		s.addRoute,
	)

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

// TODO: Log the errors for production
func ginError(ctx *gin.Context, errorCode int, message any) {
	ctx.JSON(errorCode, errorBody{errorCode, message})
}

func internalSeverError(ctx *gin.Context) {
	ginError(ctx, http.StatusInternalServerError, "Internal server error")
}

func notFound(ctx *gin.Context, message string) {
	ginError(ctx, http.StatusNotFound, message)
}

func badRequestValidation(ctx *gin.Context, bindingError error) {
	var ve validator.ValidationErrors
	if errors.As(bindingError, &ve) {
		out := internal.Map(ve, func(field validator.FieldError) errorField {
			message := ""
			switch field.Tag() {
			case "required":
				message = "This field is required"
			case "min":
				message = "This field has a minimum necessary length/size"
			case "email":
				message = "This field must be an email"
			default:
				message = field.Tag()
			}
			return errorField{field.Field(), message}
		})

		ginError(ctx, http.StatusBadRequest, out)
	}
}

func badRequest(ctx *gin.Context, message string) {
	ginError(ctx, http.StatusBadRequest, message)
}

func getPayloadUser(ctx *gin.Context, db *database.Database) (User, error) {
	if u, ok := ctx.Get(IdentityKey); ok {
		username := u.(*User).Username
		user, err := db.GetUser(username)
		if err != nil {
			return User{}, err
		}

		return userToDTO(user), nil
	}

	return User{}, errors.New("invalid payload")
}
