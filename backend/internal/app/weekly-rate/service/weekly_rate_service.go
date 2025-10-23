package service

import (
	"fmt"
	"log"

	WeeklyRateModel "app/internal/app/weekly-rate/model"
	WeeklyRateRepository "app/internal/app/weekly-rate/repository"

	"github.com/google/uuid"
)

type WeeklyRateService interface {
	GetAll() ([]WeeklyRateModel.WeeklyRate, error)
	Create(input WeeklyRateModel.CreateWeeklyRate) error
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
		log.Printf("Error fetching weekly rates: %v", err)
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
		log.Printf("Error creating weekly rate: %v", err)
		return fmt.Errorf("failed to create weekly rate")
	}
	return nil
}
