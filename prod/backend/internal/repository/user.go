package repository

import (
	"induce-master/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var m model.User
	if err := r.db.Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var m model.User
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) List() ([]*model.User, error) {
	var users []*model.User
	err := r.db.Order("rank DESC").Find(&users).Error
	return users, err
}
