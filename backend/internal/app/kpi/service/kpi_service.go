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
	GetWorkSessionTeamWeeklyTotal(startDate string, endDate string, teamUUID string) (int, error)
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

func (service *kpiService) GetWorkSessionTeamWeeklyTotal(startDate string, endDate string, teamUUID string) (int, error) {
	teamID, err := service.TeamService.GetIdByUuid(teamUUID)
	if err != nil {
		return 0, err
	}

	userIDs, err := service.TeamService.GetUserIDsByTeamID(teamID)
	if err != nil {
		return 0, err
	}

	totalWeeklyRates := 0
	for _, userID := range userIDs {
		weeklyRates, err := service.KPIRepository.GetWeeklyRatesByUserIDAndDateRange(userID, startDate, endDate)
		if err != nil {
			return 0, err
		}
		totalWeeklyRates += weeklyRates
	}

	return totalWeeklyRates, nil
}
