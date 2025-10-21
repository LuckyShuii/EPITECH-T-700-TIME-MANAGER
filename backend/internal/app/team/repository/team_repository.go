package repository

import (
	"app/internal/app/team/model"

	"gorm.io/gorm"
)

type TeamRepository interface {
	FindAll() ([]model.TeamReadAll, error)
	FindIdByUuid(id string) (teamId int, err error)
	FindByID(id int) (model.TeamReadAll, error)
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
