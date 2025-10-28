package service

import (
	BreakService "app/internal/app/break/service"
	TeamService "app/internal/app/team/service"
	UserService "app/internal/app/user/service"
	WeeklyRateService "app/internal/app/weekly-rate/service"
	"log"

	"app/internal/app/kpi/model"
	KPIRepository "app/internal/app/kpi/repository"
)

type KPIService interface {
	GetWorkSessionUserWeeklyTotal(startDate string, endDate string, userUUID string) (int, error)
	GetWorkSessionTeamWeeklyTotal(startDate string, endDate string, teamUUID string) (model.KPIWorkSessionTeamWeeklyTotalResponse, error)
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

func (service *kpiService) GetWorkSessionTeamWeeklyTotal(startDate string, endDate string, teamUUID string) (model.KPIWorkSessionTeamWeeklyTotalResponse, error) {
	teamID, err := service.TeamService.GetIdByUuid(teamUUID)
	if err != nil {
		return model.KPIWorkSessionTeamWeeklyTotalResponse{}, err
	}

	users, err := service.TeamService.GetUserIDsByTeamID(teamID)
	log.Printf("Found users for team %s: %v", teamUUID, users)
	if err != nil {
		return model.KPIWorkSessionTeamWeeklyTotalResponse{}, err
	}

	memberWeeklyRates := make([]model.KPIWorkSessionTeamMemberWeeklyTotal, 0)
	for _, user := range users {
		log.Printf("Calculating weekly rates for user %s (%s)", user.UserUUID, user.FirstName+" "+user.LastName)
		weeklyRates, err := service.KPIRepository.GetWeeklyRatesByUserIDAndDateRange(user.UserID, startDate, endDate)
		if err != nil {
			return model.KPIWorkSessionTeamWeeklyTotalResponse{}, err
		}
		memberWeeklyRates = append(memberWeeklyRates, model.KPIWorkSessionTeamMemberWeeklyTotal{
			UserUUID:  user.UserUUID,
			TotalTime: weeklyRates,
		})
		log.Printf("User %s (%s) has total time %d", user.UserUUID, user.FirstName+" "+user.LastName, weeklyRates)
	}

	totalTeamTime := 0
	for _, member := range memberWeeklyRates {
		totalTeamTime += member.TotalTime
	}

	return model.KPIWorkSessionTeamWeeklyTotalResponse{
		TotalTime: totalTeamTime,
		StartDate: startDate,
		EndDate:   endDate,
		TeamUUID:  teamUUID,
		Members:   memberWeeklyRates,
	}, nil
}
