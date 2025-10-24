package repository_test

import (
	"app/internal/app/work-session/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite schema for the work_session repository.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT,
			username TEXT,
			first_name TEXT,
			last_name TEXT,
			email TEXT,
			phone_number TEXT
		);

		CREATE TABLE work_session_active (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT,
			user_id INTEGER,
			clock_in TEXT,
			clock_out TEXT,
			status TEXT,
			duration_minutes INTEGER,
			breaks_duration_minutes INTEGER
		);

		CREATE TABLE work_session_archived (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT,
			user_id INTEGER,
			clock_in TEXT,
			clock_out TEXT,
			status TEXT,
			duration_minutes INTEGER,
			breaks_duration_minutes INTEGER
		);
	`).Error
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestCreateWorkSession(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWorkSessionRepository(db)

	uuid := "ws-uuid-1"
	err := repo.CreateWorkSession(uuid, 1, "active")
	assert.NoError(t, err)

	var count int64
	db.Table("work_session_active").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestFindIdByUuid(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWorkSessionRepository(db)

	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, status)
		VALUES ('ws-uuid-2', 1, CURRENT_TIMESTAMP, 'active')
	`)

	id, err := repo.FindIdByUuid("ws-uuid-2")
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestUpdateWorkSessionStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWorkSessionRepository(db)

	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, status)
		VALUES ('ws-uuid-3', 1, CURRENT_TIMESTAMP, 'active')
	`)

	err := repo.UpdateWorkSessionStatus("ws-uuid-3", "paused")
	assert.NoError(t, err)

	var status string
	db.Raw(`SELECT status FROM work_session_active WHERE uuid = 'ws-uuid-3'`).Scan(&status)
	assert.Equal(t, "paused", status)
}

func TestUpdateBreakDurationMinutes(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWorkSessionRepository(db)

	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, status)
		VALUES ('ws-uuid-4', 1, CURRENT_TIMESTAMP, 'active')
	`)

	err := repo.UpdateBreakDurationMinutes("ws-uuid-4", 25)
	assert.NoError(t, err)

	var breaks int
	db.Raw(`SELECT breaks_duration_minutes FROM work_session_active WHERE uuid = 'ws-uuid-4'`).Scan(&breaks)
	assert.Equal(t, 25, breaks)
}

func TestCompleteWorkSession(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWorkSessionRepository(db)

	uuid := "ws-uuid-5"
	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, status)
		VALUES (?, 1, CURRENT_TIMESTAMP, 'active')
	`, uuid)

	err := repo.CompleteWorkSession(uuid, 1, 180)
	assert.NoError(t, err)

	var row struct {
		Status          string
		DurationMinutes int
	}
	db.Raw(`SELECT status, duration_minutes FROM work_session_active WHERE uuid = ?`, uuid).Scan(&row)

	assert.Equal(t, "completed", row.Status)
	assert.Equal(t, 180, row.DurationMinutes)
}

func TestGetUserActiveWorkSession(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWorkSessionRepository(db)

	// Create user and work session
	db.Exec(`
		INSERT INTO users (id, uuid, username, first_name, last_name, email)
		VALUES (1, 'user-uuid', 'john', 'John', 'Doe', 'john@example.com')
	`)
	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, status)
		VALUES ('ws-uuid-6', 1, CURRENT_TIMESTAMP, 'active')
	`)

	session, err := repo.GetUserActiveWorkSession(1, []string{"active"})
	assert.NoError(t, err)
	assert.Equal(t, "ws-uuid-6", session.WorkSessionUUID)
	assert.Equal(t, "active", session.Status)
	assert.Equal(t, "user-uuid", session.UserUUID)
}

func TestGetWorkSessionHistoryByUserId(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWorkSessionRepository(db)

	// Create user
	db.Exec(`INSERT INTO users (id, uuid, username) VALUES (1, 'user-uuid', 'john')`)

	// Insert sessions in active + archived tables
	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, clock_out, status, duration_minutes)
		VALUES ('active-uuid', 1, '2025-10-01T10:00:00Z', '2025-10-01T12:00:00Z', 'completed', 120)
	`)
	db.Exec(`
		INSERT INTO work_session_archived (uuid, user_id, clock_in, clock_out, status, duration_minutes)
		VALUES ('archived-uuid', 1, '2025-09-01T10:00:00Z', '2025-09-01T11:00:00Z', 'completed', 60)
	`)

	start := "2025-09-01T00:00:00Z"
	end := "2025-11-01T00:00:00Z"

	results, err := repo.GetWorkSessionHistoryByUserId(1, start, end, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "active-uuid", results[0].WorkSessionUUID)
}
