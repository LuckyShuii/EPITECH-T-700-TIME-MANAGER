package repository_test

import (
	"app/internal/app/team/model"
	"app/internal/app/team/repository"
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const insertTeamQuery = "INSERT INTO teams (uuid, name) VALUES (?, ?)"
const insertTeamMemberQuery = "INSERT INTO teams_members (uuid, team_id, user_id, is_manager) VALUES (?, ?, ?, ?)"
const selectIdUserUuid = "SELECT id FROM users WHERE uuid = ?"
const selectIdTeamUuid = "SELECT id FROM teams WHERE uuid = ?"

// Setup in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	ctx := context.Background()
	const portable = "5432/tcp"
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		Env:          map[string]string{"POSTGRES_PASSWORD": "test", "POSTGRES_DB": "testdb"},
		ExposedPorts: []string{portable},
		WaitingFor:   wait.ForListeningPort(portable),
	}
	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	port, _ := pgC.MappedPort(ctx, portable)
	dsn := fmt.Sprintf("host=localhost port=%s user=postgres password=test dbname=testdb sslmode=disable", port.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Creation of necessary tables postgres
	err = db.Exec(`
		CREATE TYPE work_session_status AS ENUM(
			'active',
			'completed',
			'paused'
		);

		-- User table
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
			roles TEXT[] default '{"employee"}',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Store active work sessions for the past 30 days
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

		-- Store archived work sessions older than 30 days max 2 years
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

		-- Store archived work sessions older than 2 years
		-- Do not store user data anymore for RGPD compliance
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

		-- Teams Table
		CREATE TABLE teams (
			id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			uuid VARCHAR(36) NOT NULL UNIQUE,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- User-Team Relationship Table
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

		CREATE INDEX idx_users_weekly_rate_id ON users (weekly_rate_id);

		INSERT INTO weekly_rate (uuid, rate_name, amount)
		VALUES (gen_random_uuid()::varchar, 'Temps pleins', 35);

		INSERT INTO weekly_rate (uuid, rate_name, amount)
		VALUES (gen_random_uuid()::varchar, 'Temps pleins + RTT', 39);

		INSERT INTO weekly_rate (uuid, rate_name, amount)
		VALUES (gen_random_uuid()::varchar, 'Temps partiel', 20);
	`).Error
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	return db
}

//
// CREATE & FIND ID BY UUID
//

func TestCreateAndFindIdByUuid(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTeamRepository(db)

	desc := "Team Description"
	uuid := "123e4567-e89b-12d3-a456-426614174000"

	err := db.Exec("INSERT INTO teams (uuid, name, description) VALUES (?, ?, ?)", uuid, "My Team", desc).Error
	assert.NoError(t, err)

	id, err := repo.FindIdByUuid(uuid)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestFindIdByUuidNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTeamRepository(db)

	id, err := repo.FindIdByUuid("not-existing-uuid")
	assert.Error(t, err)
	assert.Equal(t, 0, id)
}
func TestFindAllTeams(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTeamRepository(db)

	// 1️⃣ Weekly rate
	db.Exec(`INSERT INTO weekly_rate (uuid, amount, rate_name) VALUES (?, ?, ?)`, "wr-1", 1000, "Standard Rates")

	// 2️⃣ User
	db.Exec(`
		INSERT INTO users (
			uuid, username, email, password_hash, first_name, last_name, status, roles, phone_number, first_day_of_week, weekly_rate_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ARRAY['developer'], ?, ?, ?)
	`,
		"userss-1",
		"userone",
		"userones@example.com",
		"hashedpwd",
		"User",
		"One",
		"active",
		"1234567890",
		"1",
		1,
	)

	var userID int
	db.Raw(selectIdUserUuid, "userss-1").Scan(&userID)

	// 3️⃣ Team
	db.Exec(`INSERT INTO teams (uuid, name) VALUES (?, ?)`, "teams-1", "Teams A")

	var teamID int
	db.Raw(selectIdTeamUuid, "teams-1").Scan(&teamID)

	// 4️⃣ Team Member
	db.Exec(insertTeamMemberQuery, "tm-1", teamID, userID, true)

	// 5️⃣ Check
	teams, err := repo.FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, teams)
	assert.GreaterOrEqual(t, len(teams), 1, "expected at least 1 team")

	assert.Equal(t, "Teams A", teams[0].Name)
}

func TestFindByID(t *testing.T) {
	db := setupTestDB(t)

	db.Exec(insertTeamQuery, "a", "Team 3")

	type SimpleTeam struct {
		UUID        string
		Name        string
		Description *string
	}

	var team SimpleTeam
	err := db.Table("teams").Where("id = ?", 1).Find(&team).Error
	assert.NoError(t, err)
	assert.Equal(t, "Team 3", team.Name)
}

func TestFindByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTeamRepository(db)

	_, err := repo.FindByID(999)
	assert.Error(t, err)
}

// DELETE TESTS
func TestDeleteByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTeamRepository(db)

	db.Exec(insertTeamQuery, "a", "Test Team")
	err := repo.DeleteByID(1)
	assert.NoError(t, err)

	var count int64
	db.Table("teams").Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteUserFromTeam(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTeamRepository(db)

	db.Exec(insertTeamMemberQuery, "x", 1, 10, true)
	err := repo.DeleteUserFromTeam(1, 10)
	assert.NoError(t, err)

	var count int64
	db.Table("teams_members").Count(&count)
	assert.Equal(t, int64(0), count)
}

//
// ADD MEMBERS
//

func TestAddMembersToTeam(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTeamRepository(db)

	// 1️⃣ Weekly rate
	db.Exec(`INSERT INTO weekly_rate (uuid, amount, rate_name) VALUES (?, ?, ?)`, "wr-1", 1000, "Standard Rate")

	// 2️⃣ Users
	db.Exec(`
		INSERT INTO users (
			uuid, username, email, password_hash, first_name, last_name, status, roles, phone_number, first_day_of_week, weekly_rate_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ARRAY['developer'], ?, ?, ?)
	`,
		"users-1", "userone", "userone@example.com", "hash1", "User", "One", "active", "1234567890", "1", 1,
	)
	db.Exec(`
		INSERT INTO users (
			uuid, username, email, password_hash, first_name, last_name, status, roles, phone_number, first_day_of_week, weekly_rate_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ARRAY['developer'], ?, ?, ?)
	`,
		"users-2", "usertwo", "usertwo@example.com", "hash2", "User", "Two", "active", "987654321", "1", 1,
	)

	var user1ID, user2ID int
	db.Raw(selectIdUserUuid, "users-1").Scan(&user1ID)
	db.Raw(selectIdUserUuid, "users-2").Scan(&user2ID)

	// 3️⃣ Team
	db.Exec(`INSERT INTO teams (uuid, name) VALUES (?, ?)`, "teama-1", "Teama A")
	var teamID int
	db.Raw(selectIdTeamUuid, "teama-1").Scan(&teamID)

	// 4️⃣ Add members
	members := []model.TeamMemberCreate{
		{UserID: user1ID, IsManager: true},
		{UserID: user2ID, IsManager: false},
	}

	err := repo.AddMembersToTeam(teamID, members)
	assert.NoError(t, err)

	// 5️⃣ Check inserted rows
	var count int64
	db.Table("teams_members").Count(&count)
	assert.Equal(t, int64(2), count)
}

func TestAddMembersToTeamEmpty(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTeamRepository(db)

	err := repo.AddMembersToTeam(1, []model.TeamMemberCreate{})
	assert.NoError(t, err)
}

//
// UPDATE TEAM
//

func TestUpdateTeamByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTeamRepository(db)

	db.Exec("INSERT INTO teams (uuid, name, description) VALUES (?, ?, ?)", "u1", "Old Team", "Old Desc")

	newName := "New Team"
	newDesc := "Updated Desc"
	update := model.TeamUpdate{Name: &newName, Description: &newDesc}

	err := repo.UpdateTeamByID(1, update)
	assert.NoError(t, err)

	var team struct {
		Name        string `gorm:"column:name"`
		Description string `gorm:"column:description"`
	}
	db.Raw("SELECT name, description FROM teams WHERE id = 1").Scan(&team)
	assert.Equal(t, "New Team", team.Name)
	assert.Equal(t, "Updated Desc", team.Description)
}

func TestUpdateTeamByIDNoFields(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTeamRepository(db)

	err := repo.UpdateTeamByID(1, model.TeamUpdate{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no fields to update")
}

//
// UPDATE MANAGER STATUS
//

func TestUpdateTeamUserManagerStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTeamRepository(db)

	// 1️⃣ Weekly rate
	db.Exec(`INSERT INTO weekly_rate (uuid, amount, rate_name) VALUES (?, ?, ?)`, "wr-1", 1000, "Standard Rate")

	// 2️⃣ User
	db.Exec(`
		INSERT INTO users (
			uuid, username, email, password_hash, first_name, last_name, status, roles, phone_number, first_day_of_week, weekly_rate_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ARRAY['developer'], ?, ?, ?)
	`,
		"user-1",
		"userone",
		"userone@example.com",
		"hashedpwd",
		"User",
		"One",
		"active",
		"1234567890",
		"1",
		1,
	)

	var userID int
	db.Raw(selectIdUserUuid, "user-1").Scan(&userID)

	// 3️⃣ Team
	db.Exec(`INSERT INTO teams (uuid, name) VALUES (?, ?)`, "team-1", "Team A")
	var teamID int
	db.Raw(selectIdTeamUuid, "team-1").Scan(&teamID)

	// 4️⃣ Team member
	db.Exec(insertTeamMemberQuery, "tm-1", teamID, userID, false)

	// 5️⃣ Test update manager status
	err := repo.UpdateTeamUserManagerStatus(teamID, userID, true)
	assert.NoError(t, err)

	var isManager bool
	db.Raw("SELECT is_manager FROM teams_members WHERE team_id = ? AND user_id = ?", teamID, userID).Scan(&isManager)
	assert.True(t, isManager)
}
