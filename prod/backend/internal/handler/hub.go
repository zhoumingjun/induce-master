package handler

import (
	"encoding/json"
	"log"
	"time"

	"induce-master/internal/service"

	"github.com/gorilla/websocket"
)

const (
	// 消息类型
	MsgTypeConnect     = "connect"
	MsgTypeRoomUpdate  = "room_update"
	MsgTypeGameStart   = "game_start"
	MsgTypeGameMessage = "message"
	MsgTypeGameEnd     = "game_end"
	MsgTypePing        = "ping"
	MsgTypePong        = "pong"
	MsgTypeError       = "error"
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
	// 游戏引擎
	Games map[string]*service.GameEngine
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

// RoomUpdatePayload 房间更新消息
type RoomUpdatePayload struct {
	RoomID     string                   `json:"room_id"`
	Players    []PlayerInfo             `json:"players"`
	Status     int                      `json:"status"`
}

// PlayerInfo 玩家信息
type PlayerInfo struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Ready    bool   `json:"ready"`
	Score    int    `json:"score"`
}

// GameStartPayload 游戏开始消息
type GameStartPayload struct {
	GameID      string            `json:"game_id"`
	Round       int               `json:"round"`
	TimeLimit   int               `json:"time_limit"`
	YourWord    string            `json:"your_word"`
	Opponent    PlayerInfo        `json:"opponent"`
}

// GameMessagePayload 游戏消息
type GameMessagePayload struct {
	Round    int       `json:"round"`
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
}

// GameEndPayload 游戏结束消息
type GameEndPayload struct {
	GameID    string     `json:"game_id"`
	WinnerID  string     `json:"winner_id"`
	Scores    map[string]int `json:"scores"`
	Words     map[string]string `json:"words"`
}

func NewHub(userService *service.UserService, roomService *service.RoomService, gameService *service.GameService) *Hub {
	return &Hub{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[string]*Client),
		userService:  userService,
		roomService: roomService,
		gameService: gameService,
		Games:       make(map[string]*service.GameEngine),
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
	
	// 从数据库获取房间内所有玩家
	players, _ := h.roomService.GetPlayers(roomID)
	for _, p := range players {
		if client, ok := h.Clients[p.UserID]; ok {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

// SendRoomUpdate 推送房间状态更新
func (h *Hub) SendRoomUpdate(roomID string) {
	players, err := h.roomService.GetPlayers(roomID)
	if err != nil {
		log.Printf("获取房间玩家失败: %v", err)
		return
	}

	room, err := h.roomService.GetByID(roomID)
	if err != nil {
		log.Printf("获取房间失败: %v", err)
		return
	}

	var playerInfos []PlayerInfo
	for _, p := range players {
		playerInfos = append(playerInfos, PlayerInfo{
			ID:       p.ID,
			UserID:   p.UserID,
			Username: p.UserID, // TODO: 获取用户名
			Ready:    p.Ready,
			Score:    p.Score,
		})
	}

	payload := RoomUpdatePayload{
		RoomID:  roomID,
		Players: playerInfos,
		Status:  room.Status,
	}

	msg := Message{
		Type:    MsgTypeRoomUpdate,
		Payload: payload,
	}

	h.BroadcastToRoom(roomID, msg)
}

// SendGameStart 推送游戏开始
func (h *Hub) SendGameStart(roomID string, gameID string, yourWord string, opponentID string) {
	opponent := PlayerInfo{
		UserID: opponentID,
	}

	payload := GameStartPayload{
		GameID:    gameID,
		Round:     1,
		TimeLimit: 180,
		YourWord:  yourWord,
		Opponent:  opponent,
	}

	msg := Message{
		Type:    MsgTypeGameStart,
		Payload: payload,
	}

	// 直接发送给房间内所有玩家（不通过 RoomID 过滤）
	players, _ := h.roomService.GetPlayers(roomID)
	for _, p := range players {
		h.SendToUser(p.UserID, msg)
	}
}

// CreateGame 创建游戏
func (h *Hub) CreateGame(roomID string, players []string) *service.GameEngine {
	engine := service.NewGameEngine(roomID, players, 10, 60)
	engine.Start()
	h.Games[roomID] = engine
	return engine
}

// GetGame 获取游戏引擎
func (h *Hub) GetGame(roomID string) *service.GameEngine {
	return h.Games[roomID]
}

// ProcessGameMessage 处理游戏消息
func (h *Hub) ProcessGameMessage(roomID string, userID string, username string, content string) *service.GameMessage {
	engine := h.Games[roomID]
	if engine == nil {
		return nil
	}

	// 处理消息
	msg := engine.ProcessMessage(userID, username, content)

	// 广播消息
	h.SendGameMessage(roomID, engine.Round, userID, username, content)

	// 切换到下一个玩家
	engine.NextTurn()

	// 如果说出关键词，检查是否游戏结束
	if msg.IsKeyword {
		// 广播关键词触发
		payload := map[string]interface{}{
			"user_id":  userID,
			"keyword":   engine.GetOpponentWord(userID),
			"score":    engine.Scores[userID],
		}
		h.BroadcastToRoom(roomID, Message{
			Type:    "keyword_triggered",
			Payload: payload,
		})

		// 游戏结束
		winner := userID
		h.SendGameEnd(roomID, roomID, winner, engine.Scores, engine.Words)
		delete(h.Games, roomID)

		// 更新房间状态
		h.roomService.UpdateRoomStatus(roomID, 0)
	}

	return msg
}

// SendGameMessage 推送游戏消息
func (h *Hub) SendGameMessage(roomID string, round int, userID string, username string, content string) {
	payload := GameMessagePayload{
		Round:    round,
		UserID:   userID,
		Username: username,
		Content:  content,
		Time:     time.Now(),
	}

	msg := Message{
		Type:    MsgTypeGameMessage,
		Payload: payload,
	}

	h.BroadcastToRoom(roomID, msg)
}

// SendGameEnd 推送游戏结束
func (h *Hub) SendGameEnd(roomID string, gameID string, winnerID string, scores map[string]int, words map[string]string) {
	payload := GameEndPayload{
		GameID:   gameID,
		WinnerID: winnerID,
		Scores:   scores,
		Words:    words,
	}

	msg := Message{
		Type:    MsgTypeGameEnd,
		Payload: payload,
	}

	h.BroadcastToRoom(roomID, msg)
}
