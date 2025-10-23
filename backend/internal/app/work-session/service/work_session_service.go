package service

import (
	WorkSessionModel "app/internal/app/work-session/model"
	"app/internal/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	UserService "app/internal/app/user/service"
	WorkSessionRepository "app/internal/app/work-session/repository"

	BreakRepository "app/internal/app/break/repository"

	"github.com/google/uuid"
)

type WorkSessionService interface {
	UpdateWorkSessionClocking(data WorkSessionModel.WorkSessionUpdate) (WorkSessionModel.WorkSessionUpdateResponse, error)
	GetWorkSessionStatus(userUUID string) (WorkSessionModel.WorkSessionStatus, error)
	GetWorkSessionHistory(userUUID string, startDate string, endDate string, limit int, offset int) ([]WorkSessionModel.WorkSessionReadHistory, error)
}

type workSessionService struct {
	WorkSessionRepo WorkSessionRepository.WorkSessionRepository
	UserService     UserService.UserService
	BreakRepository BreakRepository.BreakRepository
}

func NewWorkSessionService(repo WorkSessionRepository.WorkSessionRepository, userService UserService.UserService, breakRepo BreakRepository.BreakRepository) WorkSessionService {
	return &workSessionService{WorkSessionRepo: repo, UserService: userService, BreakRepository: breakRepo}
}

func (service *workSessionService) UpdateWorkSessionClocking(data WorkSessionModel.WorkSessionUpdate) (WorkSessionModel.WorkSessionUpdateResponse, error) {
	var response WorkSessionModel.WorkSessionUpdateResponse

	// 1️⃣ Get user ID from UUID
	userID, userErr := service.UserService.GetIdByUuid(data.UserUUID)
	if userErr != nil {
		response.Success = false
		return response, userErr
	}

	// 2️⃣ Check for active work session
	workSessionFound, err := service.WorkSessionRepo.GetUserActiveWorkSession(userID, []string{"active", "paused"})
	if err != nil {
		response.Success = false
		return response, err
	}

	// 3️⃣ Clock-in while already clocked-in → error
	if workSessionFound.WorkSessionUUID != "" && *data.IsClocked {
		return response, fmt.Errorf("an active work session already exists for this user, cannot clock in again")
	}

	// 4️⃣ No active session but clock-in → create session
	if workSessionFound.WorkSessionUUID == "" && *data.IsClocked {
		now := time.Now().In(time.FixedZone("Europe/Paris", 2*60*60))
		response.ClockInTime = now.Format(time.RFC3339Nano)
		response.Status = "clocked_in"
		response.Success = true

		err := service.WorkSessionRepo.CreateWorkSession(uuid.New().String(), userID, "active")
		if err != nil {
			response.Success = false
			return response, err
		}
		return response, nil
	}

	// 5️⃣ No active session and clock-out → error
	if workSessionFound.WorkSessionUUID == "" && !*data.IsClocked {
		response.Success = false
		return response, fmt.Errorf("no active work session found for this user, cannot clock out")
	}

	// 6️⃣ Clock-out process
	if workSessionFound.WorkSessionUUID != "" && !*data.IsClocked {
		return service.completeWorkSessionProcess(workSessionFound, userID)
	}

	response.Success = true
	response.Status = "clocked_in"
	return response, nil
}

func (service *workSessionService) GetWorkSessionStatus(userUUID string) (WorkSessionModel.WorkSessionStatus, error) {
	response := WorkSessionModel.WorkSessionStatus{}

	// 1️⃣ Get user ID from UUID
	userID, userErr := service.UserService.GetIdByUuid(userUUID)
	if userErr != nil {
		return WorkSessionModel.WorkSessionStatus{}, userErr
	}

	// 2️⃣ Check for active work session
	workSessionFound, err := service.WorkSessionRepo.GetUserActiveWorkSession(userID, []string{"active", "paused"})
	if err != nil {
		return WorkSessionModel.WorkSessionStatus{}, err
	}

	if workSessionFound.WorkSessionUUID != "" {
		response.WorkSessionUUID = workSessionFound.WorkSessionUUID
		response.IsClocked = true
		response.ClockInTime = &workSessionFound.ClockIn
		response.Status = workSessionFound.Status
	} else {
		response.IsClocked = false
		response.Status = "no_active_session"
	}

	return response, nil
}

func (service *workSessionService) GetWorkSessionHistory(userUUID string, startDate string, endDate string, limit int, offset int) ([]WorkSessionModel.WorkSessionReadHistory, error) {
	ctx := context.Background()

	cacheKey := fmt.Sprintf("worksession:history:%s:%s:%s:%d:%d", userUUID, startDate, endDate, limit, offset)

	cached, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var cachedHistory []WorkSessionModel.WorkSessionReadHistory
		if jsonErr := json.Unmarshal([]byte(cached), &cachedHistory); jsonErr == nil {
			return cachedHistory, nil
		}
	}

	userID, userErr := service.UserService.GetIdByUuid(userUUID)
	if userErr != nil {
		return []WorkSessionModel.WorkSessionReadHistory{}, userErr
	}

	workSessions, err := service.WorkSessionRepo.GetWorkSessionHistoryByUserId(userID, startDate, endDate, limit, offset)
	if err != nil {
		return []WorkSessionModel.WorkSessionReadHistory{}, err
	}

	data, _ := json.Marshal(workSessions)
	if setErr := db.RedisClient.Set(ctx, cacheKey, data, 30*time.Second).Err(); setErr != nil {
		log.Printf("⚠️ Error setting cache for %s : %v", cacheKey, setErr)
	}

	return workSessions, nil
}

func (service *workSessionService) completeWorkSessionProcess(workSessionFound WorkSessionModel.WorkSessionRead, userID int) (WorkSessionModel.WorkSessionUpdateResponse, error) {
	var response WorkSessionModel.WorkSessionUpdateResponse

	loc, _ := time.LoadLocation("Europe/Paris")
	t2 := time.Now().In(loc)

	// Parse the clock-in time
	t1, err := time.Parse(time.RFC3339Nano, workSessionFound.ClockIn)
	if err != nil {
		log.Println("parse error:", err)
		response.Success = false
		return response, err
	}

	// Get the duration in minutes
	duration := t2.Sub(t1)
	minutes := math.Floor(duration.Minutes() + 0.5)

	// Update the work session
	err = service.WorkSessionRepo.CompleteWorkSession(workSessionFound.WorkSessionUUID, userID, int(minutes))
	if err != nil {
		response.Success = false
		return response, err
	}

	// Get the internal ID of the work session
	workSessionId, err := service.WorkSessionRepo.FindIdByUuid(workSessionFound.WorkSessionUUID)
	if err != nil {
		response.Success = false
		return response, err
	}

	// Get total break duration
	breakDuration, err := service.BreakRepository.GetTotalBreakDurationByWorkSessionId(workSessionId)
	if err != nil {
		response.Success = false
		return response, err
	}

	// Update break duration in work session
	err = service.WorkSessionRepo.UpdateBreakDurationMinutes(workSessionFound.WorkSessionUUID, breakDuration)
	if err != nil {
		response.Success = false
		return response, err
	}

	// Delete related breaks
	err = service.BreakRepository.DeleteRelatedBreaksToWorkSession(workSessionId)
	if err != nil {
		response.Success = false
		return response, err
	}

	// Prepare response
	response.ClockInTime = workSessionFound.ClockIn
	formattedTime := t2.Format(time.RFC3339Nano)
	response.ClockOutTime = &formattedTime
	response.Status = "clocked_out"
	response.Success = true

	return response, nil
}
