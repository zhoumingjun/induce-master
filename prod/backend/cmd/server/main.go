package main

import (
	"encoding/json"
	"log"
	"net/http"

	"induce-master/internal/config"
	"induce-master/internal/handler"
	"induce-master/internal/model"
	"induce-master/internal/repository"
	"induce-master/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSMessage WebSocket 消息结构
type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := model.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化仓储
	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	gameRepo := repository.NewGameRepository(db)

	// 初始化服务
	userService := service.NewUserService(userRepo, cfg)
	roomService := service.NewRoomService(roomRepo)
	gameService := service.NewGameService(gameRepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(userService, cfg)
	
	// 初始化 WebSocket Hub (需要在 roomHandler 之前)
	hub := handler.NewHub(userService, roomService, gameService)
	go hub.Run()
	
	roomHandler := handler.NewRoomHandler(roomService, hub)
	_ = handler.NewGameHandler(gameService) // TODO: 添加游戏相关路由

	// 设置路由
	r := gin.Default()
	
	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API 路由
	v1 := r.Group("/api/v1")
	{
		// 认证
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 房间
		rooms := v1.Group("/rooms")
		rooms.Use(authHandler.AuthMiddleware())
		{
			rooms.GET("", roomHandler.List)
			rooms.POST("", roomHandler.Create)
			rooms.GET("/:id", roomHandler.Get)
			rooms.POST("/:id/join", roomHandler.Join)
			rooms.POST("/:id/leave", roomHandler.Leave)
			rooms.POST("/:id/ready", roomHandler.Ready)
			rooms.POST("/:id/start", roomHandler.Start)
			rooms.POST("/:id/message", roomHandler.SendMessage)
			rooms.GET("/:id/status", roomHandler.GetGameStatus)
		}

		// 用户
		users := v1.Group("/users")
		users.Use(authHandler.AuthMiddleware())
		{
			users.GET("/me", authHandler.Me)
			users.GET("/ranking", authHandler.Ranking)
		}
	}

	// WebSocket
	r.GET("/ws", func(c *gin.Context) {
		userID := c.Query("user_id")
		token := c.Query("token")
		if userID == "" || token == "" {
			c.JSON(400, gin.H{"error": "缺少参数"})
			return
		}

		// 验证 token
		claims, err := userService.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "无效token"})
			return
		}

		if (*claims)["user_id"] != userID {
			c.JSON(401, gin.H{"error": "token不匹配"})
			return
		}

		// 升级为 WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket升级失败: %v", err)
			return
		}

		// 创建客户端
		client := &handler.Client{
			Conn:   conn,
			UserID: userID,
			Send:   make(chan []byte, 256),
		}

		// 注册到 Hub
		hub.Register <- client

		// 启动读和写 goroutine
		go writePump(client)
		go readPump(client, hub, userService, roomService)
	})

	// 启动服务器
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// writePump 向客户端发送消息
func writePump(client *handler.Client) {
	defer client.Conn.Close()
	for {
		message, ok := <-client.Send
		if !ok {
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// readPump 读取客户端消息
func readPump(client *handler.Client, hub *handler.Hub, userService *service.UserService, roomService *service.RoomService) {
	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket 错误: %v", err)
			}
			break
		}

		// 解析消息
		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("解析消息失败: %v", err)
			continue
		}

		// 处理消息
		handleWSMessage(client, hub, wsMsg, userService, roomService)
	}
}

// handleWSMessage 处理 WebSocket 消息
func handleWSMessage(client *handler.Client, hub *handler.Hub, msg WSMessage, userService *service.UserService, roomService *service.RoomService) {
	switch msg.Type {
	case "message":
		// 处理游戏消息
		var payload struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return
		}

		// 找到玩家所在的房间 - 优先使用客户端设置的 RoomID
		roomID := client.RoomID
		
		if roomID == "" {
			// 从数据库获取玩家所在的房间
			players, _ := roomService.GetPlayerRooms(client.UserID)
			for _, p := range players {
				if p.RoomID != "" {
					roomID = p.RoomID
					break
				}
			}
		}

		if roomID != "" {
			// 获取用户名
			username := client.UserID
			if user, err := userService.GetByID(client.UserID); err == nil {
				username = user.Username
			}

			// 处理游戏消息
			hub.ProcessGameMessage(roomID, client.UserID, username, payload.Content)
		}

	case "join_room":
		// 加入房间
		var payload struct {
			RoomID string `json:"room_id"`
		}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return
		}
		client.RoomID = payload.RoomID

	case "ping":
		// 响应 pong
		response, _ := json.Marshal(map[string]string{"type": "pong"})
		client.Send <- response
	}
}
