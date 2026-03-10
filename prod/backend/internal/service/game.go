package service

type GameService struct {
	repo *repository.GameRepository
}

func NewGameService(repo *repository.GameRepository) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) Create(game *Game) error {
	return s.repo.Create(game)
}

func (s *GameService) GetByID(id string) (*Game, error) {
	return s.repo.GetByID(id)
}

func (s *GameService) Update(game *Game) error {
	return s.repo.Update(game)
}
