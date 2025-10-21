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

// GetTeamByUUID retrieves a team by its UUID.
//
// @Summary      Get team by UUID
// @Description  Returns a team and its members by the provided UUID. 🔒 Requires role: **admin**
// @Tags         Teams
// @Security     BearerAuth
// @Produce      json
// @Param        uuid   path      string  true  "Team UUID"
// @Success      200    {object}  model.TeamReadAll  "Team retrieved successfully"
// @Router       /teams/{uuid} [get]
func (handler *TeamHandler) GetTeamByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	team, err := handler.service.GetTeamByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, team)
}
