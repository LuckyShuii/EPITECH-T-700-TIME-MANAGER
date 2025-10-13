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
	UpdateWorkSessionClocking(data WorkSessionModel.WorkSessionUpdate) (bool, error)
}

type workSessionService struct {
	WorkSessionRepo WorkSessionRepository.WorkSessionRepository
	UserService     UserService.UserService
}

func NewWorkSessionService(repo WorkSessionRepository.WorkSessionRepository, userService UserService.UserService) WorkSessionService {
	return &workSessionService{WorkSessionRepo: repo, UserService: userService}
}

func (service *workSessionService) UpdateWorkSessionClocking(data WorkSessionModel.WorkSessionUpdate) (bool, error) {
	log.Println("User UUID: ", data.UserUUID, " has just clocked, status: ", *data.IsClocked)

	/**
	 * Check if user exists and get the user ID based on his UUID
	 */
	userID, userErr := service.UserService.GetIdByUuid(data.UserUUID)
	if userErr != nil {
		return false, userErr
	}

	/**
	 * Check if a work session exists for the user
	 */
	workSessionFound, err := service.WorkSessionRepo.GetUserActiveWorkSession(userID, "active")
	if err != nil {
		return false, err
	}

	/**
	 * If user is clocking in & an active work session already exists, return an error message
	 */
	if workSessionFound.WorkSessionUUID != "" && *data.IsClocked {
		return false, fmt.Errorf("an active work session already exists for this user, cannot clock in again")
	}

	/**
	 * If user is clocking in & no active work session found, start a new one
	 */
	if workSessionFound.WorkSessionUUID == "" && *data.IsClocked {
		log.Println("No active work session found for the user: ", data.UserUUID, "Creating a new one...")
		service.WorkSessionRepo.CreateWorkSession(uuid.New().String(), userID, "active")
	}

	/**
	 * If user is clocking out & and no active work session found, return an error message
	 */
	if workSessionFound.WorkSessionUUID == "" && !*data.IsClocked {
		return false, fmt.Errorf("no active work session found for this user, cannot clock out")
	}

	/**
	 * If user is clocking out & an active work session found, close it
	 */
	if workSessionFound.WorkSessionUUID != "" && !*data.IsClocked {
		log.Println("Active work session found for the user: ", data.UserUUID, "Closing it...")

		t1, err := time.Parse(time.RFC3339Nano, workSessionFound.ClockIn)
		if err != nil {
			log.Println("parse error:", err)
		}

		t2 := time.Now().UTC()

		duration := t2.Sub(t1)
		minutes := duration.Minutes()

		rounded := math.Floor(minutes + 0.5)

		service.WorkSessionRepo.CompleteWorkSession(workSessionFound.WorkSessionUUID, userID, int(rounded))
	}

	return *data.IsClocked, nil
}
