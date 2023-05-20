package controllers

import (
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"
	"github.com/krissukoco/go-gin-chat/models"
)

var (
	ErrAbortConnection = errors.New("connection is aborted either by server or client")
	ErrInvalidSchema   = errors.New("invalid message schema")
)

type ChatClient struct {
	UserId        string
	Id            string
	Authenticated bool
	Conn          *websocket.Conn
	Incoming      chan *WsBaseMessage
	ChatData      chan *models.Chat
	Exited        chan bool
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

func NewChatClient(userId string, connection *websocket.Conn) *ChatClient {
	return &ChatClient{
		UserId:   userId,
		Conn:     connection,
		Id:       newWsId(),
		Incoming: make(chan *WsBaseMessage),
		Exited:   make(chan bool),
		ChatData: make(chan *models.Chat),
	}
}

func (cl *ChatClient) convertData(data any, target any) error {
	// Marshal and Unmarshall
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, target)
	return err
}

func (cl *ChatClient) sendJson(m any) error {
	if cl.Conn == nil {
		return nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return cl.Conn.WriteMessage(1, b)
}
