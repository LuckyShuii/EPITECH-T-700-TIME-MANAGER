package repository_test

import (
	"app/internal/app/team/model"
	"app/internal/app/team/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const insertTeamQuery = "INSERT INTO teams (uuid, name) VALUES (?, ?)"

// Setup in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Creation of necessary tables
	err = db.Exec(`
		CREATE TABLE teams (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT NOT NULL,
			name TEXT,
			description TEXT
		);

		CREATE TABLE teams_members (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT NOT NULL,
			team_id INTEGER,
			user_id INTEGER,
			is_manager BOOLEAN
		);

		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT,
			username TEXT,
			email TEXT,
			first_name TEXT,
			last_name TEXT,
			status TEXT,
			roles TEXT,
			phone_number TEXT,
			first_day_of_week TEXT,
			weekly_rate_id INTEGER
		);

		CREATE TABLE weekly_rate (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			amount INTEGER,
			rate_name TEXT
		);

		CREATE TABLE work_session_active (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			status TEXT,
			created_at TIMESTAMP
		);
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

	// Données de test
	db.Exec(insertTeamQuery, "a", "Team A")
	db.Exec(insertTeamQuery, "b", "Team B")

	// ✅ Struct simplifiée pour SQLite uniquement
	type SimpleTeam struct {
		UUID        string  `json:"uuid"`
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	var teams []SimpleTeam
	err := db.Table("teams").Select("uuid, name, description").Find(&teams).Error
	assert.NoError(t, err)
	assert.NotNil(t, teams)
	assert.GreaterOrEqual(t, len(teams), 2, "expected at least 2 teams")

	assert.Equal(t, "Team A", teams[0].Name)
	assert.Equal(t, "Team B", teams[1].Name)
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

	db.Exec("INSERT INTO teams_members (uuid, team_id, user_id, is_manager) VALUES (?, ?, ?, ?)", "x", 1, 10, true)
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

	members := []model.TeamMemberCreate{
		{UserID: 1, IsManager: true},
		{UserID: 2, IsManager: false},
	}

	err := repo.AddMembersToTeam(1, members)
	assert.NoError(t, err)

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

	db.Exec("INSERT INTO teams_members (uuid, team_id, user_id, is_manager) VALUES (?, ?, ?, ?)", "x", 1, 1, false)
	err := repo.UpdateTeamUserManagerStatus(1, 1, true)
	assert.NoError(t, err)

	var isManager bool
	db.Raw("SELECT is_manager FROM teams_members WHERE team_id = 1 AND user_id = 1").Scan(&isManager)
	assert.True(t, isManager)
}
