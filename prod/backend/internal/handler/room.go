package handler

import (
	"net/http"
	"strconv"

	"induce-master/internal/service"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
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

	c.JSON(http.StatusOK, gin.H{"success": true, "room_id": roomID})
}

// Helper to parse int
func parseInt(s string, def int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return def
}
