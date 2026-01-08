package service

import (
	BreakService "app/internal/app/break/service"
	TeamService "app/internal/app/team/service"
	UserService "app/internal/app/user/service"
	WeeklyRateService "app/internal/app/weekly-rate/service"
	"fmt"
	"os"

	"app/internal/app/kpi/model"
	KPIRepository "app/internal/app/kpi/repository"

	"app/internal/app/kpi/export"
)

type KPIService interface {
	GetWorkSessionUserWeeklyTotal(startDate string, endDate string, userUUID string) (int, error)
	GetWorkSessionTeamWeeklyTotal(startDate string, endDate string, teamUUID string) (model.KPIWorkSessionTeamWeeklyTotalResponse, error)
	GetPresenceRate(startDate string, endDate string, userUUID string) (model.KPIPresenceRateResponse, error)
	ExportKPIData(startDate string, endDate string, requestedByUUID string, kpiType string, uuidToSearch string) (model.KPIExportResponse, error)
	GetAverageBreakTime(startDate string, endDate string, userUUID string) (model.KPIAverageBreakTimeResponse, error)
	GetAverageTimePerShift(startDate string, endDate string, userUUID string) (model.KPIAverageTimePerShiftResponse, error)
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
	teamID, userIdErr := service.TeamService.GetIdByUuid(teamUUID)
	if userIdErr != nil {
		return model.KPIWorkSessionTeamWeeklyTotalResponse{}, userIdErr
	}

	team, teamUuidErr := service.TeamService.GetTeamByUUID(teamUUID)
	if teamUuidErr != nil {
		return model.KPIWorkSessionTeamWeeklyTotalResponse{}, teamUuidErr
	}

	users, err := service.TeamService.GetUserIDsByTeamID(teamID)
	if err != nil {
		return model.KPIWorkSessionTeamWeeklyTotalResponse{}, err
	}

	memberWeeklyRates := make([]model.KPIWorkSessionTeamMemberWeeklyTotal, 0)
	for _, user := range users {
		weeklyRates, err := service.KPIRepository.GetWeeklyRatesByUserIDAndDateRange(user.UserID, startDate, endDate)
		if err != nil {
			return model.KPIWorkSessionTeamWeeklyTotalResponse{}, err
		}
		memberWeeklyRates = append(memberWeeklyRates, model.KPIWorkSessionTeamMemberWeeklyTotal{
			UserUUID:  user.UserUUID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			TotalTime: weeklyRates,
		})
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
		TeamName:  team.Name,
		Members:   memberWeeklyRates,
	}, nil
}

func (service *kpiService) GetPresenceRate(startDate string, endDate string, userUUID string) (model.KPIPresenceRateResponse, error) {
	userID, err := service.UserService.GetIdByUuid(userUUID)
	if err != nil {
		return model.KPIPresenceRateResponse{}, err
	}

	presenceRate, weeklyRateExpected, weeklyTimeDone, err := service.KPIRepository.GetUserPresenceRate(userID, startDate, endDate)
	if err != nil {
		return model.KPIPresenceRateResponse{}, err
	}

	data, err := service.UserService.GetUserByUUID(userUUID)
	if err != nil {
		return model.KPIPresenceRateResponse{}, err
	}

	return model.KPIPresenceRateResponse{
		FirstName:          data.FirstName,
		LastName:           data.LastName,
		UserUUID:           userUUID,
		PresenceRate:       presenceRate,
		WeeklyRateExpected: weeklyRateExpected,
		WeeklyTimeDone:     weeklyTimeDone,
	}, nil
}

func (service *kpiService) ExportKPIData(startDate string, endDate string, requestedByUUID string, kpiType string, uuidToSearch string) (model.KPIExportResponse, error) {
	filename := fmt.Sprintf("kpi_%s_%s.csv", kpiType, requestedByUUID)
	tmpPath := "/app/tmp/kpi/" + filename
	finalPath := "/app/data/kpi/" + filename

	var headers []string
	var rows [][]string

	// generate data based on kpiType
	switch kpiType {
	case "work_session_user_weekly_total":
		total, err := service.GetWorkSessionUserWeeklyTotal(startDate, endDate, uuidToSearch)
		if err != nil {
			return model.KPIExportResponse{}, err
		}

		user, err := service.UserService.GetUserByUUID(uuidToSearch)
		if err != nil {
			return model.KPIExportResponse{}, err
		}

		headers = []string{"user uuid", "firstname", "lastname", "start date", "end date", "total minutes"}
		rows = [][]string{
			{uuidToSearch, user.FirstName, user.LastName, startDate, endDate, fmt.Sprint(total)},
		}

	case "work_session_team_weekly_total":
		data, err := service.GetWorkSessionTeamWeeklyTotal(startDate, endDate, uuidToSearch)
		if err != nil {
			return model.KPIExportResponse{}, err
		}

		headers = []string{"team uuid", "team name", "start date", "end date", "total time"}
		rows = [][]string{
			{data.TeamUUID, data.TeamName, startDate, endDate, fmt.Sprint(data.TotalTime)},
		}

		headers = append(headers, "member user uuid", "firstname", "lastname", "member total minutes")
		for _, member := range data.Members {
			rows = append(rows, []string{"", "", "", "", "", member.UserUUID, member.FirstName, member.LastName, fmt.Sprint(member.TotalTime)})
		}

	case "presence_rate":
		data, err := service.GetPresenceRate(startDate, endDate, uuidToSearch)
		if err != nil {
			return model.KPIExportResponse{}, err
		}

		headers = []string{"user uuid", "firstname", "lastname", "start date", "end date", "presence rate", "weekly rate expected", "weekly time done"}
		rows = [][]string{
			{
				data.UserUUID,
				data.FirstName,
				data.LastName,
				startDate,
				endDate,
				fmt.Sprintf("%.2f", data.PresenceRate),
				fmt.Sprint(data.WeeklyRateExpected),
				fmt.Sprint(data.WeeklyTimeDone),
			},
		}

	case "weekly_average_break_time":
		data, err := service.GetAverageBreakTime(startDate, endDate, uuidToSearch)
		if err != nil {
			return model.KPIExportResponse{}, err
		}

		headers = []string{"user uuid", "firstname", "lastname", "start date", "end date", "average break time (minutes)"}
		rows = [][]string{
			{
				data.UserUUID,
				data.FirstName,
				data.LastName,
				startDate,
				endDate,
				fmt.Sprintf("%.2f", data.AverageBreakTime),
			},
		}

	case "average_time_per_shift":
		data, err := service.GetAverageTimePerShift(startDate, endDate, uuidToSearch)
		if err != nil {
			return model.KPIExportResponse{}, err
		}

		headers = []string{"user uuid", "firstname", "lastname", "start date", "end date", "average time per shift (minutes)", "total shifts", "total time (minutes)"}
		rows = [][]string{
			{
				data.UserUUID,
				data.FirstName,
				data.LastName,
				startDate,
				endDate,
				fmt.Sprintf("%.2f", data.AverageTimePerShift),
				fmt.Sprint(data.TotalShifts),
				fmt.Sprint(data.TotalTime),
			},
		}

	default:
		return model.KPIExportResponse{}, fmt.Errorf("unknown KPI type: %s", kpiType)
	}

	if err := exportCSV(headers, rows, tmpPath); err != nil {
		return model.KPIExportResponse{}, err
	}

	// if the file already exists in /data/kpi/, delete it
	if _, err := os.Stat(finalPath); err == nil {
		if err := os.Remove(finalPath); err != nil {
			return model.KPIExportResponse{}, err
		}
	}

	// move file from tmp to final destination
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return model.KPIExportResponse{}, err
	}

	// return file info
	return model.KPIExportResponse{
		File: filename,
		URL:  "/api/kpi/files/" + filename,
	}, nil
}

func exportCSV(headers []string, rows [][]string, path string) error {
	return export.ExportCSV(headers, rows, path)
}

func (service *kpiService) GetAverageBreakTime(startDate string, endDate string, userUUID string) (model.KPIAverageBreakTimeResponse, error) {
	userID, err := service.UserService.GetIdByUuid(userUUID)
	if err != nil {
		return model.KPIAverageBreakTimeResponse{}, err
	}

	averageBreakTime, err := service.KPIRepository.GetUserAverageBreakTime(userID, startDate, endDate)
	if err != nil {
		return model.KPIAverageBreakTimeResponse{}, err
	}

	data, err := service.UserService.GetUserByUUID(userUUID)
	if err != nil {
		return model.KPIAverageBreakTimeResponse{}, err
	}

	return model.KPIAverageBreakTimeResponse{
		FirstName:        data.FirstName,
		LastName:         data.LastName,
		UserUUID:         userUUID,
		AverageBreakTime: averageBreakTime,
		StartDate:        startDate,
		EndDate:          endDate,
	}, nil
}

func (service *kpiService) GetAverageTimePerShift(startDate string, endDate string, userUUID string) (model.KPIAverageTimePerShiftResponse, error) {
	userID, err := service.UserService.GetIdByUuid(userUUID)
	if err != nil {
		return model.KPIAverageTimePerShiftResponse{}, err
	}

	averageTimePerShift, totalShifts, totalTime, err := service.KPIRepository.GetUserAverageTimePerShift(userID, startDate, endDate)
	if err != nil {
		return model.KPIAverageTimePerShiftResponse{}, err
	}

	data, err := service.UserService.GetUserByUUID(userUUID)
	if err != nil {
		return model.KPIAverageTimePerShiftResponse{}, err
	}

	return model.KPIAverageTimePerShiftResponse{
		FirstName:           data.FirstName,
		LastName:            data.LastName,
		UserUUID:            userUUID,
		AverageTimePerShift: averageTimePerShift,
		TotalShifts:         totalShifts,
		TotalTime:           totalTime,
		StartDate:           startDate,
		EndDate:             endDate,
	}, nil
}
