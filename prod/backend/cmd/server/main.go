package main

import (
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
	roomHandler := handler.NewRoomHandler(roomService)
	_ = handler.NewGameHandler(gameService) // TODO: 添加游戏相关路由

	// 初始化 WebSocket Hub
	hub := handler.NewHub(userService, roomService, gameService)
	go hub.Run()

	// 设置路由
	r := gin.Default()

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
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket升级失败: %v", err)
			return
		}

		// 注册到 Hub
		hub.Register <- &handler.Client{
			Conn:    conn,
			UserID:  userID,
			Send:    make(chan []byte, 256),
		}
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
