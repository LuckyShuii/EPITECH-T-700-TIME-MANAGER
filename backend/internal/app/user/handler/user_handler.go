package handler

import (
	"net/http"

	"app/internal/app/user/model"
	"app/internal/app/user/service"

	Config "app/internal/config"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetUsers retrieves all registered users.
//
// @Summary      Get all users
// @Description  Returns a list of all registered users. 🔒 Requires role: **admin**
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   model.UserRead  "List of users retrieved successfully"
// @Router       /users [get]
func (handler *UserHandler) GetUsers(c *gin.Context) {
	users, err := handler.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// RegisterUser registers a new user.
//
// @Summary      Register a new user
// @Description  Creates a new user with the provided information. 🔒 Requires role: **admin**
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      model.UserCreate  true  "User creation payload"
// @Success      201   {object}  response.UserCreatedResponse  "User registered successfully"
// @Router       /users/register [post]
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

// DeleteUser deletes an existing user by UUID.
//
// @Summary      Delete a user
// @Description  Deletes a user by their UUID. 🔒 Requires role: **admin**
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        uuid  path      string  true  "User UUID"
// @Success      200   {object}  response.UserDeletedResponse  "User deleted successfully"
// @Router       /users/{uuid} [delete]
func (handler *UserHandler) DeleteUser(c *gin.Context) {
	userUUID := c.Param("uuid")

	if userUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user UUID in URL"})
		return
	}

	if err := handler.service.DeleteUser(userUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// UpdateUserStatus updates the status of an existing user.
//
// @Summary      Update user status
// @Description  Updates a user's status (active, disabled, or pending). 🔒 Requires role: **admin**
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      model.UserStatusUpdatePayload  true  "User status update payload"
// @Success      200   {object}  response.UserStatusUpdatedResponse  "User status updated successfully"
// @Router       /users/update-status [put]
func (handler *UserHandler) UpdateUserStatus(c *gin.Context) {
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

// UpdateUser updates the details of an existing user.
//
// @Summary      Update user details
// @Description  Updates user information such as username, email, name, and roles. Only the **UUID** is mandatory. The **username** will be auto generated with first letter of the first name + last name on every update. 🔒 Requires role: **admin**
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      model.UserUpdateEntry  true  "User update payload"
// @Success      200   {object}  response.UserUpdatedResponse  "User updated successfully"
// @Router       /users [put]
func (handler *UserHandler) UpdateUser(c *gin.Context) {
	var req model.UserUpdateEntry
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Config.ErrorMessages()["INVALID_REQUEST"]})
		return
	}

	if req.UUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing uuid field"})
		return
	}

	userID, getErr := handler.service.GetIdByUuid(req.UUID)
	if getErr != nil || userID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	updateErr := handler.service.UpdateUser(userID, req)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}
