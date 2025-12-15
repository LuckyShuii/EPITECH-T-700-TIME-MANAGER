package service

import (
	"app/internal/app/user/model"
	"app/internal/app/user/repository"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	MailModel "app/internal/app/mailer/model"
	MailerService "app/internal/app/mailer/service"

	WeeklyRateService "app/internal/app/weekly-rate/service"

	Config "app/internal/config"
)

type UserService interface {
	GetUsers() ([]model.UserRead, error)
	GetUserByEmailAuth(typeOf string, data string) (*model.UserReadJWT, error)
	RegisterUser(user model.UserCreate) error
	GetIdByUuid(id string) (int, error)
	UpdateUserStatus(userUUID string, status string) error
	DeleteUser(userUUID string) error
	DeleteUserDashboardLayout(userUUID string) error
	UpdateUser(userID int, user model.UserUpdateEntry) error
	GetUserByUUID(userUUID string) (*model.UserReadAll, error)
	SetWeeklyRateService(w WeeklyRateService.WeeklyRateService)
	GetUserDashboardLayout(userUUID string) (*model.UserDashboardLayout, error)
	UpdateUserDashboardLayout(userUUID string, layout model.UserDashboardLayoutUpdate) error
}

type userService struct {
	repo              repository.UserRepository
	WeeklyRateService WeeklyRateService.WeeklyRateService
	MailerService     MailerService.MailerService
}

func NewUserService(repo repository.UserRepository, mailerService MailerService.MailerService) UserService {
	return &userService{repo: repo, MailerService: mailerService}
}

func (s *userService) SetWeeklyRateService(w WeeklyRateService.WeeklyRateService) {
	s.WeeklyRateService = w
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
	user.UUID = uuid.New().String()

	if user.WeeklyRateUUID != nil {
		weeklyRateID, err := service.WeeklyRateService.GetIdByUuid(*user.WeeklyRateUUID)
		if err != nil {
			return fmt.Errorf("failed to fetch weekly rates: %w", err)
		}

		user.WeeklyRateUUID = nil
		user.WeeklyRateID = &weeklyRateID
	}

	if user.FirstDayOfWeek == nil {
		defaultDay := 1
		user.FirstDayOfWeek = &defaultDay
	}

	defaultPasswordHash := "$2a$10$FCvYkE0uB54aB/QykpqpOOavA7E4iDjEHeOB2xzW.Yl1b7/ThZuNq"
	user.PasswordHash = &defaultPasswordHash

	err := service.repo.RegisterUser(user)
	if err != nil {
		return err
	}

	subject := "Hi " + user.FirstName + "! Welcome to TimeManager 🙂"
	body := "Hello " + user.FirstName + ",\n\nWelcome on board! To activate your account, please set your password using the following link:\n\n" +
		Config.LoadConfig().FrontendURL + "/activate-account?user_public_key=" + user.UUID + "\n\nWe're excited to have you on board!\n\nBest regards,\nThe TimeManager Team"

	err = service.MailerService.Send(MailModel.Mail{
		To:      user.Email,
		Subject: subject,
		Body:    body,
	})

	if err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
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

	if user.WeeklyRateUUID != nil {
		// Check if the weekly rate exists
		weeklyRateID, err := service.WeeklyRateService.GetIdByUuid(*user.WeeklyRateUUID)
		if err != nil {
			return fmt.Errorf("failed to fetch weekly rates: %w", err)
		}

		user.WeeklyRateUUID = nil // Clear the UUID
		user.WeeklyRateID = &weeklyRateID
	}

	return service.repo.UpdateUser(userID, user)
}

func (service *userService) GetUserByUUID(userUUID string) (*model.UserReadAll, error) {
	return service.repo.FindByUUID(userUUID)
}

func (service *userService) GetUserDashboardLayout(userUUID string) (*model.UserDashboardLayout, error) {
	return service.repo.FindDashboardLayoutByUUID(userUUID)
}

func (service *userService) DeleteUserDashboardLayout(userUUID string) error {
	return service.repo.DeleteUserLayout(userUUID)
}

func (service *userService) UpdateUserDashboardLayout(userUUID string, layout model.UserDashboardLayoutUpdate) error {
	return service.repo.UpdateUserLayout(userUUID, layout)
}

func (service *userService) ChangeUserPassword(userUUID string, newPassword string) error {
	log.Println("Changing password for user UUID:", userUUID)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user, err := service.repo.FindByUUID(userUUID)
	if err != nil {
		return err
	}

	subject := "Your password has been changed"
	body := "Hello " + user.FirstName + ",\n\nThis is a confirmation that the password for your account " + user.Email + " has just been changed.\n\nIf you did not make this change, please contact our support team immediately."

	err = service.MailerService.Send(MailModel.Mail{
		To:      user.Email,
		Subject: subject,
		Body:    body,
	})

	if err != nil {
		return fmt.Errorf("failed to send password change email: %w", err)
	}

	userID, err := service.repo.FindIdByUuid(userUUID)
	if err != nil {
		return err
	}

	return service.repo.UpdateUserPassword(userID, string(hashedPassword))
}
