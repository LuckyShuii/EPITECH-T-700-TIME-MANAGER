package repository

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	WorkSessionModel "app/internal/app/work-session/model"
)

type WorkSessionRepository interface {
	CompleteWorkSession(uuid string, user_id int, duration int) (err error)
	CreateWorkSession(uuid string, user_id int, status string) error
	GetUserActiveWorkSession(user_id int, status []string) (workSession WorkSessionModel.WorkSessionRead, err error)
	FindIdByUuid(uuid string) (workSessionId int, err error)
	UpdateWorkSessionStatus(uuid string, status string) error
	UpdateBreakDurationMinutes(uuid string, breakDuration int) error
	GetWorkSessionHistoryByUserId(userId int, startDate string, endDate string, limit int, offset int) (workSessions []WorkSessionModel.WorkSessionReadHistory, err error)
}

type workSessionRepository struct {
	db *gorm.DB
}

func NewWorkSessionRepository(db *gorm.DB) WorkSessionRepository {
	return &workSessionRepository{db}
}

func (repo *workSessionRepository) FindIdByUuid(uuid string) (workSessionId int, err error) {
	err = repo.db.Raw("SELECT id FROM work_session_active WHERE uuid = ?", uuid).Scan(&workSessionId).Error
	if err != nil {
		return 0, err
	}
	return workSessionId, nil
}

func (repo *workSessionRepository) GetUserActiveWorkSession(userId int, status []string) (workSession WorkSessionModel.WorkSessionRead, err error) {
	var workSessionFound WorkSessionModel.WorkSessionRead
	err = repo.db.Raw(
		`SELECT 
		w.uuid as work_session_uuid, 
		w.clock_in, 
		w.clock_out, 
		w.status, 
		u.uuid as user_uuid, 
		u.username,
		u.first_name, 
		u.last_name, 
		u.email, 
		u.phone_number 
		FROM work_session_active as w 
		INNER JOIN users as u ON u.id = ? 
		WHERE w.user_id = ? 
		AND w.status IN (?) 
		ORDER BY w.clock_in 
		DESC LIMIT 1`,
		userId, userId, status,
	).Scan(&workSessionFound).Error

	if err != nil {
		return WorkSessionModel.WorkSessionRead{}, fmt.Errorf("work session not found")
	}

	return workSessionFound, nil
}

func (repo *workSessionRepository) CompleteWorkSession(uuid string, userId int, duration int) (err error) {
	log.Println("Completing work session with UUID: ", uuid, " for user ID: ", userId)
	err = repo.db.Exec(
		"UPDATE work_session_active SET clock_out = CURRENT_TIMESTAMP, status = 'completed', duration_minutes = ? WHERE uuid = ? AND user_id = ?",
		duration, uuid, userId,
	).Error
	return err
}

func (repo *workSessionRepository) CreateWorkSession(uuid string, userId int, status string) error {
	err := repo.db.Exec(
		"INSERT INTO work_session_active (uuid, user_id, clock_in, status) VALUES (?, ?, ?, ?)",
		uuid, userId, "now()", status,
	).Error
	return err
}

func (repo *workSessionRepository) UpdateWorkSessionStatus(uuid string, status string) error {
	err := repo.db.Exec(
		"UPDATE work_session_active SET status = ? WHERE uuid = ?",
		status, uuid,
	).Error
	return err
}

func (repo *workSessionRepository) UpdateBreakDurationMinutes(uuid string, breakDuration int) error {
	err := repo.db.Exec(
		"UPDATE work_session_active SET breaks_duration_minutes = ? WHERE uuid = ?",
		breakDuration, uuid,
	).Error
	return err
}

func (repo *workSessionRepository) GetWorkSessionHistoryByUserId(userId int, startDate string, endDate string, limit int, offset int) (workSessions []WorkSessionModel.WorkSessionReadHistory, err error) {
	err = repo.db.Raw(
		`SELECT
		ws.uuid AS work_session_uuid,
		u.uuid AS user_uuid,
		u.username AS username,
		ws.clock_in,
		ws.clock_out,
		ws.duration_minutes,
		ws.breaks_duration_minutes,
		ws.status
		FROM users AS u
		INNER JOIN (
			SELECT
				user_id,
				clock_in,
				clock_out,
				duration_minutes,
				breaks_duration_minutes,
				status,
				uuid
			FROM work_session_active
			UNION ALL
			SELECT
				user_id,
				clock_in,
				clock_out,
				duration_minutes,
				breaks_duration_minutes,
				status,
				uuid
			FROM work_session_archived
		) AS ws ON u.id = ws.user_id
		WHERE
			u.id = ?
			AND ws.clock_in BETWEEN ? AND ?
		ORDER BY ws.clock_in DESC
		LIMIT ? OFFSET ?`,
		userId, startDate, endDate, limit, offset,
	).Scan(&workSessions).Error

	if err != nil {
		return nil, err
	}

	return workSessions, nil
}
