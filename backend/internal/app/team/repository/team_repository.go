package repository

import (
	"app/internal/app/team/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository interface {
	FindAll() ([]model.TeamReadAll, error)
	FindIdByUuid(id string) (teamId int, err error)
	FindByID(id int) (model.TeamReadAll, error)
	DeleteByID(id int) error
	DeleteUserFromTeam(teamID int, userID int) error
	CreateTeam(teamUUID string, name string, description *string) error
	AddMembersToTeam(teamID int, members []model.TeamMemberCreate) error
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db}
}

func (repo *teamRepository) FindAll() ([]model.TeamReadAll, error) {
	var teams []model.TeamReadAll
	err := repo.db.Raw(`
		SELECT 
			t.uuid,
			t.name,
			t.description,
			json_agg(
				json_build_object(
					'user_uuid', u.uuid,
					'roles', u.roles,
					'status', u.status,
					'is_manager', tm.is_manager
				)
			) AS team_members
		FROM teams t
		JOIN teams_members tm ON tm.team_id = t.id
		JOIN users u ON u.id = tm.user_id
		LEFT JOIN work_session_active ws ON ws.user_id = u.id
		GROUP BY t.id, t.uuid, t.name, t.description
		ORDER BY t.name;
	`).Scan(&teams).Error
	return teams, err
}

func (repo *teamRepository) FindIdByUuid(uuid string) (teamId int, err error) {
	err = repo.db.Raw("SELECT id FROM teams WHERE uuid = ?", uuid).Scan(&teamId).Error
	if err != nil {
		return 0, err
	}
	return teamId, nil
}

func (repo *teamRepository) FindByID(id int) (model.TeamReadAll, error) {
	var team model.TeamReadAll
	err := repo.db.Raw(`
		SELECT 
			t.uuid,
			t.name,
			t.description,
			json_agg(
				json_build_object(
					'user_uuid', u.uuid,
					'roles', u.roles,
					'status', u.status,
					'is_manager', tm.is_manager
				)
			) AS team_members
		FROM teams t
		JOIN teams_members tm ON tm.team_id = t.id
		JOIN users u ON u.id = tm.user_id
		LEFT JOIN work_session_active ws ON ws.user_id = u.id
		WHERE t.id = ?
		GROUP BY t.id, t.uuid, t.name, t.description;
	`, id).Scan(&team).Error
	return team, err
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
