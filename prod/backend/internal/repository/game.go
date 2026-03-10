package repository

import (
	"induce-master/internal/model"

	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) Create(game *model.DBGame) error {
	return r.db.Create(game).Error
}

func (r *GameRepository) GetByID(id string) (*model.DBGame, error) {
	var game model.DBGame
	err := r.db.First(&game, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *GameRepository) Update(game *model.DBGame) error {
	return r.db.Save(game).Error
}

func (r *GameRepository) GetByRoomID(roomID string) ([]*model.DBGame, error) {
	var games []*model.DBGame
	err := r.db.Where("room_id = ?", roomID).Order("created_at DESC").Find(&games).Error
	return games, err
}
