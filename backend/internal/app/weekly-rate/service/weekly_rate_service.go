package service

import (
	"fmt"

	WeeklyRateModel "app/internal/app/weekly-rate/model"
	WeeklyRateRepository "app/internal/app/weekly-rate/repository"

	"github.com/google/uuid"
)

type WeeklyRateService interface {
	GetAll() ([]WeeklyRateModel.WeeklyRate, error)
	Create(input WeeklyRateModel.CreateWeeklyRate) error
	Update(uuid string, input WeeklyRateModel.UpdateWeeklyRate) error
}

type weeklyRateService struct {
	WeeklyRateRepo WeeklyRateRepository.WeeklyRateRepository
}

func NewWeeklyRateService(repo WeeklyRateRepository.WeeklyRateRepository) WeeklyRateService {
	return &weeklyRateService{WeeklyRateRepo: repo}
}

func (service *weeklyRateService) GetAll() ([]WeeklyRateModel.WeeklyRate, error) {
	weeklyRates, err := service.WeeklyRateRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weekly rates")
	}
	return weeklyRates, nil
}

func (service *weeklyRateService) Create(input WeeklyRateModel.CreateWeeklyRate) error {
	var newWeeklyRate WeeklyRateModel.WeeklyRate
	newWeeklyRate.RateName = input.RateName
	newWeeklyRate.Amount = input.Amount

	// Generate UUID for the new weekly rate
	u := uuid.New().String()
	newWeeklyRate.UUID = &u

	err := service.WeeklyRateRepo.Create(newWeeklyRate)
	if err != nil {
		return fmt.Errorf("failed to create weekly rate")
	}
	return nil
}

func (service *weeklyRateService) Update(uuid string, input WeeklyRateModel.UpdateWeeklyRate) error {
	weeklyRateID, err := service.WeeklyRateRepo.GetIDByUUID(uuid)
	if err != nil {
		return fmt.Errorf("failed to find weekly rate")
	}

	err = service.WeeklyRateRepo.Update(weeklyRateID, input)
	if err != nil {
		return fmt.Errorf("failed to update weekly rate: %w", err)
	}
	return nil
}
