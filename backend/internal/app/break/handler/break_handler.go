package handler

import (
	"net/http"

	"app/internal/app/break/model"
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

// UpdateBreak updates the user's break status (start or end).
//
// @Summary      Update break status
// @Description  Starts or ends a break for the current work session depending on the value of `is_breaking`. 🔒 Requires role: **any**
// @Tags         WorkSession
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      model.BreakUpdate  true  "Break update payload"
// @Success      201   {object}  model.BreakUpdateResponse  "Break updated successfully"
// @Router       /work-session/update-breaking [post]
func (handler *BreakHandler) UpdateBreak(c *gin.Context) {
	var req model.BreakUpdate

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": Config.ErrorMessages()["INVALID_REQUEST"]})
		return
	}

	if req.IsBreaking == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "is_breaking field is required"})
		return
	}

	Response, registerErr := handler.service.UpdateBreakClocking(model.BreakUpdate{
		WorkSessionUUID: req.WorkSessionUUID,
		IsBreaking:      req.IsBreaking,
	})

	if registerErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": registerErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, Response)
}
