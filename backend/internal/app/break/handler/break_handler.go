package handler

import (
	"net/http"

	BreakModel "app/internal/app/break/model"
	BreakService "app/internal/app/break/service"

	"github.com/gin-gonic/gin"

	Config "app/internal/config"
)

type BreakHandler struct {
	service BreakService.BreakService
}

func NewBreakHandler(service BreakService.BreakService) *BreakHandler {
	return &BreakHandler{service: service}
}

func (handler *BreakHandler) UpdateBreak(c *gin.Context) {
	var req BreakModel.BreakUpdate

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Config.ErrorMessages()["INVALID_REQUEST"]})
		return
	}

	if req.IsBreaking == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "is_breaking field is required"})
		return
	}

	Response, registerErr := handler.service.UpdateBreakClocking(BreakModel.BreakUpdate{
		WorkSessionUUID: req.WorkSessionUUID,
		IsBreaking:      req.IsBreaking,
	})

	if registerErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": registerErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, Response)
}
