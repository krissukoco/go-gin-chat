package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/krissukoco/go-gin-chat/models"
	"github.com/krissukoco/go-gin-chat/schema"
	"github.com/krissukoco/go-gin-chat/security"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrMessageTypeUnknown = errors.New("message type is unknown")
	ErrChatNotFound       = errors.New("chat not found")
)

type Chat struct {
	Mongo     *mongo.Database
	UserCtl   *User // bridge to user controller to get user data
	JwtSecret string
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
	// Process client websocket
	chat.clientWebsocket(client)
	// Wait until ws process finished
	log.Println("Client exited")
	client.Exited <- true
}

func (chat *Chat) clientWebsocket(cl *ChatClient) {
	for {
		mt, msg, err := cl.Conn.ReadMessage()
		if err != nil {
			log.Println("Error ReadMessage: ", err)
			break
		}
		if string(msg) == "ping" {
			cl.Conn.WriteMessage(mt, []byte("pong"))
			continue
		}
		var m WsBaseMessage
		err = json.Unmarshal(msg, &m)
		if err != nil {
			log.Println("Error Unmarshal: ", err)
			break
		}
		// Process message
		err = chat.processMessage(cl, &m)
		if err != nil {
			cl.sendJson(&WsBaseMessage{
				Type: "error",
				Data: map[string]string{
					"message": err.Error(),
				}})
			break
		}
		cl.Incoming <- &WsBaseMessage{
			Type: "message",
			Data: string(msg),
		}
	}
}

func (chat *Chat) processMessage(cl *ChatClient, m *WsBaseMessage) error {
	log.Println("Processing message: ", m)
	switch m.Type {
	// User authentication
	case "auth":
		if cl.Authenticated {
			return cl.sendJson(&WsBaseMessage{
				Type: "error",
				Data: map[string]string{"message": "already authenticated"},
			})
		}
		var authMsg WsAuthMsg
		err := cl.convertData(m.Data, &authMsg)
		if err != nil {
			log.Println("Error convertData: ", err)
			break
		}
		userId, err := security.GetUserIdFromJwt(authMsg.Token)
		if err != nil {
			log.Println("Error GetUserIdFromJwt: ", err)
			break
		}
		cl.UserId = userId
		cl.Authenticated = true
		cl.sendJson(&WsBaseMessage{
			Type: "success",
			Data: map[string]string{"message": "authenticated"},
		})
		log.Println("Client is authenticated")
	// Chat
	case "send_chat":
		if !cl.Authenticated {
			return ErrAbortConnection
		}
		var chatData WsChatMsg
		err := cl.convertData(m.Data, &chatData)
		if err != nil {
			return err
		}
		// process chat by type
		switch chatData.Type {
		case "text":
			if chatData.Text == "" {
				return ErrInvalidSchema
			}
			// Save chat to database
			now := time.Now().UnixMilli()
			chatModel := models.Chat{
				SenderId:  cl.UserId,
				ChatId:    chatData.ChatId,
				Type:      "text",
				Text:      chatData.Text,
				ReadBy:    make([]string, 0),
				CreatedAt: now,
				UpdatedAt: now,
			}
			log.Println("Chat data: ", chatModel)
			err = chat.processClientChat(cl, &chatModel)
			if err != nil {
				return err
			}
		default:
			return ErrInvalidSchema
		}
	default:
		return ErrMessageTypeUnknown
	}
	return nil
}

func (chat *Chat) processClientChat(client *ChatClient, chatData *models.Chat) error {
	log.Println("Listening to chat data...")
	if chatData.ChatId == "" {
		return errors.New("chat id cannot be empty")
	}
	// User cannot send to themselves
	if chatData.SenderId == chatData.ChatId {
		return errors.New("cannot send to yourself")
	}
	// Find group chat
	var group models.Group
	err := group.FindById(chat.Mongo, chatData.ChatId)
	if err == nil {
		// Group exists
		chatData.IsGroup = true
		err = chatData.Save(chat.Mongo)
		if err != nil {
			log.Println("ERROR saving group chat data to mongo: ", err)
			return err
		}
	}
	// Find by user id
	_, err = chat.UserCtl.GetUserById(chatData.ChatId)
	if err != nil {
		// No group nor user found
		log.Println("ERROR finding user: ", err)
		return ErrChatNotFound
	}
	log.Println("New chat data: ", chatData)
	err = chatData.Save(chat.Mongo)
	if err != nil {
		log.Println("ERROR saving chat data to mongo: ", err)
		return err
	}
	return nil
}
