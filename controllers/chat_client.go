package controllers

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/krissukoco/go-gin-chat/security"
)

type ChatClient struct {
	UserId        string
	Id            string
	Authenticated bool
	Conn          *websocket.Conn
	Incoming      chan *WsBaseMessage
	Exited        chan bool
}

func NewChatClient(userId string, connection *websocket.Conn) *ChatClient {
	return &ChatClient{
		UserId:   userId,
		Conn:     connection,
		Id:       newWsId(),
		Incoming: make(chan *WsBaseMessage),
		Exited:   make(chan bool),
	}
}

func (cl *ChatClient) Websocket(jwtSecret string) {
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
		err = cl.processMessage(&m)
		if err != nil {
			log.Println("Error processMessage: ", err)
			break
		}
		cl.Incoming <- &WsBaseMessage{
			Type: "message",
			Data: string(msg),
		}
	}
}

func (cl *ChatClient) processMessage(m *WsBaseMessage) error {
	log.Println("Processing message: ", m)
	switch m.Type {
	// User authentication
	case "auth":
		if cl.Authenticated {
			log.Println("Client is already authenticated")
			return nil
		}
		var authMsg WsAuthMsg
		b, err := json.Marshal(m.Data)
		if err != nil {
			log.Println("Error Marshal: ", err)
			break
		}
		err = json.Unmarshal(b, &authMsg)
		if err != nil {
			log.Println("Error Unmarshal: ", err)
			break
		}
		userId, err := security.GetUserIdFromJwt(authMsg.Token)
		if err != nil {
			log.Println("Error GetUserIdFromJwt: ", err)
			break
		}
		cl.UserId = userId
		cl.Authenticated = true
		log.Println("Client is authenticated")
	}
	return nil
}
