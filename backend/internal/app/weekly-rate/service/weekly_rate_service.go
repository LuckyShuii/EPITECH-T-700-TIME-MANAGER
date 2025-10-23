package service

import (
	"fmt"

	WeeklyRateModel "app/internal/app/weekly-rate/model"
	WeeklyRateRepository "app/internal/app/weekly-rate/repository"

	UserService "app/internal/app/user/service"

	Config "app/internal/config"

	"github.com/google/uuid"
)

type WeeklyRateService interface {
	GetAll() ([]WeeklyRateModel.WeeklyRate, error)
	Create(input WeeklyRateModel.CreateWeeklyRate) error
	Update(uuid string, input WeeklyRateModel.UpdateWeeklyRate) error
	Delete(uuid string) error
	AssignToUser(weeklyRateUUID string, userUUID string) error
}

type weeklyRateService struct {
	WeeklyRateRepo WeeklyRateRepository.WeeklyRateRepository
	UserService    UserService.UserService
}

func NewWeeklyRateService(repo WeeklyRateRepository.WeeklyRateRepository, userService UserService.UserService) WeeklyRateService {
	return &weeklyRateService{WeeklyRateRepo: repo, UserService: userService}
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
		return fmt.Errorf(Config.ErrorMessages()["WEEKLY_RATE_NOT_FOUND"]+": %w", err)
	}

	err = service.WeeklyRateRepo.Update(weeklyRateID, input)
	if err != nil {
		return fmt.Errorf("failed to update weekly rate: %w", err)
	}
	return nil
}

func (service *weeklyRateService) Delete(uuid string) error {
	weeklyRateID, err := service.WeeklyRateRepo.GetIDByUUID(uuid)
	if err != nil || weeklyRateID == 0 {
		return fmt.Errorf(Config.ErrorMessages()["WEEKLY_RATE_NOT_FOUND"]+": %w", err)
	}

	err = service.WeeklyRateRepo.Delete(uuid)
	if err != nil {
		return fmt.Errorf("failed to delete weekly rate: %w", err)
	}
	return nil
}

func (service *weeklyRateService) AssignToUser(weeklyRateUUID string, userUUID string) error {
	weeklyRateID, err := service.WeeklyRateRepo.GetIDByUUID(weeklyRateUUID)
	if err != nil || weeklyRateID == 0 {
		return fmt.Errorf(Config.ErrorMessages()["WEEKLY_RATE_NOT_FOUND"]+": %w", err)
	}

	userID, err := service.UserService.GetIdByUuid(userUUID)
	if err != nil || userID == 0 {
		return fmt.Errorf("failed to find user: %w", err)
	}

	err = service.WeeklyRateRepo.AssignToUser(weeklyRateID, userID)
	if err != nil {
		return fmt.Errorf("failed to assign weekly rate to user: %w", err)
	}
	return nil
}
