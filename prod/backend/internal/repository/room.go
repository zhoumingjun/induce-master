package repository

import (
	"induce-master/internal/model"

	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *model.DBRoom) error {
	return r.db.Create(room).Error
}

func (r *RoomRepository) GetByID(id string) (*model.DBRoom, error) {
	var room model.DBRoom
	err := r.db.First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepository) Update(room *model.DBRoom) error {
	return r.db.Save(room).Error
}

func (r *RoomRepository) List() ([]*model.DBRoom, error) {
	var rooms []*model.DBRoom
	err := r.db.Where("status = ?", 0).Order("created_at DESC").Find(&rooms).Error
	return rooms, err
}

func (r *RoomRepository) AddPlayer(player *model.DBRoomPlayer) error {
	return r.db.Create(player).Error
}

func (r *RoomRepository) GetPlayers(roomID string) ([]*model.DBRoomPlayer, error) {
	var players []*model.DBRoomPlayer
	err := r.db.Where("room_id = ?", roomID).Find(&players).Error
	return players, err
}

func (r *RoomRepository) RemovePlayer(roomID, userID string) error {
	return r.db.Where("room_id = ? AND user_id = ?", roomID, userID).Delete(&model.DBRoomPlayer{}).Error
}

func (r *RoomRepository) UpdatePlayer(player *model.DBRoomPlayer) error {
	return r.db.Save(player).Error
}
