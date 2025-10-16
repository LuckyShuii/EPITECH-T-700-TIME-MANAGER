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
	GetUserActiveWorkSession(user_id int, status string) (workSession WorkSessionModel.WorkSessionRead, err error)
}

type workSessionRepository struct {
	db *gorm.DB
}

func NewWorkSessionRepository(db *gorm.DB) WorkSessionRepository {
	return &workSessionRepository{db}
}

func (repo *workSessionRepository) GetUserActiveWorkSession(userId int, status string) (workSession WorkSessionModel.WorkSessionRead, err error) {
	var workSessionFound WorkSessionModel.WorkSessionRead
	err = repo.db.Raw(
		"SELECT w.uuid as work_session_uuid, w.clock_in, w.clock_out, w.status, u.uuid as user_uuid, u.username, u.first_name, u.last_name, u.email, u.phone_number FROM work_session_active as w INNER JOIN users as u ON u.id = ? WHERE w.user_id = ? AND w.status = ? ORDER BY w.clock_in DESC LIMIT 1", userId, userId, status,
	).Scan(&workSessionFound).Error

	if err != nil {
		return WorkSessionModel.WorkSessionRead{}, fmt.Errorf("work session not found")
	}

	return workSessionFound, nil
}

func (repo *workSessionRepository) CompleteWorkSession(uuid string, user_id int, duration int) (err error) {
	log.Println("Completing work session with UUID: ", uuid, " for user ID: ", user_id)
	err = repo.db.Exec(
		"UPDATE work_session_active SET clock_out = now(), status = 'completed', duration_minutes = ? WHERE uuid = ? AND user_id = ?",
		duration, uuid, user_id,
	).Error
	return err
}

func (repo *workSessionRepository) CreateWorkSession(uuid string, user_id int, status string) error {
	err := repo.db.Exec(
		"INSERT INTO work_session_active (uuid, user_id, clock_in, status) VALUES (?, ?, ?, ?)",
		uuid, user_id, "now()", status,
	).Error
	return err
}
