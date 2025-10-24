package handler

import (
	"app/internal/db"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"app/internal/app/user/model"
	"app/internal/app/user/service"

	AuthService "app/internal/app/auth/service"
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
// @Description  [Cache: 5sec] Returns a list of all registered users. 🔒 Requires role: **admin**
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   model.UserRead  "List of users retrieved successfully"
// @Router       /users [get]
func (handler *UserHandler) GetUsers(c *gin.Context) {
	// ⚠️ implement pagination - technical debt ⚠️
	ctx := c.Request.Context()

	cacheKey := "users_list"

	// Try to get users list from cache
	cachedUsers, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedUsers != "" {
		var users []model.UserRead
		if jsonErr := json.Unmarshal([]byte(cachedUsers), &users); jsonErr == nil {
			c.JSON(http.StatusOK, users)
			return
		}
	}

	users, err := handler.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cache the users list
	usersJSON, _ := json.Marshal(users)
	if setErr := db.RedisClient.Set(ctx, cacheKey, usersJSON, 5*time.Second).Err(); setErr != nil {
		log.Printf("⚠️ Error setting cache for %s : %v", cacheKey, setErr)
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
// @Router       /users/delete/:user_uuid [delete]
func (handler *UserHandler) DeleteUser(c *gin.Context) {
	userUUID := c.Param("uuid")

	if userUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user_uuid field"})
		return
	}

	deleteErr := handler.service.DeleteUser(userUUID)
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": deleteErr.Error()})
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

// GetUserByUUID retrieves a user by their UUID.
//
// @Summary      Get user by UUID
// @Description  [Cache: 5sec] Returns the details of a user identified by their UUID or not. If the UUID is not specificed it will return the current logged in users details. To query an other user data, must be manager or admin 🔒 Requires role: **all**
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        uuid   path      string  true  "User UUID"
// @Success      200    {object}  model.UserReadAll  "User details retrieved successfully"
// @Router       /users/{uuid} [get]
func (handler *UserHandler) GetUserByUUID(c *gin.Context) {
	UserUUID := c.Param("uuid")

	// if admin or manager, can get any user details
	// if normal user, can get only own details
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Config.ErrorMessages()["NO_CLAIMS"]})
		return
	}

	authClaims := claims.(*AuthService.Claims)
	roles := authClaims.Roles
	isAdminOrManager := false

	for _, role := range roles {
		if role == "admin" || role == "manager" {
			isAdminOrManager = true
			break
		}
	}

	if UserUUID != "" && !isAdminOrManager && UserUUID != authClaims.UUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to access other user details"})
		return
	}

	ctx := context.Background()

	cacheKey := "user:details:" + UserUUID

	// Try to get user details from cache
	cachedUser, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedUser != "" {
		var user model.UserReadAll
		if jsonErr := json.Unmarshal([]byte(cachedUser), &user); jsonErr == nil {
			c.JSON(http.StatusOK, user)
			return
		}
	}

	user, err := handler.service.GetUserByUUID(UserUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user.UUID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Cache the user details
	userJSON, _ := json.Marshal(user)
	if setErr := db.RedisClient.Set(ctx, cacheKey, userJSON, 5*time.Second).Err(); setErr != nil {
		log.Printf("⚠️ Error setting cache for %s : %v", cacheKey, setErr)
	}

	c.JSON(http.StatusOK, user)
}

// GetCurrentUserDashboardLayout retrieves the dashboard layout for the current user.
//
// @Summary      Get current user's dashboard layout
// @Description  Retrieve the dashboard layout configuration for the currently authenticated user. 🔒 Requires role: **all**
// @Tags         Users
// @Produce      json
// @Success      200  {object}   model.DashboardLayoutResponse
// @Router       /users/current-user-dashboard-layout [get]
func (handler *UserHandler) GetCurrentUserDashboardLayout(c *gin.Context) {
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Config.ErrorMessages()["NO_CLAIMS"]})
		return
	}

	authClaims := claims.(*AuthService.Claims)

	dashboardLayout, err := handler.service.GetUserDashboardLayout(authClaims.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if dashboardLayout.DashboardLayout == nil {
		c.JSON(http.StatusOK, gin.H{"layout": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"layout": dashboardLayout.DashboardLayout})
}

// DeleteCurrentUserDashboardLayout deletes the dashboard layout for the current user.
//
// @Summary      Delete current user's dashboard layout
// @Description  Delete the dashboard layout configuration for the currently authenticated user. 🔒 Requires role: **all**
// @Tags         Users
// @Produce      json
// @Success      200  "dashboard layout deleted successfully"
// @Router       /users/current-user-dashboard-layout/delete [delete]
func (handler *UserHandler) DeleteCurrentUserDashboardLayout(c *gin.Context) {
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Config.ErrorMessages()["NO_CLAIMS"]})
		return
	}

	authClaims := claims.(*AuthService.Claims)

	deleteErr := handler.service.DeleteUserDashboardLayout(authClaims.UUID)
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": deleteErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "dashboard layout deleted successfully"})
}

// UpdateCurrentUserDashboardLayout updates the dashboard layout for the current user.
//
// @Summary      Update current user's dashboard layout
// @Description  Update the dashboard layout configuration for the currently authenticated user. 🔒 Requires role: **all**
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body  body      model.UserDashboardLayoutUpdate  true  "Dashboard layout payload"
// @Success      200   "Dashboard layout updated successfully"
// @Router       /users/current-user-dashboard-layout/edit [put]
func (handler *UserHandler) UpdateCurrentUserDashboardLayout(c *gin.Context) {
	var req model.UserDashboardLayoutUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Config.ErrorMessages()["INVALID_REQUEST"]})
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Config.ErrorMessages()["NO_CLAIMS"]})
		return
	}

	authClaims := claims.(*AuthService.Claims)

	// Update the dashboard layout
	updateErr := handler.service.UpdateUserDashboardLayout(authClaims.UUID, req)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "dashboard layout updated successfully"})
}
