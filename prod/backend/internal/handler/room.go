package handler

import (
	"net/http"
	"strconv"

	"induce-master/internal/service"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService *service.RoomService
	hub        *Hub
}

func NewRoomHandler(roomService *service.RoomService, hub *Hub) *RoomHandler {
	return &RoomHandler{roomService: roomService, hub: hub}
}

func (h *RoomHandler) List(c *gin.Context) {
	rooms, err := h.roomService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 转换为 API 响应格式
	type RoomResponse struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		OwnerID    string `json:"owner_id"`
		Status     int    `json:"status"`
		PlayerCount int   `json:"player_count"`
	}
	var response []RoomResponse
	for _, r := range rooms {
		players, _ := h.roomService.GetPlayers(r.ID)
		response = append(response, RoomResponse{
			ID:           r.ID,
			Name:         r.Name,
			OwnerID:      r.OwnerID,
			Status:       r.Status,
			PlayerCount:  len(players),
		})
	}
	c.JSON(http.StatusOK, gin.H{"rooms": response})
}

func (h *RoomHandler) Create(c *gin.Context) {
	var req struct {
		Name      string `json:"name"`
		OwnerID   string `json:"owner_id"`
		Password  string `json:"password"`
		MaxPlayers int   `json:"max_players"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Name == "" {
		req.Name = "房间"
	}
	if req.MaxPlayers == 0 {
		req.MaxPlayers = 4
	}

	room := &service.DBRoom{
		ID:       service.GenerateUUID(),
		Name:     req.Name,
		OwnerID:  req.OwnerID,
		Password: req.Password,
		Status:   0, // waiting
	}

	err := h.roomService.Create(room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"room_id": room.ID,
		"name":    room.Name,
	})
}

func (h *RoomHandler) Get(c *gin.Context) {
	roomID := c.Param("id")
	room, err := h.roomService.GetByID(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	players, _ := h.roomService.GetPlayers(roomID)
	c.JSON(http.StatusOK, gin.H{
		"room_id":      room.ID,
		"name":         room.Name,
		"owner_id":     room.OwnerID,
		"status":       room.Status,
		"players":      players,
		"player_count": len(players),
	})
}

func (h *RoomHandler) Join(c *gin.Context) {
	roomID := c.Param("id")
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.roomService.AddPlayer(roomID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "room_id": roomID})
}

func (h *RoomHandler) Leave(c *gin.Context) {
	roomID := c.Param("id")
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.roomService.RemovePlayer(roomID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *RoomHandler) Ready(c *gin.Context) {
	roomID := c.Param("id")
	var req struct {
		UserID string `json:"user_id"`
		Ready  bool   `json:"ready"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.roomService.SetReady(roomID, req.UserID, req.Ready)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "ready": req.Ready})
}

func (h *RoomHandler) Start(c *gin.Context) {
	roomID := c.Param("id")

	// 检查房间状态
	room, err := h.roomService.GetByID(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	if room.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room is not waiting"})
		return
	}

	// 检查玩家是否都准备好了
	players, err := h.roomService.GetPlayers(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(players) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need at least 2 players"})
		return
	}

	// 检查是否都准备了
	allReady := true
	for _, p := range players {
		if !p.Ready {
			allReady = false
			break
		}
	}

	if !allReady {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not all players ready"})
		return
	}

	// 更新房间状态为游戏中
	err = h.roomService.UpdateRoomStatus(roomID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建游戏引擎
	playerIDs := make([]string, len(players))
	for i, p := range players {
		playerIDs[i] = p.UserID
	}
	game := h.hub.CreateGame(roomID, playerIDs)

	// 通知所有玩家游戏开始
	for i, p := range players {
		word := game.GetWord(p.UserID)
		var opponentID string
		if i == 0 && len(players) > 1 {
			opponentID = players[1].UserID
		} else if i > 0 {
			opponentID = players[0].UserID
		}
		
		h.hub.SendGameStart(roomID, roomID, word, opponentID)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "room_id": roomID})
}

// Helper to parse int
func parseInt(s string, def int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return def
}

// SendMessage 发送游戏消息 (供 AI Agent 使用)
func (h *RoomHandler) SendMessage(c *gin.Context) {
	roomID := c.Param("id")
	
	var req struct {
		UserID  string `json:"user_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	
	// 获取用户名
	username := req.UserID
	if user, err := h.roomService.GetUserByID(req.UserID); err == nil {
		username = user.Username
	}
	
	// 处理游戏消息
	game := h.hub.GetGame(roomID)
	if game == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game not found"})
		return
	}
	
	// 检查是否是当前玩家的回合
	currentPlayer := game.GetCurrentPlayer()
	if currentPlayer != req.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not your turn"})
		return
	}
	
	// 处理消息
	msg := h.hub.ProcessGameMessage(roomID, req.UserID, username, req.Content)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"is_keyword": msg != nil && msg.IsKeyword,
		"score": game.Scores[req.UserID],
	})
}

// GetGameStatus 获取游戏状态 (供 AI Agent 轮询)
func (h *RoomHandler) GetGameStatus(c *gin.Context) {
	roomID := c.Param("id")
	
	game := h.hub.GetGame(roomID)
	if game == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "not_started",
			"round": 0,
		})
		return
	}
	
	// 获取当前玩家
	currentPlayer := game.GetCurrentPlayer()
	
	c.JSON(http.StatusOK, gin.H{
		"status": "playing",
		"round": game.Round,
		"max_rounds": game.MaxRounds,
		"current_player": currentPlayer,
		"scores": game.Scores,
	})
}
