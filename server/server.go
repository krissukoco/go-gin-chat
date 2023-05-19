package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krissukoco/go-gin-chat/controllers"
	"github.com/krissukoco/go-gin-chat/database"
	"github.com/krissukoco/go-gin-chat/middlewares"
	"github.com/krissukoco/go-gin-chat/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Server struct {
	Engine *gin.Engine
	Pg     *gorm.DB
	Mongo  *mongo.Database
	Port   int
}

func NewDefaultServer() (*Server, error) {
	// Database connections
	pg, err := database.NewPostgres()
	if err != nil {
		return nil, err
	}
	mongoDb, err := database.NewMongo()
	if err != nil {
		return nil, err
	}

	defaultPort := 8000
	portEnv, ok := os.LookupEnv("PORT")
	if ok {
		port, err := strconv.Atoi(portEnv)
		if err == nil {
			defaultPort = port
		}
	}

	jwtSecret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	// Routing and initiating controllers and middlewares
	router := NewDefaultRouter()
	authCtl := controllers.Auth{
		Pg:        pg,
		JwtSecret: jwtSecret,
	}
	userCtl := controllers.User{
		Pg: pg,
	}
	chatCtl := controllers.Chat{
		Mongo:   mongoDb,
		UserCtl: &userCtl,
	}
	authMiddleware := middlewares.AuthMiddleware{
		JwtSecret: jwtSecret,
	}
	router.POST("/auth/login", authCtl.Login)
	router.POST("/auth/register", authCtl.Register)
	router.GET("/auth/account", authMiddleware.AuthorizationHeader, authCtl.GetAccount)
	router.GET("/users", userCtl.GetAll)
	router.GET("/users/:id", userCtl.GetById)
	router.GET("/chats", authMiddleware.AuthorizationHeader, chatCtl.GetAll)
	// router.POST("/chats", authMiddleware.AuthorizationHeader, chatCtl.Create)

	srv := &Server{
		Engine: router,
		Pg:     pg,
		Mongo:  mongoDb,
		Port:   defaultPort,
	}
	// Database auto migrate
	srv.databaseAutoMigrate()
	return srv, nil
}

func (srv *Server) databaseAutoMigrate() {
	srv.Pg.AutoMigrate(&models.User{})
}

func (srv *Server) Start() {
	srv.Engine.Run(fmt.Sprintf(":%d", srv.Port))
}
