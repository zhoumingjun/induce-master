package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:20" json:"username"`
	DisplayName string    `gorm:"size:50" json:"display_name"`
	PasswordHash string    `gorm:"size:255" json:"-"`
	AvatarURL    string    `gorm:"size:255" json:"avatar_url"`
	Rank         int       `gorm:"default:0" json:"rank"`
	Token        string    `gorm:"size:255" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type Room struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string    `gorm:"size:20" json:"name"`
	OwnerID   string    `gorm:"type:varchar(36)" json:"owner_id"`
	Password  string    `gorm:"size:6" json:"-"`
	Status    int       `gorm:"default:0" json:"status"` // 0: waiting, 1: playing
	Settings  string    `gorm:"type:json" json:"settings"`
	CreatedAt time.Time `json:"created_at"`
}

func (Room) TableName() string {
	return "rooms"
}

type RoomPlayer struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	RoomID    string    `gorm:"type:varchar(36);index" json:"room_id"`
	UserID    string    `gorm:"type:varchar(36);index" json:"user_id"`
	Ready     bool      `gorm:"default:false" json:"ready"`
	Score     int       `gorm:"default:0" json:"score"`
	JoinedAt  time.Time `json:"joined_at"`
}

func (RoomPlayer) TableName() string {
	return "room_players"
}

type Game struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	RoomID     string    `gorm:"type:varchar(36);index" json:"room_id"`
	PlayerAID  string    `gorm:"type:varchar(36)" json:"player_a_id"`
	PlayerBID  string    `gorm:"type:varchar(36)" json:"player_b_id"`
	WordA      string    `gorm:"size:50" json:"word_a"`
	WordB      string    `gorm:"size:50" json:"word_b"`
	WinnerID   string    `gorm:"type:varchar(36)" json:"winner_id"`
	ScoreA     int       `gorm:"default:0" json:"score_a"`
	ScoreB     int       `gorm:"default:0" json:"score_b"`
	Messages   string    `gorm:"type:json" json:"messages"`
	CurrentRound int    `gorm:"default:1" json:"current_round"`
	Status     int       `gorm:"default:0" json:"status"` // 0: playing, 1: finished
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

func (Game) TableName() string {
	return "games"
}

func InitDB(cfg *Config) (*gorm.DB, error) {
	// 简化版，实际需要连接 MySQL
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	// return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return nil, nil // 占位，实际需要数据库
}
