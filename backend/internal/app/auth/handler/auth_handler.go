package handler

import (
	"net/http"

	authService "app/internal/app/auth/service"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
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

	token, err := handler.service.AuthenticateUser(login.Type, login.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
