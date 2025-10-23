package handler

import (
	"net/http"

	"app/internal/app/team/model"
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

// @Summary      Delete team by UUID
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

	c.JSON(http.StatusNoContent, gin.H{"message": "Team deleted successfully"})
}

// @Description  Removes a user from a team by the provided team UUID and user UUID. 🔒 Requires role: **admin**
// @Tags         Teams
// @Security     BearerAuth
// @Summary      Remove user from team
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

	c.JSON(http.StatusNoContent, gin.H{"message": "User removed from team successfully"})
}

// CreateTeam creates a new team.
//
// @Summary      Create a new team
// @Description  Creates a new team with the provided details. You can add members or not. 🔒 Requires role: **admin**
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        team  body      model.TeamCreate  true  "Team to create"
// @Success      201 "Team created successfully"
// @Router       /teams [post]
func (handler *TeamHandler) CreateTeam(c *gin.Context) {
	var team model.TeamCreate
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if team.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team name is required"})
		return
	}

	if team.MemberUUIDs == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MemberUUIDs field is required. Cannot create an empty team"})
		return
	}

	if err := handler.service.CreateTeam(team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Team created successfully"})
}

// AddUsersToTeam adds users to an existing team.
//
// @Summary      Add users to a team
// @Description  Adds users to an existing team by their UUIDs. 🔒 Requires role: **admin**
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Param        team_users  body      model.TeamAddUsers  true  "Team UUID and User UUIDs to add"
// @Success      204 "Users added to team successfully"
// @Router       /teams/add-users [post]
func (handler *TeamHandler) AddUsersToTeam(c *gin.Context) {
	var req model.TeamAddUsers
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TeamUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team UUID is required"})
		return
	}

	if len(req.MemberUUIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one member UUID is required"})
		return
	}

	teamID, err := handler.service.GetIdByUuid(req.TeamUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var memberIDs []model.TeamMemberCreate
	for _, member := range req.MemberUUIDs {
		userID, err := handler.UserService.GetIdByUuid(member.UserUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		memberIDs = append(memberIDs, model.TeamMemberCreate{
			UserID:    userID,
			IsManager: member.IsManager,
		})
	}

	if err := handler.service.AddUsersToTeam(teamID, memberIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Users added to team successfully"})
}

// UpdateTeamByUUID updates an existing team by its UUID.
//
// @Summary      Update team by UUID
// @Description  Updates an existing team's details by its UUID. 🔒 Requires role: **admin**
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Param        uuid  path      string           true  "Team UUID"
// @Param        team  body      model.TeamUpdate  true  "Updated team details"
// @Success      200 "Team updated successfully"
// @Router       /teams/edit/{uuid} [put]
func (handler *TeamHandler) UpdateTeamByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	var team model.TeamUpdate
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamID, err := handler.service.GetIdByUuid(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := handler.service.UpdateTeamByID(teamID, team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team updated successfully"})
}

// UpdateTeamUserManagerStatus updates a team member's manager status.
//
// @Summary      Update team member manager status
// @Description  Updates a team member's manager status by team UUID and user UUID. 🔒 Requires role: **admin**
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Param        team_uuid  path      string  true  "Team UUID"
// @Param        user_uuid  path      string  true  "User UUID"
// @Param        is_manager path      bool    true  "Is Manager Status (1 = true, 0 = false)"
// @Success      200 "Status updated successfully"
// @Router       /teams/{team_uuid}/users/{user_uuid}/edit-manager-status [put]
func (handler *TeamHandler) UpdateTeamUserManagerStatus(c *gin.Context) {
	teamUUID := c.Param("team_uuid")
	userUUID := c.Param("user_uuid")
	isManager := c.Param("is_manager")

	// convert isManager to bool
	var req struct {
		IsManager bool `json:"is_manager"`
	}
	if isManager == "1" || isManager == "true" {
		req.IsManager = true
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := handler.service.UpdateTeamUserManagerStatus(teamUUID, userUUID, req.IsManager); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}
