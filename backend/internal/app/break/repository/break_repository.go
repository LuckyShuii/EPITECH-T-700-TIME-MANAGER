package repository

import (
	"fmt"

	"gorm.io/gorm"

	BreakModel "app/internal/app/break/model"
)

type BreakRepository interface {
	CompleteBreak(uuid string, workSessionId int, duration int) (err error)
	CreateBreak(uuid string, workSessionId int, status string) error
	GetWorkSessionBreak(work_session_id int, status string) (breakSession BreakModel.BreakRead, err error)
	GetTotalBreakDurationByWorkSessionId(workSessionId int) (totalDuration int, err error)
	DeleteRelatedBreaksToWorkSession(workSessionId int) error
}

type breakRepository struct {
	db *gorm.DB
}

func NewBreakRepository(db *gorm.DB) BreakRepository {
	return &breakRepository{db}
}

func (repo *breakRepository) GetWorkSessionBreak(workSessionId int, status string) (breakSession BreakModel.BreakRead, err error) {
	var breakSessionFound BreakModel.BreakRead
	err = repo.db.Raw(
		"SELECT uuid as break_uuid, start_time, end_time, status FROM breaks WHERE work_session_active_id = ? AND status = ? ORDER BY start_time DESC LIMIT 1", workSessionId, status,
	).Scan(&breakSessionFound).Error

	if err != nil {
		return BreakModel.BreakRead{}, fmt.Errorf("break session not found")
	}

	return breakSessionFound, nil
}

func (repo *breakRepository) CompleteBreak(uuid string, workSessionId int, duration int) (err error) {
	err = repo.db.Exec(
		"UPDATE breaks SET end_time = now(), status = 'completed', duration_minutes = ? WHERE uuid = ? AND work_session_active_id = ?",
		duration, uuid, workSessionId,
	).Error
	return err
}

func (repo *breakRepository) CreateBreak(uuid string, workSessionId int, status string) error {
	err := repo.db.Exec(
		"INSERT INTO breaks (uuid, work_session_active_id, start_time, status) VALUES (?, ?, ?, ?)",
		uuid, workSessionId, "now()", status,
	).Error
	return err
}

func (repo *breakRepository) GetTotalBreakDurationByWorkSessionId(workSessionId int) (totalDuration int, err error) {
	err = repo.db.Raw("SELECT COALESCE(SUM(duration_minutes), 0) FROM breaks WHERE work_session_active_id = ?", workSessionId).Scan(&totalDuration).Error
	if err != nil {
		return 0, err
	}
	return totalDuration, nil
}

func (repo *breakRepository) DeleteRelatedBreaksToWorkSession(workSessionId int) error {
	err := repo.db.Exec(
		"DELETE FROM breaks WHERE work_session_active_id = ?",
		workSessionId,
	).Error
	return err
}
