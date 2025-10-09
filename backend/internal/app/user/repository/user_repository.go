package repository

import (
	"app/internal/app/user/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]model.UserRead, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (repo *userRepository) FindAll() ([]model.UserRead, error) {
	var users []model.UserRead
	err := repo.db.Raw("SELECT uuid, username, email, first_name, last_name, phone_number, roles, created_at, updated_at FROM users").Scan(&users).Error
	return users, err
}
