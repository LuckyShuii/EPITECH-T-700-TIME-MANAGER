package repository

import (
	"app/internal/app/team/model"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository interface {
	FindAll() ([]model.TeamReadAll, error)
	FindIdByUuid(id string) (teamId int, err error)
	FindUserIDsByTeamID(teamID int) ([]int, error)
	FindByID(id int) (model.TeamReadAll, error)
	DeleteByID(id int) error
	DeleteUserFromTeam(teamID int, userID int) error
	CreateTeam(teamUUID string, name string, description *string) error
	AddMembersToTeam(teamID int, members []model.TeamMemberCreate) error
	UpdateTeamByID(id int, updatedTeam model.TeamUpdate) error
	UpdateTeamUserManagerStatus(teamID int, userID int, isManager bool) error
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db}
}

func (repo *teamRepository) FindAll() ([]model.TeamReadAll, error) {
	var teams []model.TeamReadAll

	query := `
		SELECT 
			t.uuid,
			t.name,
			t.description,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'user_uuid', u.uuid,
					'roles', u.roles,
					'status', u.status,
					'work_session_status', COALESCE(ws.status, 'completed'),
					'username', u.username,
					'email', u.email,
					'first_name', u.first_name,
					'last_name', u.last_name,
					'phone_number', u.phone_number,
					'first_day_of_week', u.first_day_of_week,
					'is_manager', tm.is_manager,
					'weekly_rate', COALESCE(wr.amount, 0),
					'weekly_rate_name', wr.rate_name
				)
			) AS team_members
		FROM teams t
		LEFT JOIN teams_members tm ON tm.team_id = t.id
		LEFT JOIN users u ON u.id = tm.user_id
		LEFT JOIN weekly_rate wr ON wr.id = u.weekly_rate_id
		LEFT JOIN LATERAL (
			SELECT 
				wsa.status
			FROM work_session_active wsa
			WHERE wsa.user_id = u.id
			ORDER BY wsa.created_at DESC
			LIMIT 1
		) ws ON TRUE
		GROUP BY 
			t.id, t.uuid, t.name, t.description
		ORDER BY 
			t.name;
	`

	err := repo.db.Raw(query).Scan(&teams).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams: %w", err)
	}

	return teams, nil
}

func (repo *teamRepository) FindIdByUuid(uuid string) (teamId int, err error) {
	err = repo.db.Raw("SELECT id FROM teams WHERE uuid = ?", uuid).Scan(&teamId).Error
	if err != nil {
		return 0, err
	}
	if teamId == 0 {
		return 0, fmt.Errorf("team not found")
	}
	return teamId, nil
}

func (repo *teamRepository) FindByID(id int) (model.TeamReadAll, error) {
	var team model.TeamReadAll

	query := `
		SELECT 
			t.uuid,
			t.name,
			t.description,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'user_uuid', u.uuid,
					'roles', u.roles,
					'status', u.status,
					'work_session_status', COALESCE(ws.status, 'completed'),
					'is_manager', tm.is_manager,
					'username', u.username,
					'email', u.email,
					'first_name', u.first_name,
					'last_name', u.last_name,
					'first_day_of_week', u.first_day_of_week,
					'phone_number', u.phone_number,
					'weekly_rate', COALESCE(wr.amount, 0),
					'weekly_rate_name', wr.rate_name
				)
			) AS team_members
		FROM teams t
		LEFT JOIN teams_members tm ON tm.team_id = t.id
		LEFT JOIN users u ON u.id = tm.user_id
		LEFT JOIN weekly_rate wr ON wr.id = u.weekly_rate_id
		LEFT JOIN LATERAL (
			SELECT 
				wsa.status
			FROM work_session_active wsa
			WHERE wsa.user_id = u.id
			ORDER BY wsa.created_at DESC
			LIMIT 1
		) ws ON TRUE
		WHERE t.id = ?
		GROUP BY 
			t.id, t.uuid, t.name, t.description;
	`

	err := repo.db.Raw(query, id).Scan(&team).Error
	if err != nil {
		return team, fmt.Errorf("failed to fetch team: %w", err)
	}
	if team.UUID == "" {
		return model.TeamReadAll{}, fmt.Errorf("team with id %d not found", id)
	}
	return team, nil
}

func (repo *teamRepository) DeleteByID(id int) error {
	return repo.db.Exec("DELETE FROM teams WHERE id = ?", id).Error
}

func (repo *teamRepository) DeleteUserFromTeam(teamID int, userID int) error {
	return repo.db.Exec("DELETE FROM teams_members WHERE team_id = ? AND user_id = ?", teamID, userID).Error
}

func (repo *teamRepository) CreateTeam(teamUUID string, name string, description *string) error {
	var team model.TeamReadAll
	err := repo.db.Raw(
		"INSERT INTO teams (uuid, name, description) VALUES (?, ?, ?) RETURNING uuid, name, description",
		teamUUID, name, description,
	).Scan(&team).Error
	return err
}

func (repo *teamRepository) AddMembersToTeam(teamID int, members []model.TeamMemberCreate) error {
	if len(members) == 0 {
		return nil
	}

	query := "INSERT INTO teams_members (uuid, team_id, user_id, is_manager) VALUES "
	values := []interface{}{}

	for i, m := range members {
		uuid := uuid.New().String()
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?, ?)"
		values = append(values, uuid, teamID, m.UserID, m.IsManager)
	}

	return repo.db.Exec(query, values...).Error
}

func (repo *teamRepository) UpdateTeamByID(id int, updatedTeam model.TeamUpdate) error {
	updateData := make(map[string]any)

	if updatedTeam.Name != nil {
		updateData["name"] = *updatedTeam.Name
	}
	if updatedTeam.Description != nil {
		updateData["description"] = *updatedTeam.Description
	}

	if len(updateData) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := "UPDATE teams SET "
	args := []interface{}{}
	i := 0

	for field, value := range updateData {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("\"%s\" = ?", field)
		args = append(args, value)
		i++
	}

	query += " WHERE id = ?"
	args = append(args, id)

	return repo.db.Exec(query, args...).Error
}

func (repo *teamRepository) UpdateTeamUserManagerStatus(teamID int, userID int, isManager bool) error {
	return repo.db.Exec("UPDATE teams_members SET is_manager = ? WHERE team_id = ? AND user_id = ?", isManager, teamID, userID).Error
}

func (repo *teamRepository) FindUserIDsByTeamID(teamID int) ([]int, error) {
	var userIDs []int
	err := repo.db.Raw("SELECT user_id FROM teams_members WHERE team_id = ?", teamID).Scan(&userIDs).Error
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}
