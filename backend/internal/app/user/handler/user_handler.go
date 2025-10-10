package handler

import (
	"net/http"

	"app/internal/app/user/model"
	"app/internal/app/user/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (handler *UserHandler) GetUsers(c *gin.Context) {
	users, err := handler.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (handler *UserHandler) RegisterUser(c *gin.Context) {
	var req model.UserCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Email == "" || req.Password == "" || req.Username == "" || req.FirstName == "" || req.LastName == "" || len(req.Roles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}

	registerErr := handler.service.RegisterUser(req)
	if registerErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": registerErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}
