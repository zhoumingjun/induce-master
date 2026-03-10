package service

import (
	"time"

	"induce-master/internal/model"
	"induce-master/internal/repository"
)

// DBRoom 房间数据库模型（导出给 handler 使用）
type DBRoom struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	Password  string    `json:"-"`
	Status    int       `json:"status"`
	Settings  string    `json:"settings"`
	CreatedAt time.Time `json:"created_at"`
}

type RoomService struct {
	repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) Create(room *DBRoom) error {
	dbRoom := &model.DBRoom{
		ID:        room.ID,
		Name:      room.Name,
		OwnerID:   room.OwnerID,
		Password:  room.Password,
		Status:    room.Status,
		Settings:  room.Settings,
		CreatedAt: time.Now(),
	}
	return s.repo.Create(dbRoom)
}

func (s *RoomService) GetByID(id string) (*model.DBRoom, error) {
	return s.repo.GetByID(id)
}

func (s *RoomService) List() ([]*model.DBRoom, error) {
	return s.repo.List()
}

func (s *RoomService) AddPlayer(roomID, userID string) error {
	player := &model.DBRoomPlayer{
		ID:       GenerateUUID(),
		RoomID:   roomID,
		UserID:   userID,
		Ready:    false,
		JoinedAt: time.Now(),
	}
	return s.repo.AddPlayer(player)
}

func (s *RoomService) RemovePlayer(roomID, userID string) error {
	return s.repo.RemovePlayer(roomID, userID)
}

func (s *RoomService) SetReady(roomID, userID string, ready bool) error {
	players, err := s.repo.GetPlayers(roomID)
	if err != nil {
		return err
	}
	for _, p := range players {
		if p.UserID == userID {
			p.Ready = ready
			return s.repo.UpdatePlayer(p)
		}
	}
	return nil
}

func (s *RoomService) GetPlayers(roomID string) ([]*model.DBRoomPlayer, error) {
	return s.repo.GetPlayers(roomID)
}

func (s *RoomService) UpdateRoomStatus(roomID string, status int) error {
	room, err := s.repo.GetByID(roomID)
	if err != nil {
		return err
	}
	room.Status = status
	return s.repo.Update(room)
}
