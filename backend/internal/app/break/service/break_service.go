package service

import (
	"fmt"
	"log"
	"math"
	"time"

	BreakModel "app/internal/app/break/model"
	BreakRepository "app/internal/app/break/repository"
	WorkSessionRepository "app/internal/app/work-session/repository"

	"github.com/google/uuid"
)

type BreakService interface {
	UpdateBreakClocking(data BreakModel.BreakUpdate) (BreakModel.BreakUpdateResponse, error)
}

type breakService struct {
	BreakRepo       BreakRepository.BreakRepository
	WorkSessionRepo WorkSessionRepository.WorkSessionRepository
}

func NewBreakService(repo BreakRepository.BreakRepository, workSessionRepo WorkSessionRepository.WorkSessionRepository) BreakService {
	return &breakService{BreakRepo: repo, WorkSessionRepo: workSessionRepo}
}

func (service *breakService) UpdateBreakClocking(data BreakModel.BreakUpdate) (BreakModel.BreakUpdateResponse, error) {
	var response BreakModel.BreakUpdateResponse

	/**
	 * Check if work session exists and get the work session ID based on its UUID
	 */
	WorkSessionID, err := service.WorkSessionRepo.FindIdByUuid(data.WorkSessionUUID)
	if err != nil {
		response.Success = false
		return response, err
	}

	/**
	 * Check if a break session exists for the work session
	 */
	breakSessionFound, err := service.BreakRepo.GetWorkSessionBreak(WorkSessionID, "active")
	if err != nil {
		response.Success = false
		return response, err
	}

	if breakSessionFound.BreakUUID != "" {
		response.StartTime = breakSessionFound.StartTime
	}

	/**
	 * If user starting a break & an active break session already exists, return an error message
	 */
	if breakSessionFound.BreakUUID != "" && *data.IsBreaking {
		return response, fmt.Errorf("an active break session already exists for this work session, cannot start a break again")
	}

	/**
	 * If user is starting a break & no active break session found, start a new one
	 */
	if breakSessionFound.BreakUUID == "" && *data.IsBreaking {
		response.StartTime = time.Now().In(time.FixedZone("Europe/Paris", 2*60*60)).Format(time.RFC3339Nano)
		service.BreakRepo.CreateBreak(uuid.New().String(), WorkSessionID, "active")
		service.WorkSessionRepo.UpdateWorkSessionStatus(data.WorkSessionUUID, "paused")
	}

	/**
	 * If user is stopping a break & and no active break session found, return an error message
	 */
	if breakSessionFound.BreakUUID == "" && !*data.IsBreaking {
		response.Success = false
		return response, fmt.Errorf("no active break session found for this work session, cannot stop break")
	}

	/**
	 * If user is stopping a break & an active break session found, close it
	 */
	if breakSessionFound.BreakUUID != "" && !*data.IsBreaking {
		t1, err := time.Parse(time.RFC3339Nano, breakSessionFound.StartTime)
		if err != nil {
			log.Println("parse error:", err)
		}

		loc, _ := time.LoadLocation("Europe/Paris")

		t2 := time.Now().In(loc)

		duration := t2.Sub(t1)
		minutes := duration.Minutes()

		rounded := math.Floor(minutes + 0.5)

		service.BreakRepo.CompleteBreak(breakSessionFound.BreakUUID, WorkSessionID, int(rounded))
		service.WorkSessionRepo.UpdateWorkSessionStatus(data.WorkSessionUUID, "active")

		clockOutTimeStr := t2.Format(time.RFC3339Nano)
		response.EndTime = &clockOutTimeStr
	}

	response.Success = true
	if !*data.IsBreaking {
		response.Status = "break_ended"
	} else {
		response.Status = "break_started"
	}

	return response, nil
}
