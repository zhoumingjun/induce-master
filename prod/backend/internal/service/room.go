package service

type RoomService struct {
	repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) Create(room *Room) error {
	return s.repo.Create(room)
}

func (s *RoomService) GetByID(id string) (*Room, error) {
	return s.repo.GetByID(id)
}

func (s *RoomService) List() ([]Room, error) {
	return s.repo.List()
}

func (s *RoomService) Join(roomID, userID string) error {
	return nil
}
