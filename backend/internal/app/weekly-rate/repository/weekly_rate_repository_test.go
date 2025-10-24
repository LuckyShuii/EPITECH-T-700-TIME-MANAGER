package repository_test

import (
	"app/internal/app/weekly-rate/model"
	"app/internal/app/weekly-rate/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database schema for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Create tables used by the repository
	err = db.Exec(`
		CREATE TABLE weekly_rate (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT NOT NULL,
			rate_name TEXT,
			amount INTEGER
		);
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT,
			weekly_rate_id INTEGER
		);
	`).Error
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	return db
}

func TestCreateWeeklyRate(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWeeklyRateRepository(db)

	uuid := "uuid-1"
	newRate := model.WeeklyRate{
		UUID:     &uuid,
		RateName: "Standard",
		Amount:   450,
	}

	err := repo.Create(newRate)
	assert.NoError(t, err)

	var count int64
	db.Table("weekly_rate").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestGetAllWeeklyRates(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWeeklyRateRepository(db)

	db.Exec(`
		INSERT INTO weekly_rate (uuid, rate_name, amount)
		VALUES ('uuid-1', 'Basic', 100), ('uuid-2', 'Premium', 200)
	`)

	rates, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, rates, 2)
	assert.Equal(t, "Basic", rates[0].RateName)
}

func TestGetIDByUUID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWeeklyRateRepository(db)

	db.Exec(`
		INSERT INTO weekly_rate (uuid, rate_name, amount)
		VALUES ('uuid-3', 'Custom', 300)
	`)

	id, err := repo.GetIDByUUID("uuid-3")
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestUpdateWeeklyRate(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWeeklyRateRepository(db)

	db.Exec(`
		INSERT INTO weekly_rate (uuid, rate_name, amount)
		VALUES ('uuid-4', 'OldName', 100)
	`)

	updateData := model.UpdateWeeklyRate{
		RateName: "NewName",
		Amount:   200,
	}

	err := repo.Update(1, updateData)
	assert.NoError(t, err)

	var result struct {
		RateName string
		Amount   int
	}
	db.Raw("SELECT rate_name, amount FROM weekly_rate WHERE id = 1").Scan(&result)

	assert.Equal(t, "NewName", result.RateName)
	assert.Equal(t, 200, result.Amount)
}

func TestDeleteWeeklyRate(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWeeklyRateRepository(db)

	db.Exec(`
		INSERT INTO weekly_rate (uuid, rate_name, amount)
		VALUES ('uuid-5', 'Temp', 123)
	`)

	err := repo.Delete("uuid-5")
	assert.NoError(t, err)

	var count int64
	db.Table("weekly_rate").Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestAssignToUser(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWeeklyRateRepository(db)

	// Create a rate and user
	db.Exec(`INSERT INTO weekly_rate (uuid, rate_name, amount) VALUES ('uuid-6', 'Assign', 999)`)
	db.Exec(`INSERT INTO users (uuid) VALUES ('user-1')`)

	err := repo.AssignToUser(1, 1)
	assert.NoError(t, err)

	var user struct {
		WeeklyRateID int
	}
	db.Raw("SELECT weekly_rate_id FROM users WHERE id = 1").Scan(&user)

	assert.Equal(t, 1, user.WeeklyRateID)
}
