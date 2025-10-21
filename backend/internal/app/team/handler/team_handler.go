package handler

import (
	"net/http"

	"app/internal/app/team/service"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	service service.TeamService
}

func NewTeamHandler(service service.TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

// GetTeams retrieves all teams.
//
// @Summary      Get all teams
// @Description  Returns a list of all registered teams & their members. 🔒 Requires role: **admin**
// @Tags         Teams
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   model.TeamReadAll  "List of teams retrieved successfully"
// @Router       /teams [get]
func (handler *TeamHandler) GetTeams(c *gin.Context) {
	teams, err := handler.service.GetTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teams)
}
