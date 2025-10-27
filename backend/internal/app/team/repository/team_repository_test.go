package repository_test

import (
	"app/internal/app/team/model"
	"app/internal/app/team/repository"
	"app/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

const insertTeamQuery = "INSERT INTO teams (uuid, name) VALUES (?, ?)"
const insertTeamMemberQuery = "INSERT INTO teams_members (uuid, team_id, user_id, is_manager) VALUES (?, ?, ?, ?)"
const selectIdUserUuid = "SELECT id FROM users WHERE uuid = ?"
const selectIdTeamUuid = "SELECT id FROM teams WHERE uuid = ?"

//
// CREATE & FIND ID BY UUID
//

func TestCreateTeam(t *testing.T) {
	db := test.ResetDB(t)
	repo := repository.NewTeamRepository(db)

	desc := "Team Description"
	uuid := "123e4567-e89b-12d3-a456-426614174000"

	err := repo.CreateTeam(uuid, "My Team2", &desc)
	assert.NoError(t, err)

	var team struct {
		UUID        string
		Name        string
		Description *string
	}
	err = db.Table("teams").Where("uuid = ?", uuid).First(&team).Error
	assert.NoError(t, err)
	assert.Equal(t, "My Team2", team.Name)
	assert.NotNil(t, team.Description)
	assert.Equal(t, desc, *team.Description)
}

func TestCreateAndFindIdByUuid(t *testing.T) {
	db := test.ResetDB(t)
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
	db := test.ResetDB(t)
	repo := repository.NewTeamRepository(db)

	id, err := repo.FindIdByUuid("not-existing-uuid")
	assert.Error(t, err)
	assert.Equal(t, 0, id)
}
func TestFindAllTeams(t *testing.T) {
	db := test.ResetDB(t)
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
	db := test.ResetDB(t)
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
	db := test.ResetDB(t)
	repo := repository.NewTeamRepository(db)

	_, err := repo.FindByID(999)
	assert.Error(t, err)
}

// DELETE TESTS
func TestDeleteByID(t *testing.T) {
	db := test.ResetDB(t)
	repo := repository.NewTeamRepository(db)

	db.Exec(insertTeamQuery, "a", "Test Team")
	err := repo.DeleteByID(1)
	assert.NoError(t, err)

	var count int64
	db.Table("teams").Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteUserFromTeam(t *testing.T) {
	db := test.ResetDB(t)
	repo := repository.NewTeamRepository(db)

	// 1️⃣ Insert user
	db.Exec(`INSERT INTO users (uuid, username, email, password_hash, status)
	         VALUES (?, ?, ?, ?, ?)`,
		"x", "userx", "userx@example.com", "hashedpwd", "active")

	var userID int
	db.Raw(`SELECT id FROM users WHERE uuid = ?`, "x").Scan(&userID)

	// 2️⃣ Insert team
	db.Exec(insertTeamQuery, "team-12", "Team 12")

	var teamID int
	db.Raw(`SELECT id FROM teams WHERE uuid = ?`, "team-12").Scan(&teamID)

	// 3️⃣ Insert member
	db.Exec(insertTeamMemberQuery, "tm-1", teamID, userID, true)

	// 4️⃣ Delete and check
	err := repo.DeleteUserFromTeam(teamID, userID)
	assert.NoError(t, err)

	var count int64
	db.Table("teams_members").Count(&count)
	assert.Equal(t, int64(0), count)
}

//
// ADD MEMBERS
//

func TestAddMembersToTeam(t *testing.T) {
	db := test.ResetDB(t)
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
	db := test.ResetDB(t)
	repo := repository.NewTeamRepository(db)

	err := repo.AddMembersToTeam(1, []model.TeamMemberCreate{})
	assert.NoError(t, err)
}

//
// UPDATE TEAM
//

func TestUpdateTeamByID(t *testing.T) {
	db := test.ResetDB(t)
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
	db := test.ResetDB(t)
	repo := repository.NewTeamRepository(db)

	err := repo.UpdateTeamByID(1, model.TeamUpdate{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no fields to update")
}

//
// UPDATE MANAGER STATUS
//

func TestUpdateTeamUserManagerStatus(t *testing.T) {
	db := test.ResetDB(t)
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
