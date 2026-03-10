package repository

import (
	"induce-master/internal/model"
	"induce-master/internal/service"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *service.User) error {
	m := &model.User{
		ID:           user.ID,
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		PasswordHash: user.PasswordHash,
		AvatarURL:    user.AvatarURL,
		Rank:         user.Rank,
	}
	return r.db.Create(m).Error
}

func (r *UserRepository) GetByUsername(username string) (*service.User, error) {
	var m model.User
	if err := r.db.Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return &service.User{
		ID:           m.ID,
		Username:     m.Username,
		DisplayName:  m.DisplayName,
		PasswordHash: m.PasswordHash,
		AvatarURL:    m.AvatarURL,
		Rank:         m.Rank,
	}, nil
}

func (r *UserRepository) GetByID(id string) (*service.User, error) {
	var m model.User
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &service.User{
		ID:           m.ID,
		Username:     m.Username,
		DisplayName:  m.DisplayName,
		PasswordHash: m.PasswordHash,
		AvatarURL:    m.AvatarURL,
		Rank:         m.Rank,
	}, nil
}
