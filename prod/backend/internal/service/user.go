package service

import (
	"induce-master/internal/config"
	"induce-master/internal/model"
	"induce-master/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	repo   *repository.UserRepository
	config *config.Config
}

func NewUserService(repo *repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		repo:   repo,
		config: cfg,
	}
}

func (s *UserService) Create(user *model.User) error {
	return s.repo.Create(user)
}

func (s *UserService) GetByUsername(username string) (*model.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *UserService) GetByID(id string) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
