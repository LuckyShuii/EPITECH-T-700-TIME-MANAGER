package service

import (
	"app/internal/app/user/model"
	"app/internal/app/user/repository"
	"fmt"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUsers() ([]model.UserRead, error)
	GetUserByEmailAuth(typeOf string, data string) (*model.UserReadJWT, error)
	RegisterUser(user model.UserCreate) error
	GetIdByUuid(id string) (int, error)
	UpdateUserStatus(userUUID string, status string) error
	DeleteUser(userUUID string) error
	UpdateUser(userID int, user model.UserUpdateEntry) error
	GetUserByUUID(userUUID string) (*model.UserReadAll, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (service *userService) GetUsers() ([]model.UserRead, error) {
	return service.repo.FindAll()
}

func (service *userService) GetUserByEmailAuth(typeOf string, data string) (*model.UserReadJWT, error) {
	return service.repo.FindByTypeAuth(typeOf, data)
}

func (service *userService) GetIdByUuid(id string) (int, error) {
	userID, err := service.repo.FindIdByUuid(id)
	return userID, err
}

func (service *userService) RegisterUser(user model.UserCreate) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	user.UUID = uuid.New().String()

	return service.repo.RegisterUser(user)
}

func (service *userService) DeleteUser(userUUID string) error {
	return service.repo.DeleteUser(userUUID)
}

func (service *userService) UpdateUserStatus(userUUID string, status string) error {
	return service.repo.UpdateUserStatus(userUUID, status)
}

func (service *userService) UpdateUser(userID int, user model.UserUpdateEntry) error {
	// user.username should not be trusted.
	if user.FirstName != nil && user.LastName != nil {
		user.Username = new(string)
		*user.Username = fmt.Sprintf("%c%s", unicode.ToLower(rune((*user.FirstName)[0])), strings.ToLower(*user.LastName))

	}

	return service.repo.UpdateUser(userID, user)
}

func (service *userService) GetUserByUUID(userUUID string) (*model.UserReadAll, error) {
	return service.repo.FindByUUID(userUUID)
}
