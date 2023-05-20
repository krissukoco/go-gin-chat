package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/krissukoco/go-gin-chat/controllers"
	"github.com/krissukoco/go-gin-chat/middlewares"
	"github.com/krissukoco/go-gin-chat/security"
)

func newDefaultRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	return r
}

func (srv *Server) setupRouter() error {
	jwtSecret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		return fmt.Errorf("JWT_SECRET is not set")
	}
	security.SetJwtSecret(jwtSecret)

	// Middlewares
	authMiddleware := middlewares.AuthMiddleware{
		JwtSecret: jwtSecret,
	}
	// Routers
	router := newDefaultRouter()
	authCtl := controllers.Auth{
		Pg:        srv.Pg,
		JwtSecret: jwtSecret,
	}
	userCtl := controllers.User{
		Pg: srv.Pg,
	}
	chatCtl := controllers.Chat{
		Mongo:     srv.Mongo,
		UserCtl:   &userCtl,
		JwtSecret: jwtSecret,
	}
	router.POST("/auth/login", authCtl.Login)
	router.POST("/auth/register", authCtl.Register)
	router.GET("/auth/account", authMiddleware.AuthorizationHeader, authCtl.GetAccount)
	router.GET("/users", userCtl.GetAll)
	router.GET("/users/:id", userCtl.GetById)
	router.GET("/chats", authMiddleware.AuthorizationHeader, chatCtl.GetAll)
	// Websockets
	ws := router.Group("/ws", middlewares.WebsocketMiddleware)
	ws.GET("/chats", func(c *gin.Context) {
		chatCtl.ChatWebsocketHandler(c, srv.NewClient)
	})

	srv.Engine = router
	return nil
}
