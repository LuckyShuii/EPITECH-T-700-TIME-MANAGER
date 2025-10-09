package service

import (
	"app/internal/app/user/model"
	"app/internal/app/user/repository"
)

type UserService interface {
	GetUsers() ([]model.UserRead, error)
	GetUserByEmailAuth(typeOf string, data string) (*model.UserReadJWT, error)
	RegisterUser(user model.UserCreate) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (service *userService) GetUsers() ([]model.UserRead, error) {
	return service.repo.FindAll()
}

func (service *userService) GetUserByEmailAuth(typeOf string, data string) (*model.UserReadJWT, error) {
	return service.repo.FindByTypeAuth(typeOf, data)
}

func (service *userService) RegisterUser(user model.UserCreate) error {
	return service.repo.RegisterUser(user)
}
