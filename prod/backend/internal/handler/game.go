package handler

import (
	"net/http"

	"induce-master/internal/service"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService *service.GameService
}

func NewGameHandler(gameService *service.GameService) *GameHandler {
	return &GameHandler{gameService: gameService}
}

func (h *GameHandler) Get(c *gin.Context) {
	gameID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"game_id": gameID,
		"status":  "playing",
	})
}

func (h *GameHandler) SendMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *GameHandler) SubmitWord(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *GameHandler) Guess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}
