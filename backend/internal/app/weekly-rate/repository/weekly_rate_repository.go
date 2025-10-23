package repository

import (
	"fmt"

	"gorm.io/gorm"

	WeeklyRateModel "app/internal/app/weekly-rate/model"
)

type WeeklyRateRepository interface {
	GetAll() ([]WeeklyRateModel.WeeklyRate, error)
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
