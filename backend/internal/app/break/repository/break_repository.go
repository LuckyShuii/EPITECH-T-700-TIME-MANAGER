package repository

import (
	"fmt"

	"gorm.io/gorm"

	BreakModel "app/internal/app/break/model"
)

type BreakRepository interface {
	CompleteBreak(uuid string, workSessionId int, duration int) error
	CreateBreak(uuid string, workSessionId int, status string) error
	GetWorkSessionBreak(work_session_id int, status string) (BreakModel.BreakRead, error)
	GetTotalBreakDurationByWorkSessionId(workSessionId int) (int, error)
	DeleteRelatedBreaksToWorkSession(workSessionId int) error
}

type breakRepository struct {
	db *gorm.DB
}

func NewBreakRepository(db *gorm.DB) BreakRepository {
	return &breakRepository{db}
}

func (repo *breakRepository) GetWorkSessionBreak(workSessionId int, status string) (BreakModel.BreakRead, error) {
	var breakSessionFound BreakModel.BreakRead
	result := repo.db.Raw(`
		SELECT uuid as break_uuid, start_time, end_time, status
		FROM breaks
		WHERE work_session_active_id = ? AND status = ?
		ORDER BY start_time DESC
		LIMIT 1
	`, workSessionId, status).Scan(&breakSessionFound)

	if result.Error != nil {
		return BreakModel.BreakRead{}, result.Error
	}

	if breakSessionFound.BreakUUID == "" {
		return BreakModel.BreakRead{}, fmt.Errorf("break session not found")
	}

	return breakSessionFound, nil
}

func (repo *breakRepository) CompleteBreak(uuid string, workSessionId int, duration int) error {
	result := repo.db.Exec(`
		UPDATE breaks
		SET end_time = CURRENT_TIMESTAMP,
		    status = 'completed',
		    duration_minutes = ?
		WHERE uuid = ? AND work_session_active_id = ?
	`, duration, uuid, workSessionId)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no break found to update")
	}

	return nil
}

func (repo *breakRepository) CreateBreak(uuid string, workSessionId int, status string) error {
	return repo.db.Exec(`
		INSERT INTO breaks (uuid, work_session_active_id, start_time, status)
		VALUES (?, ?, CURRENT_TIMESTAMP, ?)
	`, uuid, workSessionId, status).Error
}

func (repo *breakRepository) GetTotalBreakDurationByWorkSessionId(workSessionId int) (int, error) {
	var totalDuration int
	err := repo.db.Raw(`
		SELECT COALESCE(SUM(duration_minutes), 0)
		FROM breaks
		WHERE work_session_active_id = ?
	`, workSessionId).Scan(&totalDuration).Error
	return totalDuration, err
}

func (repo *breakRepository) DeleteRelatedBreaksToWorkSession(workSessionId int) error {
	return repo.db.Exec(`
		DELETE FROM breaks
		WHERE work_session_active_id = ?
	`, workSessionId).Error
}
