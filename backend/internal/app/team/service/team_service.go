package service

import (
	"app/internal/app/team/model"
	"app/internal/app/team/repository"
)

type TeamService interface {
	GetTeams() ([]model.TeamReadAll, error)
	GetIdByUuid(id string) (int, error)
	GetTeamByUUID(uuid string) (model.TeamReadAll, error)
	DeleteTeamByID(id int) error
	RemoveUserFromTeam(teamID int, userID int) error
}

type teamService struct {
	repo repository.TeamRepository
}

func NewTeamService(repo repository.TeamRepository) TeamService {
	return &teamService{repo}
}

func (service *teamService) GetTeams() ([]model.TeamReadAll, error) {
	return service.repo.FindAll()
}

func (service *teamService) GetIdByUuid(id string) (int, error) {
	teamID, err := service.repo.FindIdByUuid(id)
	return teamID, err
}

func (service *teamService) GetTeamByUUID(uuid string) (model.TeamReadAll, error) {
	teamID, err := service.GetIdByUuid(uuid)
	if err != nil {
		return model.TeamReadAll{}, err
	}
	return service.repo.FindByID(teamID)
}

func (service *teamService) DeleteTeamByID(id int) error {
	return service.repo.DeleteByID(id)
}

func (service *teamService) RemoveUserFromTeam(teamID int, userID int) error {
	return service.repo.DeleteUserFromTeam(teamID, userID)
}
