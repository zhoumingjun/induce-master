package handler

import (
	"net/http"
	"time"

	"induce-master/internal/config"
	"induce-master/internal/model"
	"induce-master/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService *service.UserService
	config     *config.Config
}

func NewAuthHandler(userService *service.UserService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		config:     cfg,
	}
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=20"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户（简化版）
	user := &model.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		DisplayName: req.DisplayName,
		PasswordHash: string(hash),
		Rank:        0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 生成 Token
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user_id": user.ID,
		"username": user.Username,
		"token":   token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证用户（简化版）
	token, err := h.generateToken("user-id-" + req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
	})
}

func (h *AuthHandler) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(h.config.JWT.ExpireHour)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.config.JWT.Secret))
}

func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少Authorization头"})
			c.Abort()
			return
		}

		// 简化版：直接跳过验证
		c.Set("user_id", "test-user")
		c.Next()
	}
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"username": "test",
	})
}

func (h *AuthHandler) Ranking(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ranking": []gin.H{},
	})
}
