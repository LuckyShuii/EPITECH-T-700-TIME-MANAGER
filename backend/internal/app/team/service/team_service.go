package service

import (
	"app/internal/app/team/model"
	"app/internal/app/team/repository"
	UserService "app/internal/app/user/service"

	"github.com/google/uuid"
)

type TeamService interface {
	GetTeams() ([]model.TeamReadAll, error)
	GetIdByUuid(id string) (int, error)
	GetUserIDsByTeamID(teamID int) ([]model.TeamMemberLight, error)
	GetTeamByUUID(uuid string) (model.TeamReadAll, error)
	DeleteTeamByID(id int) error
	RemoveUserFromTeam(teamID int, userID int) error
	CreateTeam(newTeam model.TeamCreate) error
	AddUsersToTeam(teamID int, members []model.TeamMemberCreate) error
	UpdateTeamByID(id int, updatedTeam model.TeamUpdate) error
	UpdateTeamUserManagerStatus(teamUUID string, userUUID string, isManager bool) error
}

type teamService struct {
	repo        repository.TeamRepository
	UserService UserService.UserService
}

func NewTeamService(repo repository.TeamRepository, userService UserService.UserService) TeamService {
	return &teamService{repo, userService}
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

func (service *teamService) CreateTeam(newTeam model.TeamCreate) error {
	teamUUID := uuid.New().String()

	if err := service.repo.CreateTeam(teamUUID, newTeam.Name, newTeam.Description); err != nil {
		return err
	}

	teamID, err := service.GetIdByUuid(teamUUID)
	if err != nil {
		return err
	}

	var members []model.TeamMemberCreate
	if newTeam.MemberUUIDs != nil {
		for _, member := range *newTeam.MemberUUIDs {
			userID, err := service.UserService.GetIdByUuid(member.UserUUID)
			if err != nil {
				return err
			}
			members = append(members, model.TeamMemberCreate{
				UserID:    userID,
				IsManager: member.IsManager,
			})
		}
	}

	if len(members) > 0 {
		if err := service.repo.AddMembersToTeam(teamID, members); err != nil {
			return err
		}
	}

	return nil
}

func (service *teamService) AddUsersToTeam(teamID int, members []model.TeamMemberCreate) error {
	return service.repo.AddMembersToTeam(teamID, members)
}

func (service *teamService) UpdateTeamByID(id int, updatedTeam model.TeamUpdate) error {
	return service.repo.UpdateTeamByID(id, updatedTeam)
}

func (service *teamService) UpdateTeamUserManagerStatus(teamUUID string, userUUID string, isManager bool) error {
	teamID, err := service.GetIdByUuid(teamUUID)
	if err != nil {
		return err
	}

	userID, err := service.UserService.GetIdByUuid(userUUID)
	if err != nil {
		return err
	}

	return service.repo.UpdateTeamUserManagerStatus(teamID, userID, isManager)
}

func (service *teamService) GetUserIDsByTeamID(teamID int) ([]model.TeamMemberLight, error) {
	return service.repo.FindUserIDsByTeamID(teamID)
}
