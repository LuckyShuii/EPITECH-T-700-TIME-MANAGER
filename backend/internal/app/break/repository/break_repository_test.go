package repository_test

import (
	"app/internal/app/break/repository"
	"app/internal/test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const WORK_SESSION_NEW_ERROR = "Failed to add work session: %v"

func addWorkSession(t *testing.T, uuid string) (int, error) {
	db := test.ResetDB(t)

	err := db.Exec(`
		INSERT INTO users (uuid, username, email, password_hash, status)
		VALUES (?, ?, ?, ?, 'active')
	`, uuid, "testuser", "testuser@example.com", "hashedpassword").Error
	if err != nil {
		return 0, err
	}

	var userId int
	err = db.Raw("SELECT id FROM users WHERE uuid = ?", uuid).Scan(&userId).Error
	if err != nil {
		return 0, err
	}

	err = db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, status)
		VALUES (?, ?, ?, 'active')
	`, uuid, userId, time.Now().Format(time.RFC3339)).Error

	if err != nil {
		return 0, err
	}

	var workSessionID int
	err = db.Raw("SELECT id FROM work_session_active WHERE uuid = ?", uuid).Scan(&workSessionID).Error
	if err != nil {
		return 0, err
	}

	return workSessionID, nil
}

func TestCreateBreak(t *testing.T) {
	uuid := "123e4567-e89b-12d3-a456-426614174001"
	db := test.ResetDB(t)
	repo := repository.NewBreakRepository(db)

	workSessionID, err := addWorkSession(t, uuid)

	if err != nil {
		t.Fatalf(WORK_SESSION_NEW_ERROR, err)
	}

	err = repo.CreateBreak(uuid, workSessionID, "active")
	assert.NoError(t, err)

	var count int64
	db.Table("breaks").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestCompleteBreak(t *testing.T) {
	uuid := "123e4567-e89b-12d3-a456-426614174002"
	db := test.ResetDB(t)
	repo := repository.NewBreakRepository(db)

	workSessionID, err := addWorkSession(t, uuid)

	if err != nil {
		t.Fatalf(WORK_SESSION_NEW_ERROR, err)
	}

	// Insert row manually so it can be updated later
	db.Exec(`
		INSERT INTO breaks (uuid, work_session_active_id, start_time, status)
		VALUES (?, ?, ?, 'active')
	`, uuid, workSessionID, time.Now().Format(time.RFC3339))

	err = repo.CompleteBreak(uuid, workSessionID, 15)
	assert.NoError(t, err)

	var req struct {
		Status          string
		DurationMinutes int
	}

	db.Raw("SELECT status, duration_minutes FROM breaks WHERE uuid = ?", uuid).Scan(&req)

	assert.Equal(t, "completed", req.Status)
	assert.Equal(t, 15, req.DurationMinutes)
}

func TestCompleteBreakNotFound(t *testing.T) {
	const uuid = "123e4567-e8oi-12d3-a456-426614174002"
	db := test.ResetDB(t)
	repo := repository.NewBreakRepository(db)

	err := repo.CompleteBreak(uuid, 10, 15)
	assert.Error(t, err)
}

func TestGetWorkSessionBreakNotFound(t *testing.T) {
	db := test.ResetDB(t)
	repo := repository.NewBreakRepository(db)

	br, err := repo.GetWorkSessionBreak(999, "active")

	assert.NoError(t, err)
	assert.NotNil(t, br)
	assert.Empty(t, br.UUID)
	assert.Empty(t, br.Status)
}

func TestGetWorkSessionBreakFound(t *testing.T) {
	uuid := "123e4567-e89b-12d3-a456-426614174003"
	db := test.ResetDB(t)
	repo := repository.NewBreakRepository(db)

	workSessionID, err := addWorkSession(t, uuid)

	if err != nil {
		t.Fatalf(WORK_SESSION_NEW_ERROR, err)
	}

	start := time.Now().Format(time.RFC3339)
	db.Exec(`
		INSERT INTO breaks (uuid, work_session_active_id, start_time, status)
		VALUES (?, ?, ?, 'active')
	`, uuid, workSessionID, start)

	breakFound, err := repo.GetWorkSessionBreak(workSessionID, "active")

	assert.NoError(t, err)
	assert.Equal(t, uuid, breakFound.BreakUUID)
	assert.Equal(t, "active", breakFound.Status)
}

func TestGetTotalBreakDurationByWorkSessionId(t *testing.T) {
	uuid1 := "123e4567-e89b-12d3-a456-426614174004"
	uuid2 := "223e4567-e89b-12d3-a456-426614174006"
	db := test.ResetDB(t)
	repo := repository.NewBreakRepository(db)

	workSessionID, err := addWorkSession(t, uuid1)

	if err != nil {
		t.Fatalf(WORK_SESSION_NEW_ERROR, err)
	}

	db.Exec(`
		INSERT INTO breaks (uuid, work_session_active_id, duration_minutes, start_time)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?)
	`, uuid1, workSessionID, 10, time.Now().Format(time.RFC3339), uuid2, workSessionID, 30, time.Now().Format(time.RFC3339),
	)

	total, err := repo.GetTotalBreakDurationByWorkSessionId(workSessionID)
	assert.NoError(t, err)
	assert.Equal(t, 40, total)
}

func TestDeleteRelatedBreaksToWorkSession(t *testing.T) {
	uuid := "123e4567-e89b-12d3-a456-426614174005"
	db := test.ResetDB(t)
	repo := repository.NewBreakRepository(db)

	workSessionID, err := addWorkSession(t, uuid)

	if err != nil {
		t.Fatalf(WORK_SESSION_NEW_ERROR, err)
	}

	// Insert a break related to the work session
	db.Exec(`
		INSERT INTO breaks (uuid, work_session_active_id, status, start_time)
		VALUES (?, ?, 'completed', ?)
	`, uuid, workSessionID, time.Now().Format(time.RFC3339))

	err = repo.DeleteRelatedBreaksToWorkSession(workSessionID)
	assert.NoError(t, err)

	var count int64
	db.Table("breaks").Count(&count)
	assert.Equal(t, int64(0), count)
}
