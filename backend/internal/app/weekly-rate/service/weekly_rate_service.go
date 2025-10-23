package service

import (
	"fmt"
	"log"

	WeeklyRateModel "app/internal/app/weekly-rate/model"
	WeeklyRateRepository "app/internal/app/weekly-rate/repository"
)

type WeeklyRateService interface {
	GetAll() ([]WeeklyRateModel.WeeklyRate, error)
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
