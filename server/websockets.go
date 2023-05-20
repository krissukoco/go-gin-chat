package server

import (
	"log"

	"github.com/krissukoco/go-gin-chat/controllers"
)

type WebsocketManager struct {
	ChatClients []*controllers.ChatClient
	// IncomingClient chan *controllers.ChatClient
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		ChatClients: make([]*controllers.ChatClient, 0),
	}
}

func (m *WebsocketManager) RegisterChatClient(client *controllers.ChatClient) {
	log.Println("Registering chat client: ", client.Id)
	for _, c := range m.ChatClients {
		if c.Id == client.Id {
			return
		}
	}

	m.ChatClients = append(m.ChatClients, client)
	log.Printf("Registered chat client: %v\n", client.Id)
}

func (m *WebsocketManager) UnregisterChatClient(client *controllers.ChatClient) {
	log.Println("Removing chat client: ", client.Id)
	for i, c := range m.ChatClients {
		if c.Id == client.Id {
			m.ChatClients = append(m.ChatClients[:i], m.ChatClients[i+1:]...)
			return
		}
	}
}

func (m *WebsocketManager) Run(stop chan bool, newClient chan *controllers.ChatClient) {
	for {
		select {
		case <-stop:
			return
		case client := <-newClient:
			if client != nil {
				m.RegisterChatClient(client)
			}
		default:
			for _, client := range m.ChatClients {
				select {
				case <-client.Exited:
					log.Println("client exited")
					m.UnregisterChatClient(client)
				case msg := <-client.Incoming:
					log.Println("incoming message", msg)
				}
			}
		}
	}
}
