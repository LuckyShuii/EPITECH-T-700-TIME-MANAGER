package service

import (
	"app/internal/app/user/model"
	"app/internal/app/user/repository"
	"app/internal/config"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	MailTemplate "app/internal/app/mailer/template"

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
	ChangeUserPassword(token string, newPassword string) error
	ResetPassword(userEmail string, userUUID string) error
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

	defaultPasswordHash := config.LoadConfig().FixturesPassword

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash default password: %w", err)
	}
	hashedPasswordStr := string(hashedPassword)
	user.PasswordHash = &hashedPasswordStr

	err = service.repo.RegisterUser(user)
	if err != nil {
		return err
	}

	// create a JWT token for account activation with user uuid as key inside and 15 min
	secret := Config.LoadConfig().JWTSecret
	claims := jwt.MapClaims{
		"user_uuid": user.UUID,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
		"iat":       time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	activationToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return fmt.Errorf("failed to generate activation token: %w", err)
	}

	subject := "Hi " + user.FirstName + "! Welcome to TimeManager 🙂"

	body := MailTemplate.BaseMailTemplate(
		"Welcome to TimeManager",
		fmt.Sprintf(
			"Hello %s %s,<br><br>Welcome on board! To activate your account, please set your password by clicking the button below.<br><br>Your username is: <span style=\"font-weight:bold;\">%s</span><br><br>We are excited to have you with us!<br><br>Best regards,<br>The TimeManager Team",
			user.FirstName,
			user.LastName,
			user.Username,
		),
		"Activate my account",
		Config.LoadConfig().FrontendURL+"/activate-account?token="+activationToken,
	)

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

func (service *userService) ChangeUserPassword(token string, newPassword string) error {
	log.Println("Changing password for user token:", token)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	secret := Config.LoadConfig().JWTSecret
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !parsedToken.Valid {
		return fmt.Errorf("invalid or expired token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	userUUID, ok := claims["user_uuid"].(string)
	if !ok {
		return fmt.Errorf("user_uuid not found in token claims")
	}

	user, err := service.repo.FindByUUID(userUUID)
	if err != nil {
		return err
	}

	subject := "Hi " + user.FirstName + "! Your password has been changed 👍🏻"

	body := MailTemplate.BaseMailTemplate(
		"Password Changed Successfully",
		fmt.Sprintf(
			"Hello %s,<br><br>This is a confirmation that the password for your account %s has just been changed.<br><br>If you did not make this change, please contact our support team immediately.",
			user.FirstName,
			user.Email,
		),
		"",
		"",
	)

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

func (service *userService) ResetPassword(userEmail string, userUUID string) error {
	// Create JWT token for password reset with 15 min expiration
	secret := Config.LoadConfig().JWTSecret
	claims := jwt.MapClaims{
		"user_uuid": userUUID,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
		"iat":       time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	resetToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	link := Config.LoadConfig().FrontendURL + "/reset-password?token=" + resetToken

	subject := "Hi! Reset your TimeManager password 🔐"

	body := MailTemplate.BaseMailTemplate(
		"Reset Your Password",
		"Hello,<br><br>We received a request to reset the password for your account associated with this email address.<br><br>To reset your password please click the button below.",
		"Reset my password",
		link,
	)

	err = service.MailerService.Send(MailModel.Mail{
		To:      userEmail,
		Subject: subject,
		Body:    body,
	})

	if err != nil {
		return fmt.Errorf("failed to send reset password email: %w", err)
	}

	return nil
}
