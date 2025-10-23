package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"app/internal/app/work-session/model"
	WorkSessionService "app/internal/app/work-session/service"

	AuthService "app/internal/app/auth/service"

	Config "app/internal/config"

	"github.com/gin-gonic/gin"
)

type WorkSessionHandler struct {
	service WorkSessionService.WorkSessionService
}

func (handler *WorkSessionHandler) isValidISO8601(date string) bool {
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.999999",
	}

	for _, layout := range layouts {
		if _, err := time.Parse(layout, date); err == nil {
			return true
		}
	}
	return false
}

func NewWorkSessionHandler(service WorkSessionService.WorkSessionService) *WorkSessionHandler {
	return &WorkSessionHandler{service: service}
}

// UpdateWorkSessionClocking godoc
// @Summary      Update work session clocking status
// @Description  Starts or stops a work session for the authenticated user 🔒 Requires role: **any**
// @Tags         WorkSession
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      model.WorkSessionUpdate  true  "Clocking payload"
// @Success      201   {object}  model.WorkSessionUpdateResponse  "Work session updated successfully"
// @Router       /work-session/update-clocking [post]
func (handler *WorkSessionHandler) UpdateWorkSessionClocking(c *gin.Context) {
	var req model.WorkSessionUpdate

	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Config.ErrorMessages()["NO_CLAIMS"]})
		return
	}

	userUUID := claims.(*AuthService.Claims).UUID

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Config.ErrorMessages()["INVALID_REQUEST"]})
		return
	}

	if req.IsClocked == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "is_clocked field is required"})
		return
	}

	Response, registerErr := handler.service.UpdateWorkSessionClocking(model.WorkSessionUpdate{
		UserUUID:  userUUID,
		IsClocked: req.IsClocked,
	})

	if registerErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": registerErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, Response)
}

// GetWorkSessionStatus godoc
// @Summary      Get current work session status
// @Description  Returns whether the authenticated user is currently clocked in 🔒 Requires role: **any**
// @Tags         WorkSession
// @Security     BearerAuth
// @Produce      json
// @Success      200   {object}  model.WorkSessionStatus  "Current work session status"
// @Router       /work-session/status [get]
func (handler *WorkSessionHandler) GetWorkSessionStatus(c *gin.Context) {
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Config.ErrorMessages()["NO_CLAIMS"]})
		return
	}

	userUUID := claims.(*AuthService.Claims).UUID

	status, err := handler.service.GetWorkSessionStatus(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetWorkSessionHistory godoc
// @Summary      Get user work session history
// @Description  Returns paginated history of work sessions between a start and end date 🔒 Requires role: **admin or manager** to get history from any users. If the user is employee, only their own history can be accessed.
// @Tags         WorkSession
// @Security     BearerAuth
// @Produce      json
// @Param        user_uuid    query     string  false  "UUID of the user (optional, defaults to authenticated user)"
// @Param        start_date   query     string  true   "Start date in ISO 8601 format (from 2 years ago up to now)"
// @Param        end_date     query     string  true   "End date in ISO 8601 format (can't be in the future)"
// @Param        limit        query     int     false  "Number of results to return (default 50)"
// @Param        offset       query     int     false  "Pagination offset (default 0)"
// @Success      200   {array}  model.WorkSessionReadHistory  "List of work session history entries"
// @Router       /work-session/history [get]
func (handler *WorkSessionHandler) GetWorkSessionHistory(c *gin.Context) {
	userUUID, err := handler.getUserUUIDFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	startDate, endDate, err := handler.validateDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, offset := handler.parsePaginationParams(c)

	history, err := handler.service.GetWorkSessionHistory(userUUID, startDate, endDate, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if history == nil {
		history = []model.WorkSessionReadHistory{}
	}

	c.JSON(http.StatusOK, history)
}

func (handler *WorkSessionHandler) getUserUUIDFromClaims(c *gin.Context) (string, error) {
	claims, exists := c.Get("userClaims")
	if !exists {
		return "", fmt.Errorf("%s", Config.ErrorMessages()["NO_CLAIMS"])
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

	if c.Query("user_uuid") != "" && isAdminOrManager {
		return c.Query("user_uuid"), nil
	}

	if authClaims.UUID == "" {
		return "", fmt.Errorf("user_uuid is required")
	}

	return authClaims.UUID, nil
}

func (handler *WorkSessionHandler) validateDateRange(c *gin.Context) (string, string, error) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if !handler.isValidISO8601(startDate) || !handler.isValidISO8601(endDate) {
		return "", "", fmt.Errorf("start_date and end_date must be valid ISO 8601 timestamps")
	}

	now := time.Now().Format(time.RFC3339Nano)
	twoYearsAgo := time.Now().AddDate(-2, 0, 0).Format(time.RFC3339Nano)

	if startDate < twoYearsAgo {
		return "", "", fmt.Errorf("date range cannot exceed 2 years from the current date")
	}

	if endDate > now {
		return "", "", fmt.Errorf("end_date cannot be in the future")
	}

	return startDate, endDate, nil
}

func (handler *WorkSessionHandler) parsePaginationParams(c *gin.Context) (int, int) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
