package model

import "time"

// Room 房间
type Room struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	OwnerID     string       `json:"owner_id"`
	MaxPlayers  int          `json:"max_players"`
	Status      RoomStatus   `json:"status"` // waiting, playing, finished
	Difficulty  string       `json:"difficulty"`
	CreatedAt   time.Time    `json:"created_at"`
	Players     []*RoomPlayer `json:"players"`
}

// RoomStatus 房间状态
type RoomStatus string

const (
	RoomStatusWaiting  RoomStatus = "waiting"
	RoomStatusPlaying  RoomStatus = "playing"
	RoomStatusFinished RoomStatus = "finished"
)

// RoomPlayer 房间玩家
type RoomPlayer struct {
	ID       string    `json:"id"`
	RoomID   string    `json:"room_id"`
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Ready    bool      `json:"ready"`
	JoinedAt time.Time `json:"joined_at"`
}

// Game 游戏
type Game struct {
	ID         string     `json:"id"`
	RoomID     string     `json:"room_id"`
	Status     GameStatus `json:"status"`
	Round      int        `json:"round"`
	MaxRounds  int        `json:"max_rounds"`
	TimeLimit  int        `json:"time_limit"` // 秒
	Words      map[string]string `json:"words"` // userID -> keyword
	CreatedAt  time.Time `json:"created_at"`
	Messages   []*GameMessage `json:"messages"`
	Scores     map[string]int `json:"scores"`
}

// GameStatus 游戏状态
type GameStatus string

const (
	GameStatusWaiting GameStatus = "waiting"
	GameStatusPlaying GameStatus = "playing"
	GameStatusFinished GameStatus = "finished"
)

// GameMessage 游戏消息
type GameMessage struct {
	ID        string    `json:"id"`
	GameID    string    `json:"game_id"`
	Round     int       `json:"round"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Word 关键词
type Word struct {
	ID         string   `json:"id"`
	Category   string   `json:"category"`
	Content    string   `json:"content"`
	Difficulty string   `json:"difficulty"`
}

// WordCategory 词库分类
var WordCategories = []string{"水果", "动物", "城市", "电影", "成语"}

// Words 词库
var Words = map[string][]string{
	"水果":  {"苹果", "香蕉", "葡萄", "车厘子", "西瓜", "草莓", "橙子", "桃子", "梨", "芒果"},
	"动物":  {"猫", "狗", "兔子", "考拉", "熊猫", "老虎", "狮子", "大象", "长颈鹿", "企鹅"},
	"城市":  {"北京", "上海", "成都", "纽约", "东京", "巴黎", "伦敦", "悉尼", "首尔", "香港"},
	"电影":  {"阿凡达", "复仇者联盟", "盗梦空间", "泰坦尼克号", "哈利波特", "星球大战", "指环王", "速度与激情"},
	"成语":  {"画蛇添足", "掩耳盗铃", "亡羊补牢", "守株待兔", "刻舟求剑", "叶公好龙", "井底之蛙", "胸有成竹"},
}
