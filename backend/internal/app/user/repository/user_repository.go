package repository

import (
	"app/internal/app/user/model"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]model.UserRead, error)
	FindByTypeAuth(typeOf string, data string) (*model.UserReadJWT, error)
	RegisterUser(user model.UserCreate) error
	FindIdByUuid(id string) (userId int, err error)
	UpdateUserStatus(userUUID string, status string) error
	UpdateUserLayout(userUUID string, layout model.UserDashboardLayoutUpdate) error
	DeleteUser(userUUID string) error
	DeleteUserLayout(userUUID string) error
	UpdateUser(userID int, user model.UserUpdateEntry) error
	FindByUUID(userUUID string) (*model.UserReadAll, error)
	FindDashboardLayoutByUUID(userUUID string) (*model.UserDashboardLayout, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (repo *userRepository) FindAll() ([]model.UserRead, error) {
	var users []model.UserRead
	err := repo.db.Raw(`
		SELECT
			u.uuid,
			u.username,
			u.email,
			u.first_name,
			u.last_name,
			u.phone_number,
			u.roles,
			u.status,
			u.first_day_of_week,
			wr.rate_name AS weekly_rate_name,
			COALESCE(wr.amount, 0) AS weekly_rate,
			u.created_at,
			u.updated_at,
			COALESCE(ws.status, 'completed') AS work_session_status
		FROM users u
		LEFT JOIN LATERAL (
			SELECT
				wsa.status
			FROM work_session_active wsa
			WHERE wsa.user_id = u.id
			ORDER BY wsa.created_at DESC
			LIMIT 1
		) ws ON TRUE
		LEFT JOIN weekly_rate wr ON wr.id = u.weekly_rate_id
		ORDER BY u.created_at DESC
	`).Scan(&users).Error
	return users, err
}

func (repo *userRepository) FindIdByUuid(uuid string) (int, error) {
	var userId int
	err := repo.db.Raw("SELECT id FROM users WHERE uuid = ?", uuid).Scan(&userId).Error
	if err != nil {
		return 0, err
	}
	if userId == 0 {
		return 0, fmt.Errorf("record not found")
	}
	return userId, nil
}

func (repo *userRepository) FindByTypeAuth(typeOf string, data string) (*model.UserReadJWT, error) {
	var user model.UserReadJWT

	if typeOf != "email" && typeOf != "username" {
		return nil, fmt.Errorf("invalid type: %s", typeOf)
	}

	query := fmt.Sprintf("SELECT id, uuid, first_day_of_week, email, roles, first_name, last_name, username, phone_number, password_hash FROM users WHERE %s = ?", typeOf)

	err := repo.db.Raw(query, data).Scan(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) RegisterUser(user model.UserCreate) error {
	err := repo.db.Exec(
		"INSERT INTO users (uuid, first_name, last_name, email, username, phone_number, roles, password_hash, weekly_rate_id, first_day_of_week) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.UUID, user.FirstName, user.LastName, user.Email, user.Username, user.PhoneNumber, user.Roles, user.PasswordHash, user.WeeklyRateID, user.FirstDayOfWeek,
	).Error
	return err
}

func (repo *userRepository) DeleteUser(userUUID string) error {
	err := repo.db.Exec("DELETE FROM users WHERE uuid = ?", userUUID).Error
	return err
}

func (repo *userRepository) UpdateUserStatus(userUUID string, status string) error {
	err := repo.db.Exec("UPDATE users SET status = ? WHERE uuid = ?", status, userUUID).Error
	return err
}

func (repo *userRepository) UpdateUser(userID int, user model.UserUpdateEntry) error {
	updateData := make(map[string]any)

	if user.Username != nil {
		updateData["username"] = *user.Username
	}
	if user.Email != nil {
		updateData["email"] = *user.Email
	}
	if user.FirstName != nil {
		updateData["first_name"] = *user.FirstName
	}
	if user.LastName != nil {
		updateData["last_name"] = *user.LastName
	}
	if user.PhoneNumber != nil {
		updateData["phone_number"] = *user.PhoneNumber
	}
	if user.Roles != nil {
		updateData["roles"] = *user.Roles
	}
	if user.Status != nil {
		updateData["status"] = *user.Status
	}

	if user.WeeklyRateID != nil {
		updateData["weekly_rate_id"] = *user.WeeklyRateID
	}

	if user.FirstDayOfWeek != nil {
		updateData["first_day_of_week"] = *user.FirstDayOfWeek
	}

	if len(updateData) == 0 {
		return fmt.Errorf("no fields to update")
	}

	err := repo.db.Table("users").Where("id = ?", userID).Updates(updateData).Error
	return err
}

func (repo *userRepository) FindByUUID(userUUID string) (*model.UserReadAll, error) {
	var user model.UserReadAll

	// ✅ Si c'est SQLite (tests unitaires)
	if repo.db.Dialector.Name() == "sqlite" {
		err := repo.db.Raw(`
			SELECT 
				u.uuid,
				u.username,
				u.email,
				u.first_name,
				u.last_name,
				u.phone_number,
				u.first_day_of_week,
				u.roles,
				u.status,
				COALESCE(wr.amount, 0) AS weekly_rate,
				wr.rate_name AS weekly_rate_name
			FROM users u
			LEFT JOIN weekly_rate wr ON wr.id = u.weekly_rate_id
			WHERE u.uuid = ?
			LIMIT 1
		`, userUUID).Scan(&user).Error

		if err != nil {
			return nil, err
		}

		// Pas de JSON_AGG ici, donc on simule un tableau vide
		user.Teams = []model.UserTeamMemberInfo{}
		return &user, nil
	}

	// ✅ Sinon : version PostgreSQL (production)
	err := repo.db.Raw(`
		SELECT 
			u.uuid,
			u.username,
			u.email,
			u.first_name,
			u.last_name,
			u.phone_number,
			u.first_day_of_week,
			u.roles,
			u.status,
			COALESCE(wr.amount, 0) AS weekly_rate,
			wr.rate_name AS weekly_rate_name,
			COALESCE(ws.status, 'completed') AS work_session_status,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'team_name', t.name,
						'team_description', t.description,
						'team_uuid', t.uuid,
						'is_manager', tm.is_manager
					)
				) FILTER (WHERE t.uuid IS NOT NULL),
				'[]'
			) AS teams
		FROM users u
		LEFT JOIN teams_members tm ON tm.user_id = u.id
		LEFT JOIN teams t ON t.id = tm.team_id
		LEFT JOIN weekly_rate wr ON wr.id = u.weekly_rate_id
		LEFT JOIN LATERAL (
			SELECT 
				wsa.status
			FROM work_session_active wsa
			WHERE wsa.user_id = u.id
			ORDER BY wsa.created_at DESC
			LIMIT 1
		) ws ON TRUE
		WHERE u.uuid = ?
		GROUP BY u.id, ws.status, wr.amount, wr.rate_name
	`, userUUID).Scan(&user).Error

	if err != nil {
		return nil, err
	}

	if user.TeamsRaw != "" {
		if err := json.Unmarshal([]byte(user.TeamsRaw), &user.Teams); err != nil {
			return nil, fmt.Errorf("failed to unmarshal teams JSON: %w", err)
		}
	} else {
		user.Teams = []model.UserTeamMemberInfo{}
	}

	return &user, nil
}

func (repo *userRepository) FindDashboardLayoutByUUID(userUUID string) (*model.UserDashboardLayout, error) {
	var raw string

	err := repo.db.Raw(`
		SELECT COALESCE(dashboard_layout, '[]')
		FROM users
		WHERE uuid = ?
	`, userUUID).Scan(&raw).Error
	if err != nil {
		return nil, err
	}

	var layout model.JSONLayout
	if err := json.Unmarshal([]byte(raw), &layout); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dashboard layout: %w", err)
	}

	if raw == "" {
		raw = "[]"
	}
	if err := json.Unmarshal([]byte(raw), &layout); err != nil {
		layout = model.JSONLayout{}
	}

	return &model.UserDashboardLayout{
		DashboardLayout: layout,
	}, nil
}

func (repo *userRepository) DeleteUserLayout(userUUID string) error {
	err := repo.db.Exec("UPDATE users SET dashboard_layout = NULL WHERE uuid = ?", userUUID).Error
	return err
}

func (repo *userRepository) UpdateUserLayout(userUUID string, layout model.UserDashboardLayoutUpdate) error {
	layoutJSON, err := json.Marshal(layout.Layout)
	if err != nil {
		return fmt.Errorf("failed to marshal layout to JSON: %w", err)
	}

	err = repo.db.Exec("UPDATE users SET dashboard_layout = ? WHERE uuid = ?", layoutJSON, userUUID).Error
	return err
}
