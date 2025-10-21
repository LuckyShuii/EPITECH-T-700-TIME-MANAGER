package handler

import (
	"net/http"

	"app/internal/app/team/service"
	UserService "app/internal/app/user/service"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	service     service.TeamService
	UserService UserService.UserService
}

func NewTeamHandler(service service.TeamService, userService UserService.UserService) *TeamHandler {
	return &TeamHandler{service: service, UserService: userService}
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
// @Description  Returns a team and its members by the provided UUID. 🔒 Requires role: **any**
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

// @Description  Deletes a team by the provided UUID. 🔒 Requires role: **admin**
// @Tags         Teams
// @Security     BearerAuth
// @Param        uuid   path      string  true  "Team UUID"
// @Success      204    "Team deleted successfully"
// @Router       /teams/{uuid} [delete]
func (handler *TeamHandler) DeleteTeamByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	teamID, err := handler.service.GetIdByUuid(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = handler.service.DeleteTeamByID(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Description  Removes a user from a team by the provided team UUID and user UUID. 🔒 Requires role: **admin**
// @Tags         Teams
// @Security     BearerAuth
// @Param        team_uuid   path      string  true  "Team UUID"
// @Param        user_uuid   path      string  true  "User UUID"
// @Success      204    "User removed from team successfully"
// @Router       /teams/users/{team_uuid}/{user_uuid} [delete]
func (handler *TeamHandler) RemoveUserFromTeam(c *gin.Context) {
	teamUUID := c.Param("team_uuid")
	userUUID := c.Param("user_uuid")

	teamID, err := handler.service.GetIdByUuid(teamUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, err := handler.UserService.GetIdByUuid(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = handler.service.RemoveUserFromTeam(teamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
