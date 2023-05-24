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
	"github.com/krissukoco/go-gin-chat/utils"
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

type WsBaseMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
type WsAuthMsg struct {
	Token string `json:"token"`
}
type WsChatMsg struct {
	Type   string `json:"type"`
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}
type WsChatData struct {
	*models.Chat
	Receiver *models.User  `json:"receiver"`
	Group    *models.Group `json:"group"`
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

func (chat *Chat) GetAllChats(userId string) ([]*models.ChatRoom, error) {
	// Get all chats
	chats, err := models.GetUserChatRooms(chat.Mongo, userId)
	for _, room := range chats {
		// Get user data if room is user chat
		if room.ChatId[:2] == "u_" {
			room.User, err = chat.UserCtl.GetUserById(room.ChatId)
			if err != nil {
				return nil, err
			}
		}
	}
	return chats, err
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
	// Wait until ws process finished then send exited
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
		err := utils.ConvertStruct(m.Data, &authMsg)
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
		// Find user
		user, err := chat.UserCtl.GetUserById(userId)
		if err != nil {
			log.Println("Error FindUserById: ", err)
			break
		}
		cl.Authenticated = true
		cl.sendJson(&WsBaseMessage{
			Type: "success",
			Data: map[string]interface{}{
				"message": "authenticated",
				"user":    user,
			},
		})
		log.Println("Client is authenticated")
	// Chat
	case "send_chat":
		if !cl.Authenticated {
			return ErrAbortConnection
		}
		var chatData WsChatMsg
		err := utils.ConvertStruct(m.Data, &chatData)
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
			chatExtended, err := chat.processClientChat(cl, &chatModel)
			if err != nil {
				return err
			}
			// Send back to sender
			chatMsg := &WsBaseMessage{
				Type: "chat_sent",
				Data: chatExtended,
			}
			if err = cl.sendJson(chatMsg); err != nil {
				return err
			}
			// Send to receiver(s)
			chatMsg.Type = "new_chat"
			cl.Out <- chatMsg

		default:
			return ErrInvalidSchema
		}
	case "get_chats":
		if !cl.Authenticated {
			return ErrAbortConnection
		}
		chats, err := chat.GetAllChats(cl.UserId)
		if err != nil {
			return err
		}
		cl.sendJson(&WsBaseMessage{
			Type: "chats",
			Data: chats,
		})

	default:
		return ErrMessageTypeUnknown
	}
	return nil
}

func (chat *Chat) processClientChat(client *ChatClient, chatData *models.Chat) (*WsChatData, error) {
	log.Println("Listening to chat data...")
	if chatData.ChatId == "" {
		return nil, errors.New("chat id cannot be empty")
	}
	// User cannot send to themselves
	if chatData.SenderId == chatData.ChatId {
		return nil, errors.New("cannot send to yourself")
	}
	var data WsChatData
	// Find group chat
	var group models.Group
	err := group.FindById(chat.Mongo, chatData.ChatId)
	if err == nil {
		// Group exists
		chatData.IsGroup = true
		data.Group = &group
		err = chatData.Save(chat.Mongo)
		if err != nil {
			log.Println("ERROR saving group chat data to mongo: ", err)
			return nil, err
		}
	}

	// Find by user id
	user, err := chat.UserCtl.GetUserById(chatData.ChatId)
	if err != nil && !chatData.IsGroup {
		// No group nor user found
		log.Println("ERROR finding user: ", err)
		return nil, ErrChatNotFound
	}
	data.Receiver = user
	log.Println("New chat data: ", chatData)
	err = chatData.Save(chat.Mongo)
	if err != nil {
		log.Println("ERROR saving chat data to mongo: ", err)
		return nil, err
	}
	data.Chat = chatData
	return &data, nil
}
