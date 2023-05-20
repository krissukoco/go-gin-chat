package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/krissukoco/go-gin-chat/schema"
	"go.mongodb.org/mongo-driver/mongo"
)

type Chat struct {
	Mongo     *mongo.Database
	UserCtl   *User // bridge to user controller to get user data
	JwtSecret string
}

type WsBaseMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type WsAuthMsg struct {
	Token string `json:"token"`
}

func (chat *Chat) GetAll(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.JSON(500, &schema.ErrorResponse{
			Code:    schema.ErrInternalServer,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(200, gin.H{"username": username})
}

func newWsId() string {
	return "ws_" + uuid.NewString()
}

func (chat *Chat) ChatWebsocketHandler(c *gin.Context, newClient chan *ChatClient) {
	wsIntf, exists := c.Get("ws")
	if !exists {
		c.JSON(500, &schema.ErrorResponse{
			Code:    schema.ErrInternalServer,
			Message: "Internal server error",
		})
		return
	}
	ws, ok := wsIntf.(*websocket.Conn)
	if !ok {
		c.JSON(500, &schema.ErrorResponse{
			Code:    schema.ErrInternalServer,
			Message: "Internal server error",
		})
		return
	}

	// Create and process new client
	client := NewChatClient("anon", ws)
	newClient <- client
	client.Websocket(chat.JwtSecret)
	// Wait until ws process finished
	log.Println("Client exited")
	client.Exited <- true
}
