package handler

import (
	"net/http"
	"strconv"

	authService "app/internal/app/auth/service"
	"app/internal/app/user/model"
	"app/internal/config"

	"github.com/gin-gonic/gin"

	"app/internal/app/common/response"
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

// LoginHandler authenticates a user using email or username and password.
//
// @Summary      Login a user
// @Description  Authenticates a user via email or username and returns a JWT cookie on success.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "Login credentials"
// @Success      200   {object}  response.MessageResponse  "logged in successfully"
// @Router       /auth/login [post]
func (handler *AuthHandler) LoginHandler(c *gin.Context) {
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
	if err != nil || user.UUID == "" || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	jwtExpirationHours, err := strconv.Atoi(config.LoadConfig().JWTExpirationHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid JWT expiration configuration"})
		return
	}
	expiration := jwtExpirationHours * 3600

	c.SetSameSite(http.SameSiteStrictMode)

	c.SetCookie(
		"token",
		token,
		expiration,
		"/",
		config.LoadConfig().FrontendURL,
		config.LoadConfig().ProjectStatus == "PROD",
		true,
	)

	c.JSON(http.StatusOK, response.MessageResponse{Message: "logged in successfully"})
}

// MeHandler retrieves the current logged-in user's information.
//
// @Summary      Get current user info
// @Description  Retrieves information about the currently logged-in user based on the JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  model.UserMeJWT  "Current user information"
// @Router       /auth/me [get]
func (handler *AuthHandler) MeHandler(c *gin.Context) {
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

	c.JSON(http.StatusOK, model.UserMeJWT{
		UserUUID:    claims.UUID,
		Roles:       claims.Roles,
		Email:       claims.Email,
		Username:    claims.Username,
		FirstName:   claims.FirstName,
		LastName:    claims.LastName,
		PhoneNumber: claims.PhoneNumber,
	})
}

// LogoutHandler logs out the current user by clearing the JWT cookie.
//
// @Summary      Logout the current user
// @Description  Clears the JWT cookie and invalidates the user's session.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.MessageResponse  "Logout successful"
// @Router       /auth/logout [post]
func (handler *AuthHandler) LogoutHandler(c *gin.Context) {
	c.SetCookie(
		"token",
		"",
		-1,
		"/",
		config.LoadConfig().FrontendURL,
		config.LoadConfig().ProjectStatus == "PROD",
		true,
	)

	c.JSON(http.StatusOK, response.MessageResponse{Message: "logged out successfully"})
}
