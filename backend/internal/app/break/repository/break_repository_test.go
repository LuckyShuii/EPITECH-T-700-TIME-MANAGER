package repository_test

import (
	"app/internal/app/break/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB initializes a temporary SQLite in-memory database.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.Exec(`
		CREATE TABLE breaks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT,
			work_session_active_id INTEGER,
			start_time TEXT,
			end_time TEXT,
			status TEXT,
			duration_minutes INTEGER
		);
	`).Error
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	return db
}

func TestCreateBreak(t *testing.T) {
	const uuid = "uuid-1"
	db := setupTestDB(t)
	repo := repository.NewBreakRepository(db)

	err := repo.CreateBreak(uuid, 10, "active")
	assert.NoError(t, err)

	var count int64
	db.Table("breaks").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestCompleteBreak(t *testing.T) {
	const uuid = "uuid-1"
	db := setupTestDB(t)
	repo := repository.NewBreakRepository(db)

	// Insert row manually so it can be updated later
	db.Exec(`
		INSERT INTO breaks (uuid, work_session_active_id, start_time, status)
		VALUES (?, 10, ?, 'active')
	`, uuid, time.Now().Format(time.RFC3339))

	err := repo.CompleteBreak(uuid, 10, 15)
	assert.NoError(t, err)

	var req struct {
		Status          string
		DurationMinutes int
	}

	db.Raw("SELECT status, duration_minutes FROM breaks WHERE uuid = ?", uuid).Scan(&req)

	assert.Equal(t, "completed", req.Status)
	assert.Equal(t, 15, req.DurationMinutes)
}

func TestGetWorkSessionBreakFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBreakRepository(db)

	start := time.Now().Format(time.RFC3339)
	db.Exec(`
		INSERT INTO breaks (uuid, work_session_active_id, start_time, status)
		VALUES ('uuid-2', 99, ?, 'active')
	`, start)

	breakFound, err := repo.GetWorkSessionBreak(99, "active")

	assert.NoError(t, err)
	assert.Equal(t, "uuid-2", breakFound.BreakUUID)
	assert.Equal(t, "active", breakFound.Status)
}

func TestGetTotalBreakDurationByWorkSessionId(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBreakRepository(db)

	db.Exec(`
		INSERT INTO breaks (uuid, work_session_active_id, duration_minutes)
		VALUES ('b1', 20, 10), ('b2', 20, 30)
	`)

	total, err := repo.GetTotalBreakDurationByWorkSessionId(20)
	assert.NoError(t, err)
	assert.Equal(t, 40, total)
}

func TestDeleteRelatedBreaksToWorkSession(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBreakRepository(db)

	db.Exec(`
		INSERT INTO breaks (uuid, work_session_active_id, status)
		VALUES ('b1', 25, 'completed')
	`)

	err := repo.DeleteRelatedBreaksToWorkSession(25)
	assert.NoError(t, err)

	var count int64
	db.Table("breaks").Count(&count)
	assert.Equal(t, int64(0), count)
}
