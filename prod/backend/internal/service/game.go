package service

import (
	"induce-master/internal/model"
	"induce-master/internal/repository"
)

type GameService struct {
	repo *repository.GameRepository
}

func NewGameService(repo *repository.GameRepository) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) Create(game *model.DBGame) error {
	return s.repo.Create(game)
}

func (s *GameService) GetByID(id string) (*model.DBGame, error) {
	return s.repo.GetByID(id)
}

func (s *GameService) Update(game *model.DBGame) error {
	return s.repo.Update(game)
}
