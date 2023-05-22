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
	In            chan *WsBaseMessage
	Out           chan *WsBaseMessage
	ChatData      chan *models.Chat
	Exited        chan bool
}

func NewChatClient(userId string, connection *websocket.Conn) *ChatClient {
	return &ChatClient{
		UserId:   userId,
		Conn:     connection,
		Id:       newWsId(),
		In:       make(chan *WsBaseMessage),
		Out:      make(chan *WsBaseMessage),
		Exited:   make(chan bool),
		ChatData: make(chan *models.Chat),
	}
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
