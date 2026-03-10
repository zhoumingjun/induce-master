package model

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
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

// DBRoom 数据库房间模型
type DBRoom struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string    `gorm:"size:20" json:"name"`
	OwnerID   string    `gorm:"type:varchar(36)" json:"owner_id"`
	Password  string    `gorm:"size:6" json:"-"`
	Status    int       `gorm:"default:0" json:"status"` // 0: waiting, 1: playing
	Settings  string    `gorm:"type:json" json:"settings"`
	CreatedAt time.Time `json:"created_at"`
}

func (DBRoom) TableName() string {
	return "rooms"
}

// DBRoomPlayer 数据库房间玩家模型
type DBRoomPlayer struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	RoomID    string    `gorm:"type:varchar(36);index" json:"room_id"`
	UserID    string    `gorm:"type:varchar(36);index" json:"user_id"`
	Ready     bool      `gorm:"default:false" json:"ready"`
	Score     int       `gorm:"default:0" json:"score"`
	JoinedAt  time.Time `json:"joined_at"`
}

func (DBRoomPlayer) TableName() string {
	return "room_players"
}

// DBGame 数据库游戏模型
type DBGame struct {
	ID            string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	RoomID        string     `gorm:"type:varchar(36);index" json:"room_id"`
	PlayerAID     string     `gorm:"type:varchar(36)" json:"player_a_id"`
	PlayerBID     string     `gorm:"type:varchar(36)" json:"player_b_id"`
	WordA         string     `gorm:"size:50" json:"word_a"`
	WordB         string     `gorm:"size:50" json:"word_b"`
	WinnerID      string     `gorm:"type:varchar(36)" json:"winner_id"`
	ScoreA        int        `gorm:"default:0" json:"score_a"`
	ScoreB        int        `gorm:"default:0" json:"score_b"`
	Messages      string     `gorm:"type:json" json:"messages"`
	CurrentRound  int        `gorm:"default:1" json:"current_round"`
	Status        int        `gorm:"default:0" json:"status"` // 0: playing, 1: finished
	CreatedAt     time.Time  `json:"created_at"`
	FinishedAt    *time.Time `json:"finished_at"`
}

func (DBGame) TableName() string {
	return "games"
}

func InitDB() (*gorm.DB, error) {
	// 使用 SQLite 本地数据库
	dsn := "induce_master.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&User{}, &DBRoom{}, &DBRoomPlayer{}, &DBGame{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return db, nil
}
