package test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	pgC    testcontainers.Container
	ctx    = context.Background()
	once   sync.Once // to avoid multiple initializations
	schema = `
	CREATE TYPE work_session_status AS ENUM('active', 'completed', 'paused');

	CREATE TABLE users (
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL UNIQUE,
		username VARCHAR(100) NOT NULL UNIQUE,
		email VARCHAR(320) NOT NULL UNIQUE,
		password_hash VARCHAR(100) NOT NULL,
		status VARCHAR(50) NOT NULL,
		first_day_of_week INT DEFAULT 1,
		dashboard_layout JSON DEFAULT NULL,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		phone_number VARCHAR(15),
		roles TEXT[] DEFAULT '{"employee"}',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE work_session_active (
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL UNIQUE,
		user_id INT NOT NULL,
		clock_in TIMESTAMP NOT NULL,
		clock_out TIMESTAMP,
		duration_minutes INT,
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
		status work_session_status DEFAULT 'active',
		archived_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);

	CREATE TABLE work_session_history (
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL UNIQUE,
		clock_in TIMESTAMP NOT NULL,
		clock_out TIMESTAMP,
		duration_minutes INT,
		status work_session_status DEFAULT 'active',
		archived_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE teams (
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL UNIQUE,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE teams_members (
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL UNIQUE,
		user_id INT NOT NULL,
		team_id INT NOT NULL,
		is_manager BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (team_id) REFERENCES teams (id) ON DELETE CASCADE
	);

	CREATE TABLE weekly_rate (
		id SERIAL PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL UNIQUE,
		rate_name VARCHAR(255) NOT NULL,
		amount SMALLINT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	ALTER TABLE users
	ADD COLUMN weekly_rate_id INT,
	ADD CONSTRAINT fk_weekly_rate FOREIGN KEY (weekly_rate_id) REFERENCES weekly_rate (id);

	CREATE TYPE break_status AS ENUM(
		'active',
		'completed'
	);

	CREATE TABLE breaks (
		id SERIAL PRIMARY KEY,
		uuid UUID NOT NULL UNIQUE,
		start_time TIMESTAMP NOT NULL,
		end_time TIMESTAMP,
		duration_minutes INTEGER,
		status break_status DEFAULT 'active',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		work_session_active_id INTEGER NOT NULL REFERENCES work_session_active (id) ON DELETE CASCADE
	);

	ALTER TABLE work_session_active
	ADD COLUMN breaks_duration_minutes INTEGER DEFAULT 0;

	ALTER TABLE work_session_archived
	ADD COLUMN breaks_duration_minutes INTEGER DEFAULT 0;

	ALTER TABLE work_session_history
	ADD COLUMN breaks_duration_minutes INTEGER DEFAULT 0;

	CREATE INDEX idx_users_weekly_rate_id ON users (weekly_rate_id);

	INSERT INTO weekly_rate (uuid, rate_name, amount) VALUES (gen_random_uuid()::varchar, 'Temps pleins', 35);
	INSERT INTO weekly_rate (uuid, rate_name, amount) VALUES (gen_random_uuid()::varchar, 'Temps pleins + RTT', 39);
	INSERT INTO weekly_rate (uuid, rate_name, amount) VALUES (gen_random_uuid()::varchar, 'Temps partiel', 20);

	INSERT INTO users (uuid, username, email, password_hash, status, dashboard_layout) VALUES ('123e4567-e89b-12d3-a456-426614174009', 'newuser', 'new@example.com', 'hashed_password', 'pending', '[{"layout": {"i":"a","x":0,"y":0,"w":4,"h":2,"minW":2,"minH":2,"static":false},"widgets":"a"}]');
	INSERT INTO users (uuid, username, email, password_hash, status) VALUES ('123e4567-e89b-12d3-a456-426614174010', 'delete_user', 'delete@example.com', 'hashed_password', 'pending');
	`
)

// Init once for all tests
func SetupTestDB() *gorm.DB {
	once.Do(func() {
		const port = "5432/tcp"

		req := testcontainers.ContainerRequest{
			Image: "postgres:16",
			Env: map[string]string{
				"POSTGRES_PASSWORD": "test",
				"POSTGRES_DB":       "testdb",
				"POSTGRES_USER":     "postgres",
			},
			ExposedPorts: []string{port},
			WaitingFor: wait.ForAll(
				wait.ForListeningPort(port),
				wait.ForLog("database system is ready to accept connections").
					WithStartupTimeout(60*time.Second),
			),
		}

		var err error
		pgC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			panic(fmt.Sprintf("❌ Failed to start container: %v", err))
		}

		// 🧩 Get host + port
		mappedPort, _ := pgC.MappedPort(ctx, port)
		host, _ := pgC.Host(ctx)
		dsn := fmt.Sprintf("host=%s port=%s user=postgres password=test dbname=testdb sslmode=disable", host, mappedPort.Port())

		// 🔁 Retry connection
		for i := 0; i < 5; i++ {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				sqlDB, _ := db.DB()
				if sqlDB.Ping() == nil {
					break
				}
			}
			time.Sleep(time.Second)
		}
		if err != nil {
			panic(fmt.Sprintf("❌ Failed to connect DB: %v", err))
		}

		// 🧱 Init schema once
		if err := db.Exec(schema).Error; err != nil {
			panic(fmt.Sprintf("❌ Failed to init schema: %v", err))
		}

		fmt.Println("✅ Test DB initialized successfully!")
	})

	return db
}

// ResetDB resets the database to a clean state before each test
func ResetDB(t *testing.T) *gorm.DB {
	t.Helper()

	if db == nil {
		db = SetupTestDB()
	}

	tables := []string{
		"teams_members", "teams", "users", "weekly_rate",
		"work_session_active", "work_session_archived", "work_session_history",
	}
	for _, tbl := range tables {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", tbl))
	}
	return db
}

// Cleanup is called at the end of all tests (optional)
func Cleanup() {
	if pgC != nil {
		_ = pgC.Terminate(ctx)
	}
	os.Exit(0)
}
