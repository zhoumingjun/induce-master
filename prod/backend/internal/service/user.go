package service

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	ID           string
	Username     string
	DisplayName  string
	PasswordHash string
	AvatarURL    string
	Rank         int
}

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

func (s *UserService) Create(user *User) error {
	return s.repo.Create(user)
}

func (s *UserService) GetByUsername(username string) (*User, error) {
	return s.repo.GetByUsername(username)
}

func (s *UserService) GetByID(id string) (*User, error) {
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
