package handler

import (
	"encoding/json"
	"log"

	"induce-master/internal/service"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// 注册客户端
	Register chan *Client
	// 注销客户端
	Unregister chan *Client
	// 客户端映射
	Clients map[string]*Client
	// 服务
	userService  *service.UserService
	roomService *service.RoomService
	gameService *service.GameService
}

type Client struct {
	Conn    *websocket.Conn
	UserID  string
	Send    chan []byte
	RoomID  string
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func NewHub(userService *service.UserService, roomService *service.RoomService, gameService *service.GameService) *Hub {
	return &Hub{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[string]*Client),
		userService:  userService,
		roomService: roomService,
		gameService: gameService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = client
			log.Printf("客户端注册: %s", client.UserID)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
				log.Printf("客户端注销: %s", client.UserID)
			}
		}
	}
}

func (h *Hub) SendToUser(userID string, msg Message) {
	if client, ok := h.Clients[userID]; ok {
		data, _ := json.Marshal(msg)
		select {
		case client.Send <- data:
		default:
			delete(h.Clients, userID)
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID string, msg Message) {
	data, _ := json.Marshal(msg)
	for _, client := range h.Clients {
		if client.RoomID == roomID {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}
