package handler

import (
	"net/http"
	"strconv"
	"time"

	WorkSessionModel "app/internal/app/work-session/model"
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

func (handler *WorkSessionHandler) UpdateWorkSessionClocking(c *gin.Context) {
	var req WorkSessionModel.WorkSessionUpdate

	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Config.ErrorMessages()["NO_CLAIMS"]})
		return
	}

	userUUID := claims.(*AuthService.Claims).UUID

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.IsClocked == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "is_clocked field is required"})
		return
	}

	Response, registerErr := handler.service.UpdateWorkSessionClocking(WorkSessionModel.WorkSessionUpdate{
		UserUUID:  userUUID,
		IsClocked: req.IsClocked,
	})

	if registerErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": registerErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, Response)
}

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

func (handler *WorkSessionHandler) GetWorkSessionHistory(c *gin.Context) {
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Config.ErrorMessages()["NO_CLAIMS"]})
		return
	}

	userUUID := ""

	// if the param user_uuid is not present, use the uuid from the claims
	if c.Query("user_uuid") != "" {
		userUUID = c.Query("user_uuid")
	} else {
		userUUID = claims.(*AuthService.Claims).UUID
	}

	if userUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_uuid is required"})
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	if !handler.isValidISO8601(startDate) || !handler.isValidISO8601(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date must be valid ISO 8601 timestamps"})
		return
	}

	nowTimestamp := time.Now().Format(time.RFC3339Nano)
	twoYearsAgo := time.Now().AddDate(-2, 0, 0).Format(time.RFC3339Nano)

	if !handler.isValidISO8601(startDate) || !handler.isValidISO8601(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date must be valid ISO 8601 timestamps"})
		return
	}

	if startDate < twoYearsAgo || endDate > nowTimestamp {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date range cannot exceed 2 years from the current date"})
		return
	}

	if limit <= 0 {
		limit = 50
	}

	if offset < 0 {
		offset = 0
	}

	history, err := handler.service.GetWorkSessionHistory(userUUID, startDate, endDate, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if history == nil {
		history = []WorkSessionModel.WorkSessionReadHistory{}
	}

	c.JSON(http.StatusOK, history)

}
