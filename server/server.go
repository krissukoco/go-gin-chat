package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krissukoco/go-gin-chat/controllers"
	"github.com/krissukoco/go-gin-chat/database"
	"github.com/krissukoco/go-gin-chat/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Server struct {
	Engine    *gin.Engine
	Pg        *gorm.DB
	Mongo     *mongo.Database
	WsManager *WebsocketManager
	Port      int
	NewClient chan *controllers.ChatClient
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

	wsManager := NewWebsocketManager()

	// Router
	srv := &Server{
		Pg:        pg,
		Mongo:     mongoDb,
		Port:      defaultPort,
		WsManager: wsManager,
		NewClient: make(chan *controllers.ChatClient),
	}
	srv.WsManager.IncomingClient = srv.NewClient
	err = srv.setupRouter()
	if err != nil {
		return nil, err
	}

	srv.databaseAutoMigrate()

	// Run WS manager
	stop := make(chan bool)
	go srv.WsManager.Run(stop)
	return srv, nil
}

func (srv *Server) databaseAutoMigrate() {
	srv.Pg.AutoMigrate(&models.User{})
}

func (srv *Server) Start() {
	srv.Engine.Run(fmt.Sprintf(":%d", srv.Port))
}
