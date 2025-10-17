package handler

import (
	"log"
	"net/http"

	authService "app/internal/app/auth/service"
	"app/internal/config"
	"strconv"

	"github.com/gin-gonic/gin"

	"app/internal/app/user/model"
)

type LoginRequest struct {
	model.UserLogin
}

type AuthHandler struct {
	service authService.AuthService
}

func NewAuthHandler(service authService.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (handler *AuthHandler) LoginHandler(c *gin.Context) {
	/**
	 * Login Module
	 * Allow the user to log in with either their email or username and password
	 * Return a JWT token if successful
	 */
	var req LoginRequest

	var login struct {
		Login string `json:"login"`
		Type  string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.ErrorMessages()["INVALID_REQUEST"]})
		return
	}

	if req.Email == nil && req.Username == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or username required"})
		return
	}

	if req.Email == nil {
		login.Login = *req.Username
		login.Type = "username"
	} else {
		login.Login = *req.Email
		login.Type = "email"
	}

	user, token, err := handler.service.AuthenticateUser(login.Type, login.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if user.UUID == "" || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	jwtExpirationHours, err := strconv.Atoi(config.LoadConfig().JWTExpirationHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid JWT expiration configuration"})
		return
	}
	expiration := jwtExpirationHours * 3600

	c.SetCookie(
		"token",
		token,
		expiration,
		"/",
		config.LoadConfig().FrontendURL,
		config.LoadConfig().ProjectStatus == "PROD",
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "logged in"})
}

func (handler *AuthHandler) MeHandler(c *gin.Context) {
	/**
	 * Me Module
	 * Return the current logged in user based on the JWT token
	 */
	tokenStr, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token"})
		return
	}

	claims, err := handler.service.ValidateJWT(tokenStr)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_uuid":    claims.UUID,
		"roles":        claims.Roles,
		"email":        claims.Email,
		"username":     claims.Username,
		"first_name":   claims.FirstName,
		"last_name":    claims.LastName,
		"phone_number": claims.PhoneNumber,
	})

	log.Println("USER ID FROM CONTEXT:", c.GetString("userID"))
}

func (handler *AuthHandler) LogoutHandler(c *gin.Context) {
	/**
	 * Logout Module
	 * Clear the JWT token cookie
	 */
	c.SetCookie(
		"token",
		"",
		-1,
		"/",
		config.LoadConfig().FrontendURL,
		config.LoadConfig().ProjectStatus == "PROD",
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
