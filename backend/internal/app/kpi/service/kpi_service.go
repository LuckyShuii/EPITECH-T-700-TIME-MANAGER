package service

import (
	BreakService "app/internal/app/break/service"
	TeamService "app/internal/app/team/service"
	UserService "app/internal/app/user/service"
	WeeklyRateService "app/internal/app/weekly-rate/service"

	KPIRepository "app/internal/app/kpi/repository"
)

type KPIService interface {
	GetWorkSessionUserWeeklyTotal(startDate string, endDate string, userUUID string) (int, error)
}

type kpiService struct {
	BreakService      BreakService.BreakService
	TeamService       TeamService.TeamService
	UserService       UserService.UserService
	WeeklyRateService WeeklyRateService.WeeklyRateService
	KPIRepository     KPIRepository.KPIRepository
}

func NewKPIService(breakService BreakService.BreakService, teamService TeamService.TeamService, userService UserService.UserService, weeklyRateService WeeklyRateService.WeeklyRateService, kpiRepository KPIRepository.KPIRepository) KPIService {
	return &kpiService{
		BreakService:      breakService,
		TeamService:       teamService,
		UserService:       userService,
		WeeklyRateService: weeklyRateService,
		KPIRepository:     kpiRepository,
	}
}

func (service *kpiService) GetWorkSessionUserWeeklyTotal(startDate string, endDate string, userUUID string) (int, error) {
	userID, err := service.UserService.GetIdByUuid(userUUID)
	if err != nil {
		return 0, err
	}

	weeklyRates, err := service.KPIRepository.GetWeeklyRatesByUserIDAndDateRange(userID, startDate, endDate)
	if err != nil {
		return 0, err
	}

	return weeklyRates, nil
}
