package server

import (
	"log"

	"github.com/krissukoco/go-gin-chat/controllers"
)

type WebsocketManager struct {
	ChatClients    []*controllers.ChatClient
	IncomingClient chan *controllers.ChatClient
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

func (m *WebsocketManager) Broadcast(msg any, userId string) {
	for _, client := range m.ChatClients {
		if client.UserId == userId {
			log.Println("broadcast message to client: ", client.Id)
			client.Conn.WriteJSON(msg)
		}
	}
}

func (m *WebsocketManager) Run(stop chan bool) {
	for {
		select {
		case <-stop:
			return
		case client := <-m.IncomingClient:
			if client != nil {
				m.RegisterChatClient(client)
			}
		default:
			for _, client := range m.ChatClients {
				select {
				case <-client.Exited:
					log.Println("client exited")
					m.UnregisterChatClient(client)
				case msg := <-client.Out:
					log.Println("message from client: ", msg)
					if msg.Type == "new_chat" {
						newChat, ok := msg.Data.(*controllers.WsChatData)
						if !ok {
							log.Println("invalid message schema for 'new_chat'")
							continue
						}
						if newChat.IsGroup && newChat.Group != nil {
							// Send to each members of group
							for _, member := range newChat.Group.MemberIds {
								m.Broadcast(msg, member)
							}
						} else {
							// Send to user
							m.Broadcast(msg, newChat.Receiver.Id)
						}

					}
				default:
					continue
				}
			}
		}
	}
}
