package handler

import (
	"net/http"

	WorkSessionModel "app/internal/app/work-session/model"
	WorkSessionService "app/internal/app/work-session/service"

	AuthService "app/internal/app/auth/service"

	"github.com/gin-gonic/gin"
)

type WorkSessionHandler struct {
	service WorkSessionService.WorkSessionService
}

func NewWorkSessionHandler(service WorkSessionService.WorkSessionService) *WorkSessionHandler {
	return &WorkSessionHandler{service: service}
}

func (handler *WorkSessionHandler) UpdateWorkSessionClocking(c *gin.Context) {
	var req WorkSessionModel.WorkSessionUpdate

	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
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
