package service

import (
	WorkSessionModel "app/internal/app/work-session/model"
	"fmt"
	"log"
	"math"
	"time"

	UserService "app/internal/app/user/service"
	WorkSessionRepository "app/internal/app/work-session/repository"

	"github.com/google/uuid"
)

type WorkSessionService interface {
	UpdateWorkSessionClocking(data WorkSessionModel.WorkSessionUpdate) (WorkSessionModel.WorkSessionUpdateResponse, error)
}

type workSessionService struct {
	WorkSessionRepo WorkSessionRepository.WorkSessionRepository
	UserService     UserService.UserService
}

func NewWorkSessionService(repo WorkSessionRepository.WorkSessionRepository, userService UserService.UserService) WorkSessionService {
	return &workSessionService{WorkSessionRepo: repo, UserService: userService}
}

func (service *workSessionService) UpdateWorkSessionClocking(data WorkSessionModel.WorkSessionUpdate) (WorkSessionModel.WorkSessionUpdateResponse, error) {
	var response WorkSessionModel.WorkSessionUpdateResponse

	/**
	 * Check if user exists and get the user ID based on his UUID
	 */
	userID, userErr := service.UserService.GetIdByUuid(data.UserUUID)
	if userErr != nil {
		response.Success = false
		return response, userErr
	}

	/**
	 * Check if a work session exists for the user
	 */
	workSessionFound, err := service.WorkSessionRepo.GetUserActiveWorkSession(userID, "active")
	if err != nil {
		response.Success = false
		return response, err
	}

	if workSessionFound.WorkSessionUUID != "" {
		response.ClockInTime = workSessionFound.ClockIn
	}

	/**
	 * If user is clocking in & an active work session already exists, return an error message
	 */
	if workSessionFound.WorkSessionUUID != "" && *data.IsClocked {
		return response, fmt.Errorf("an active work session already exists for this user, cannot clock in again")
	}

	/**
	 * If user is clocking in & no active work session found, start a new one
	 */
	if workSessionFound.WorkSessionUUID == "" && *data.IsClocked {
		response.ClockInTime = time.Now().In(time.FixedZone("Europe/Paris", 2*60*60)).Format(time.RFC3339Nano)
		service.WorkSessionRepo.CreateWorkSession(uuid.New().String(), userID, "active")
	}

	/**
	 * If user is clocking out & and no active work session found, return an error message
	 */
	if workSessionFound.WorkSessionUUID == "" && !*data.IsClocked {
		response.Success = false
		return response, fmt.Errorf("no active work session found for this user, cannot clock out")
	}

	/**
	 * If user is clocking out & an active work session found, close it
	 */
	if workSessionFound.WorkSessionUUID != "" && !*data.IsClocked {
		t1, err := time.Parse(time.RFC3339Nano, workSessionFound.ClockIn)
		if err != nil {
			log.Println("parse error:", err)
		}

		loc, _ := time.LoadLocation("Europe/Paris")

		t2 := time.Now().In(loc)

		duration := t2.Sub(t1)
		minutes := duration.Minutes()

		rounded := math.Floor(minutes + 0.5)

		service.WorkSessionRepo.CompleteWorkSession(workSessionFound.WorkSessionUUID, userID, int(rounded))

		clockOutTimeStr := t2.Format(time.RFC3339Nano)
		response.ClockOutTime = &clockOutTimeStr
	}

	response.Success = true
	if !*data.IsClocked {
		response.Status = "clocked_out"
	} else {
		response.Status = "clocked_in"
	}

	return response, nil
}
