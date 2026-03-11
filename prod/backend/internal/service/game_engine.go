package service

import (
	"math/rand"
	"time"

	"induce-master/internal/model"
)

// GameEngine 游戏引擎
type GameEngine struct {
	RoomID     string
	Players    []string // 用户ID列表
	Words      map[string]string // userID -> 关键词
	Scores     map[string]int
	Round      int
	MaxRounds  int
	TimeLimit  int // 每轮时间(秒)
	CurrentTurn int // 当前发言玩家索引
	Messages   []GameMessage
	Status     GameStatus
}

// GameMessage 游戏消息
type GameMessage struct {
	Round     int       `json:"round"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Time      time.Time `json:"time"`
	IsKeyword bool      `json:"is_keyword"` // 是否说出关键词
}

// GameStatus 游戏状态
type GameStatus int

const (
	GameStatusWaiting  GameStatus = 0
	GameStatusPlaying  GameStatus = 1
	GameStatusFinished GameStatus = 2
)

// NewGameEngine 创建游戏引擎
func NewGameEngine(roomID string, players []string, maxRounds int, timeLimit int) *GameEngine {
	return &GameEngine{
		RoomID:    roomID,
		Players:   players,
		Words:     make(map[string]string),
		Scores:    make(map[string]int),
		Round:     1,
		MaxRounds: maxRounds,
		TimeLimit: timeLimit,
		Messages:  make([]GameMessage, 0),
		Status:    GameStatusWaiting,
	}
}

// AssignWords 分配关键词
func (g *GameEngine) AssignWords() {
	// 每个玩家分配不同的关键词
	categories := []string{"水果", "动物", "城市", "电影", "成语"}
	
	for i, playerID := range g.Players {
		// 轮流使用不同分类
		category := categories[i%len(categories)]
		words := model.Words[category]
		word := words[rand.Intn(len(words))]
		g.Words[playerID] = word
	}
}

// GetWord 获取指定玩家的关键词
func (g *GameEngine) GetWord(userID string) string {
	return g.Words[userID]
}

// GetOpponentWord 获取对手的关键词（用于判定）
func (g *GameEngine) GetOpponentWord(userID string) string {
	for pid, word := range g.Words {
		if pid != userID {
			return word
		}
	}
	return ""
}

// ProcessMessage 处理游戏消息，返回是否说出关键词
func (g *GameEngine) ProcessMessage(userID, username, content string) *GameMessage {
	msg := GameMessage{
		Round:    g.Round,
		UserID:   userID,
		Username: username,
		Content:  content,
		Time:     time.Now(),
	}

	// 检查是否说出对手的关键词（严格匹配）
	opponentWord := g.GetOpponentWord(userID)
	if opponentWord != "" && containsKeyword(content, opponentWord) {
		msg.IsKeyword = true
		// 说出关键词扣分
		g.Scores[userID] -= 10
	}

	g.Messages = append(g.Messages, msg)
	return &msg
}

// containsKeyword 检查消息是否严格包含关键词
// 严格匹配：消息中必须完整包含关键词（可以有前后文）
func containsKeyword(content, keyword string) bool {
	// 精确匹配：内容必须包含关键词完整文本
	contentClean := normalizeString(content)
	keywordClean := normalizeString(keyword)
	
	if len(contentClean) == 0 || len(keywordClean) == 0 {
		return false
	}
	
	// 检查是否包含完整关键词
	return findSubstring(contentClean, keywordClean)
}

// normalizeString 规范化字符串（移除空格和标点，统一大小写）
func normalizeString(s string) string {
	result := make([]rune, 0)
	for _, r := range s {
		// 只保留中文、字母、数字
		if (r >= '\u4e00' && r <= '\u9fff') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			// 字母转小写
			if r >= 'A' && r <= 'Z' {
				r = r + 32
			}
			result = append(result, r)
		}
	}
	return string(result)
}

// findSubstring 查找子串
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// NextTurn 下一回合
func (g *GameEngine) NextTurn() {
	g.CurrentTurn = (g.CurrentTurn + 1) % len(g.Players)
	
	// 如果回到第一个玩家，增加回合数
	if g.CurrentTurn == 0 {
		g.Round++
		if g.Round > g.MaxRounds {
			g.Status = GameStatusFinished
		}
	}
}

// Start 开始游戏
func (g *GameEngine) Start() {
	g.Status = GameStatusPlaying
	g.AssignWords()
}

// GetCurrentPlayer 获取当前玩家
func (g *GameEngine) GetCurrentPlayer() string {
	if len(g.Players) == 0 {
		return ""
	}
	return g.Players[g.CurrentTurn]
}

// IsFinished 检查游戏是否结束
func (g *GameEngine) IsFinished() bool {
	return g.Status == GameStatusFinished
}

// GetWinner 获取获胜者
func (g *GameEngine) GetWinner() string {
	if len(g.Players) == 0 {
		return ""
	}
	
	// 分数最高的获胜
	winner := g.Players[0]
	maxScore := g.Scores[winner]
	
	for pid, score := range g.Scores {
		if score > maxScore {
			maxScore = score
			winner = pid
		}
	}
	
	return winner
}

// GetScores 获取所有玩家分数
func (g *GameEngine) GetScores() map[string]int {
	return g.Scores
}

// GetMessages 获取游戏消息
func (g *GameEngine) GetMessages() []GameMessage {
	return g.Messages
}

// GetGameInfo 获取游戏信息
func (g *GameEngine) GetGameInfo() map[string]interface{} {
	info := make(map[string]interface{})
	info["room_id"] = g.RoomID
	info["round"] = g.Round
	info["max_rounds"] = g.MaxRounds
	info["time_limit"] = g.TimeLimit
	info["status"] = g.Status
	info["words"] = g.Words
	info["scores"] = g.Scores
	info["current_player"] = g.GetCurrentPlayer()
	return info
}
