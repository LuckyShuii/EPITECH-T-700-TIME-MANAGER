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
)

// setupKPITestDB initializes a PostgreSQL test database using testcontainers.
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
		if err := pgC.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
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

	var db *gorm.DB
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err := db.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					break
				}
			}
		}

		if i < maxRetries-1 {
			waitTime := time.Duration(i+1) * time.Second
			t.Logf("Connection attempt %d failed, retrying in %v... Error: %v", i+1, waitTime, err)
			time.Sleep(waitTime)
		}
	}

	if err != nil {
		t.Fatalf("Failed to connect to Postgres after %d attempts: %v", maxRetries, err)
	}

	schema := `
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
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			phone_number VARCHAR(15),
			roles TEXT[] DEFAULT '{"employee"}',
			weekly_rate_id INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_weekly_rate FOREIGN KEY (weekly_rate_id) REFERENCES weekly_rate (id)
		);

		CREATE TABLE work_session_active (
			id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			uuid VARCHAR(36) NOT NULL UNIQUE,
			user_id INT NOT NULL,
			clock_in TIMESTAMP NOT NULL,
			clock_out TIMESTAMP,
			duration_minutes INT,
			breaks_duration_minutes INT DEFAULT 0,
			status work_session_status DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		);

		CREATE TABLE work_session_archived (
			id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			uuid VARCHAR(36) NOT NULL UNIQUE,
			user_id INT NOT NULL,
			clock_in TIMESTAMP NOT NULL,
			clock_out TIMESTAMP,
			duration_minutes INT,
			breaks_duration_minutes INT DEFAULT 0,
			status work_session_status DEFAULT 'active',
			archived_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		);
	`

	if err := db.Exec(schema).Error; err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	return db
}

func TestGetWeeklyRatesByUserIDAndDateRange(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	// Insert test user
	db.Exec(`INSERT INTO users (id, uuid, username) VALUES (1, 'user-uuid-1', 'testuser')`)

	// Insert work sessions in active table
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, status) 
		VALUES (1, '2026-01-06 09:00:00', 480, 'completed')`)
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, status) 
		VALUES (1, '2026-01-07 09:00:00', 450, 'completed')`)

	// Insert work sessions in archived table
	db.Exec(`INSERT INTO work_session_archived (user_id, clock_in, duration_minutes) 
		VALUES (1, '2026-01-08 09:00:00', 500)`)

	totalMinutes, err := repo.GetWeeklyRatesByUserIDAndDateRange(1, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 1430, totalMinutes) // 480 + 450 + 500
}

func TestGetWeeklyRatesByUserIDAndDateRangeNoData(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	db.Exec(`INSERT INTO users (id, uuid, username) VALUES (2, 'user-uuid-2', 'testuser2')`)

	totalMinutes, err := repo.GetWeeklyRatesByUserIDAndDateRange(2, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 0, totalMinutes)
}

func TestGetUserPresenceRate(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	// Insert test user with weekly rate
	db.Exec(`INSERT INTO weekly_rate (id, rate_name, amount) VALUES (1, 'Standard', 40)`)
	db.Exec(`INSERT INTO users (id, uuid, username, weekly_rate_id) VALUES (1, 'user-uuid-1', 'testuser', 1)`)

	// Insert work sessions totaling 2400 minutes = 40 hours
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, status) 
		VALUES (1, '2026-01-06 09:00:00', 1200, 'completed')`)
	db.Exec(`INSERT INTO work_session_archived (user_id, clock_in, duration_minutes) 
		VALUES (1, '2026-01-07 09:00:00', 1200)`)

	presenceRate, weeklyRateExpected, weeklyTimeDone, err := repo.GetUserPresenceRate(1, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 100.0, presenceRate)
	assert.Equal(t, 40.0, weeklyRateExpected)
	assert.Equal(t, 40.0, weeklyTimeDone)
}

func TestGetUserPresenceRateWithoutWeeklyRate(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	// Insert test user without weekly rate (should default to 40)
	db.Exec(`INSERT INTO users (id, uuid, username) VALUES (2, 'user-uuid-2', 'testuser2')`)

	// Insert work sessions totaling 1200 minutes = 20 hours
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, status) 
		VALUES (2, '2026-01-06 09:00:00', 1200, 'completed')`)

	presenceRate, weeklyRateExpected, weeklyTimeDone, err := repo.GetUserPresenceRate(2, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 50.0, presenceRate)
	assert.Equal(t, 40.0, weeklyRateExpected)
	assert.Equal(t, 20.0, weeklyTimeDone)
}

func TestGetUserAverageBreakTime(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	// Insert test user
	db.Exec(`INSERT INTO users (id, uuid, username) VALUES (1, 'user-uuid-1', 'testuser')`)

	// Insert work sessions with breaks
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, breaks_duration_minutes, status) 
		VALUES (1, '2026-01-06 09:00:00', 480, 30, 'completed')`)
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, breaks_duration_minutes, status) 
		VALUES (1, '2026-01-07 09:00:00', 450, 45, 'completed')`)
	db.Exec(`INSERT INTO work_session_archived (user_id, clock_in, duration_minutes, breaks_duration_minutes) 
		VALUES (1, '2026-01-08 09:00:00', 500, 25)`)

	averageBreakTime, err := repo.GetUserAverageBreakTime(1, "2026-01-06 00:00:00", "2026-01-10 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 20.0, averageBreakTime) // 100 minutes / 5 days = 20 minutes/day
}

func TestGetUserAverageBreakTimeNoBreaks(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	db.Exec(`INSERT INTO users (id, uuid, username) VALUES (2, 'user-uuid-2', 'testuser2')`)

	// Insert work sessions with no breaks
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, breaks_duration_minutes, status) 
		VALUES (2, '2026-01-06 09:00:00', 480, 0, 'completed')`)

	averageBreakTime, err := repo.GetUserAverageBreakTime(2, "2026-01-06 00:00:00", "2026-01-10 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, averageBreakTime)
}

func TestGetUserAverageTimePerShift(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	// Insert test user
	db.Exec(`INSERT INTO users (id, uuid, username) VALUES (1, 'user-uuid-1', 'testuser')`)

	// Insert work sessions
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, status) 
		VALUES (1, '2026-01-06 09:00:00', 480, 'completed')`)
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, status) 
		VALUES (1, '2026-01-07 09:00:00', 450, 'completed')`)
	db.Exec(`INSERT INTO work_session_archived (user_id, clock_in, duration_minutes) 
		VALUES (1, '2026-01-08 09:00:00', 500)`)

	averageTime, totalShifts, totalTime, err := repo.GetUserAverageTimePerShift(1, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 476.67, averageTime)
	assert.Equal(t, 3, totalShifts)
	assert.Equal(t, 1430, totalTime)
}

func TestGetUserAverageTimePerShiftNoData(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	db.Exec(`INSERT INTO users (id, uuid, username) VALUES (2, 'user-uuid-2', 'testuser2')`)

	averageTime, totalShifts, totalTime, err := repo.GetUserAverageTimePerShift(2, "2026-01-06 00:00:00", "2026-01-08 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, averageTime)
	assert.Equal(t, 0, totalShifts)
	assert.Equal(t, 0, totalTime)
}

func TestGetUserAverageTimePerShiftSingleShift(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	db.Exec(`INSERT INTO users (id, uuid, username) VALUES (3, 'user-uuid-3', 'testuser3')`)

	// Single shift of 400 minutes
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, status) 
		VALUES (3, '2026-01-06 09:00:00', 400, 'completed')`)

	averageTime, totalShifts, totalTime, err := repo.GetUserAverageTimePerShift(3, "2026-01-06 00:00:00", "2026-01-06 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 400.0, averageTime)
	assert.Equal(t, 1, totalShifts)
	assert.Equal(t, 400, totalTime)
}

func TestGetUserAverageTimePerShiftMultipleWeeks(t *testing.T) {
	db := setupKPITestDB(t)
	repo := repository.NewKPIRepository(db)

	db.Exec(`INSERT INTO users (id, uuid, username) VALUES (4, 'user-uuid-4', 'testuser4')`)

	// Multiple shifts across different weeks
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, status) 
		VALUES (4, '2026-01-06 09:00:00', 480, 'completed')`)
	db.Exec(`INSERT INTO work_session_active (user_id, clock_in, duration_minutes, status) 
		VALUES (4, '2026-01-13 09:00:00', 460, 'completed')`)
	db.Exec(`INSERT INTO work_session_archived (user_id, clock_in, duration_minutes) 
		VALUES (4, '2026-01-20 09:00:00', 440)`)
	db.Exec(`INSERT INTO work_session_archived (user_id, clock_in, duration_minutes) 
		VALUES (4, '2026-01-27 09:00:00', 420)`)

	averageTime, totalShifts, totalTime, err := repo.GetUserAverageTimePerShift(4, "2026-01-01 00:00:00", "2026-01-31 23:59:59")
	assert.NoError(t, err)
	assert.Equal(t, 450.0, averageTime)
	assert.Equal(t, 4, totalShifts)
	assert.Equal(t, 1800, totalTime)
}
