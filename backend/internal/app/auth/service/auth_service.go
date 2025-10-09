package service

import (
	"errors"
	"strconv"
	"time"

	userService "app/internal/app/user/service"

	"app/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

/**
 * {Claims} are the JWT parameters
 */
type Claims struct {
	UUID  string   `json:"uuid"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

type AuthService interface {
	AuthenticateUser(typeOf string, data string, password string) (string, error)
	GenerateJWT(uuid string, roles []string) (string, error)
}

type authService struct {
	userService userService.UserService
}

func NewAuthService(userService userService.UserService) AuthService {
	return &authService{userService}
}

func (service *authService) GenerateJWT(uuid string, roles []string) (string, error) {
	secret := config.LoadConfig().JWTSecret
	expiration := config.LoadConfig().JWTExpirationHours

	if secret == "" {
		return "", errors.New("missing JWT_SECRET")
	}

	// Convert expiration from string to int with Atoi from strconv package
	expirationInt, err := strconv.Atoi(expiration)
	if err != nil {
		return "", errors.New("invalid JWT expiration value")
	}

	// Set expiration time to current time + expiration hours
	expirationTime := time.Now().Add(time.Duration(expirationInt) * time.Hour)

	claims := &Claims{
		UUID:  uuid,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the JWT token with claims and sign it using the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (service *authService) ValidateJWT(tokenStr string) (*Claims, error) {
	secret := config.LoadConfig().JWTSecret
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Authenticates a user using the userService
func (service *authService) AuthenticateUser(typeOf string, data string, password string) (string, error) {
	// get user by email or username
	user, err := service.userService.GetUserByEmailAuth(typeOf, data)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := service.GenerateJWT(user.UUID, user.Roles)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
