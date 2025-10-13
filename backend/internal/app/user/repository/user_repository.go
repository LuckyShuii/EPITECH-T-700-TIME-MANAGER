package repository

import (
	"app/internal/app/user/model"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]model.UserRead, error)
	FindByTypeAuth(typeOf string, data string) (*model.UserReadJWT, error)
	RegisterUser(user model.UserCreate) error
	FindIdByUuid(id string) (userId int, err error)
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

func (repo *userRepository) FindIdByUuid(uuid string) (userId int, err error) {
	err = repo.db.Raw("SELECT id FROM users WHERE uuid = ?", uuid).Scan(&userId).Error
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (repo *userRepository) FindByTypeAuth(typeOf string, data string) (*model.UserReadJWT, error) {
	var user model.UserReadJWT

	if typeOf != "email" && typeOf != "username" {
		return nil, fmt.Errorf("invalid type: %s", typeOf)
	}

	query := fmt.Sprintf("SELECT id, uuid, email, roles, first_name, last_name, username, phone_number, password_hash FROM users WHERE %s = ?", typeOf)

	err := repo.db.Raw(query, data).Scan(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) RegisterUser(user model.UserCreate) error {
	err := repo.db.Exec(
		"INSERT INTO users (uuid, first_name, last_name, email, username, phone_number, roles, password_hash) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.UUID, user.FirstName, user.LastName, user.Email, user.Username, user.PhoneNumber, user.Roles, user.PasswordHash,
	).Error
	return err
}
