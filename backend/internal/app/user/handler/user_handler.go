package handler

import (
	"net/http"

	"app/internal/app/user/model"
	"app/internal/app/user/service"

	"github.com/gin-gonic/gin"

	Config "app/internal/config"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": Config.ErrorMessages()["INVALID_REQUEST"]})
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

func (handler *UserHandler) DeleteUser(c *gin.Context) {
	var req struct {
		UserUUID string `json:"user_uuid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Config.ErrorMessages()["INVALID_REQUEST"]})
		return
	}

	if req.UserUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user_uuid field"})
		return
	}

	deleteErr := handler.service.DeleteUser(req.UserUUID)
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": deleteErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// UpdateUserStatus
func (handler *UserHandler) UpdateUserStatus(c *gin.Context) {
	// status is either active, disabled or pending
	var req struct {
		UserUUID string `json:"user_uuid"`
		Status   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Config.ErrorMessages()["INVALID_REQUEST"]})
		return
	}

	if req.UserUUID == "" || req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}

	if req.Status != "active" && req.Status != "disabled" && req.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status value"})
		return
	}

	err := handler.service.UpdateUserStatus(req.UserUUID, req.Status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "user status updated successfully",
		"new_status": req.Status,
	})
}
