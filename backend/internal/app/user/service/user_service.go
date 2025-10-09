package service

import (
	"app/internal/app/user/model"
	"app/internal/app/user/repository"
)

type UserService interface {
	GetUsers() ([]model.UserRead, error)
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
