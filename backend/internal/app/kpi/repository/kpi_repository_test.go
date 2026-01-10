package repository_test

import (
	"app/internal/app/kpi/repository"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupKPITestDB starts a Postgres testcontainer and creates a minimal schema
// compatible with KPI repository queries. No external SQL files needed.
func setupKPITestDB(t *testing.T) *gorm.DB {
	t.Helper()

	ctx := context.Background()
	const portable = "5432/tcp"

	req := testcontainers.ContainerRequest{
		Image: "postgres:16",
		Env: map[string]string{
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "postgres",
		},
		ExposedPorts: []string{portable},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(portable),
			wait.ForLog("database system is ready to accept connections").
				WithStartupTimeout(60*time.Second).
				WithOccurrence(2),
		),
	}

	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	t.Cleanup(func() {
		_ = pgC.Terminate(ctx)
	})

	port, err := pgC.MappedPort(ctx, portable)
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
	}

	host, err := pgC.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=postgres password=test dbname=testdb sslmode=disable connect_timeout=10",
		host, port.Port(),
	)

	// Connect (retry a bit)
	var db *gorm.DB
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			// keep tests quiet; remove this if you want to see every SQL query
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			sqlDB, e := db.DB()
			if e == nil {
				if e = sqlDB.Ping(); e == nil {
					break
				}
			}
		}
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	if err != nil {
		t.Fatalf("Failed to connect to Postgres after %d attempts: %v", maxRetries, err)
	}

	schema := `
		CREATE EXTENSION IF NOT EXISTS pgcrypto;

		CREATE TYPE work_session_status AS ENUM('active', 'completed', 'paused');

		CREATE TABLE weekly_rate (
			id SERIAL PRIMARY KEY,
			uuid VARCHAR(36) NOT NULL UNIQUE,
			rate_name VARCHAR(255) NOT NULL,
			amount SMALLINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE users (
			id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			uuid VARCHAR(36) NOT NULL UNIQUE,
			username VARCHAR(100) NOT NULL UNIQUE,
			email VARCHAR(320) NOT NULL UNIQUE,
			password_hash VARCHAR(100) NOT NULL,
			status VARCHAR(15) NOT NULL DEFAULT 'pending',
			roles TEXT[] DEFAULT '{"employee"}',
			weekly_rate_id INT NULL REFERENCES weekly_rate (id),
			dashboard_layout JSON DEFAULT NULL,
			first_day_of_week INT DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE work_session_active (
			id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			uuid VARCHAR(36) NOT NULL UNIQUE,
			user_id INT NOT NULL REFERENCES users(id),
			clock_in TIMESTAMP NOT NULL,
			clock_out TIMESTAMP,
			duration_minutes INT,
			status work_session_status DEFAULT 'active',
			breaks_duration_minutes INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE work_session_archived (
			id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			uuid VARCHAR(36) NOT NULL UNIQUE,
			user_id INT NOT NULL REFERENCES users(id),
			clock_in TIMESTAMP NOT NULL,
			clock_out TIMESTAMP,
			duration_minutes INT,
			status work_session_status DEFAULT 'active',
			breaks_duration_minutes INTEGER DEFAULT 0,
			archived_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	if err := db.Exec(schema).Error; err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func insertUser(t *testing.T, db *gorm.DB, uuid, username string, weeklyRateID *int) int {
	t.Helper()

	email := fmt.Sprintf("%s@example.com", username)

	var id int
	if weeklyRateID != nil {
		err := db.Raw(`
			INSERT INTO users (uuid, username, email, password_hash, status, weekly_rate_id)
			VALUES (?, ?, ?, 'hash', 'active', ?)
			RETURNING id
		`, uuid, username, email, *weeklyRateID).Scan(&id).Error
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}
		return id
	}

	err := db.Raw(`
		INSERT INTO users (uuid, username, email, password_hash, status)
		VALUES (?, ?, ?, 'hash', 'active')
		RETURNING id
	`, uuid, username, email).Scan(&id).Error
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}
	return id
}

func insertWeeklyRate(t *testing.T, db *gorm.DB, uuid, name string, amount int) int {
	t.Helper()

	var id int
	err := db.Raw(`
		INSERT INTO weekly_rate (uuid, rate_name, amount)
		VALUES (?, ?, ?)
		RETURNING id
	`, uuid, name, amount).Scan(&id).Error
	if err != nil {
		t.Fatalf("failed to insert weekly_rate: %v", err)
	}
	return id
}

func TestGetWeeklyRatesByUserIDAndDateRange(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	userID := insertUser(t, db, "user-uuid-1", "testuser", nil)

	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-a-1', ?, '2026-01-06 09:00:00', 480, 'completed')
	`, userID)
	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-a-2', ?, '2026-01-07 09:00:00', 450, 'completed')
	`, userID)
	db.Exec(`
		INSERT INTO work_session_archived (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-ar-1', ?, '2026-01-08 09:00:00', 500, 'completed')
	`, userID)

	totalMinutes, err := repo.GetWeeklyRatesByUserIDAndDateRange(userID, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 1430, totalMinutes)
}

func TestGetWeeklyRatesByUserIDAndDateRangeNoData(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	userID := insertUser(t, db, "user-uuid-2", "testuser2", nil)

	totalMinutes, err := repo.GetWeeklyRatesByUserIDAndDateRange(userID, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 0, totalMinutes)
}

func TestGetUserPresenceRate(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	wrID := insertWeeklyRate(t, db, "wr-1", "Standard", 40)
	userID := insertUser(t, db, "user-uuid-1", "testuser", &wrID)

	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-a-1', ?, '2026-01-06 09:00:00', 1200, 'completed')
	`, userID)
	db.Exec(`
		INSERT INTO work_session_archived (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-ar-1', ?, '2026-01-07 09:00:00', 1200, 'completed')
	`, userID)

	presenceRate, weeklyRateExpected, weeklyTimeDone, err := repo.GetUserPresenceRate(userID, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 100.0, presenceRate)
	assert.Equal(t, 40.0, weeklyRateExpected)
	assert.Equal(t, 40.0, weeklyTimeDone)
}

func TestGetUserPresenceRateWithoutWeeklyRate(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	userID := insertUser(t, db, "user-uuid-2", "testuser2", nil)

	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-a-1', ?, '2026-01-06 09:00:00', 1200, 'completed')
	`, userID)

	presenceRate, weeklyRateExpected, weeklyTimeDone, err := repo.GetUserPresenceRate(userID, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 50.0, presenceRate)
	assert.Equal(t, 40.0, weeklyRateExpected)
	assert.Equal(t, 20.0, weeklyTimeDone)
}

func TestGetUserAverageBreakTime(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	userID := insertUser(t, db, "user-uuid-1", "testuser", nil)

	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, breaks_duration_minutes, status)
		VALUES ('ws-a-1', ?, '2026-01-06 09:00:00', 480, 30, 'completed')
	`, userID)
	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, breaks_duration_minutes, status)
		VALUES ('ws-a-2', ?, '2026-01-07 09:00:00', 450, 45, 'completed')
	`, userID)
	db.Exec(`
		INSERT INTO work_session_archived (uuid, user_id, clock_in, duration_minutes, breaks_duration_minutes, status)
		VALUES ('ws-ar-1', ?, '2026-01-08 09:00:00', 500, 25, 'completed')
	`, userID)

	averageBreakTime, err := repo.GetUserAverageBreakTime(userID, "2026-01-06 00:00:00", "2026-01-10 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 20.0, averageBreakTime) // (30 + 45 + 25) / 5
}

func TestGetUserAverageBreakTimeNoBreaks(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	userID := insertUser(t, db, "user-uuid-2", "testuser2", nil)

	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, breaks_duration_minutes, status)
		VALUES ('ws-a-1', ?, '2026-01-06 09:00:00', 480, 0, 'completed')
	`, userID)

	averageBreakTime, err := repo.GetUserAverageBreakTime(userID, "2026-01-06 00:00:00", "2026-01-10 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, averageBreakTime)
}

func TestGetUserAverageTimePerShift(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	userID := insertUser(t, db, "user-uuid-1", "testuser", nil)

	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-a-1', ?, '2026-01-06 09:00:00', 480, 'completed')
	`, userID)
	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-a-2', ?, '2026-01-07 09:00:00', 450, 'completed')
	`, userID)
	db.Exec(`
		INSERT INTO work_session_archived (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-ar-1', ?, '2026-01-08 09:00:00', 500, 'completed')
	`, userID)

	averageTime, totalShifts, totalTime, err := repo.GetUserAverageTimePerShift(userID, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 476.67, averageTime)
	assert.Equal(t, 3, totalShifts)
	assert.Equal(t, 1430, totalTime)
}

func TestGetUserAverageTimePerShiftNoData(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	userID := insertUser(t, db, "user-uuid-2", "testuser2", nil)

	averageTime, totalShifts, totalTime, err := repo.GetUserAverageTimePerShift(userID, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, averageTime)
	assert.Equal(t, 0, totalShifts)
	assert.Equal(t, 0, totalTime)
}

func TestGetUserAverageTimePerShiftSingleShift(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	userID := insertUser(t, db, "user-uuid-3", "testuser3", nil)

	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-a-1', ?, '2026-01-06 09:00:00', 400, 'completed')
	`, userID)

	averageTime, totalShifts, totalTime, err := repo.GetUserAverageTimePerShift(userID, "2026-01-06 00:00:00", "2026-01-06 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 400.0, averageTime)
	assert.Equal(t, 1, totalShifts)
	assert.Equal(t, 400, totalTime)
}

func TestGetUserAverageTimePerShiftMultipleWeeks(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	userID := insertUser(t, db, "user-uuid-4", "testuser4", nil)

	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-a-1', ?, '2026-01-06 09:00:00', 480, 'completed')
	`, userID)
	db.Exec(`
		INSERT INTO work_session_active (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-a-2', ?, '2026-01-13 09:00:00', 460, 'completed')
	`, userID)
	db.Exec(`
		INSERT INTO work_session_archived (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-ar-1', ?, '2026-01-20 09:00:00', 440, 'completed')
	`, userID)
	db.Exec(`
		INSERT INTO work_session_archived (uuid, user_id, clock_in, duration_minutes, status)
		VALUES ('ws-ar-2', ?, '2026-01-27 09:00:00', 420, 'completed')
	`, userID)

	averageTime, totalShifts, totalTime, err := repo.GetUserAverageTimePerShift(userID, "2026-01-01 00:00:00", "2026-01-31 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 450.0, averageTime)
	assert.Equal(t, 4, totalShifts)
	assert.Equal(t, 1800, totalTime)
}
