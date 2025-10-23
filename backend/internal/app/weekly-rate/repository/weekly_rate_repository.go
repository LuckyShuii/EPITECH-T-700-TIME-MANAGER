package repository

import (
	"fmt"

	"gorm.io/gorm"

	WeeklyRateModel "app/internal/app/weekly-rate/model"
)

type WeeklyRateRepository interface {
	GetAll() ([]WeeklyRateModel.WeeklyRate, error)
	GetIDByUUID(uuid string) (int, error)
	Create(input WeeklyRateModel.WeeklyRate) error
	Update(id int, input WeeklyRateModel.UpdateWeeklyRate) error
}

type weeklyRateRepository struct {
	db *gorm.DB
}

func NewWeeklyRateRepository(db *gorm.DB) WeeklyRateRepository {
	return &weeklyRateRepository{db}
}

func (repo *weeklyRateRepository) GetAll() ([]WeeklyRateModel.WeeklyRate, error) {
	var weeklyRates []WeeklyRateModel.WeeklyRate
	err := repo.db.Raw(`
		SELECT uuid, rate_name, amount
		FROM weekly_rate
	`).Scan(&weeklyRates).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weekly rates: %w", err)
	}
	return weeklyRates, nil
}

func (repo *weeklyRateRepository) GetIDByUUID(uuid string) (int, error) {
	var weeklyRateID int
	err := repo.db.Raw(`
		SELECT id
		FROM weekly_rate
		WHERE uuid = ?
	`, uuid).Scan(&weeklyRateID).Error
	if err != nil {
		return 0, fmt.Errorf("failed to find weekly rate ID by UUID: %w", err)
	}
	return weeklyRateID, nil
}

func (repo *weeklyRateRepository) Create(input WeeklyRateModel.WeeklyRate) error {
	result := repo.db.Exec(`
		INSERT INTO weekly_rate (uuid, rate_name, amount)
		VALUES (?, ?, ?)
	`, input.UUID, input.RateName, input.Amount)
	if result.Error != nil {
		return fmt.Errorf("failed to create weekly rate: %w", result.Error)
	}
	return nil
}

func (repo *weeklyRateRepository) Update(id int, input WeeklyRateModel.UpdateWeeklyRate) error {
	updateData := make(map[string]any)

	if input.RateName != "" {
		updateData["rate_name"] = input.RateName
	}
	if input.Amount != 0 {
		updateData["amount"] = input.Amount
	}

	result := repo.db.Model(&WeeklyRateModel.WeeklyRate{}).Table("weekly_rate").Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return fmt.Errorf("failed to update weekly rate: %w", result.Error)
	}
	return nil
}
