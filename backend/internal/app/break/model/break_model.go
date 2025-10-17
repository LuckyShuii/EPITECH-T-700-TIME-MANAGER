package model

import (
	"time"
)

type BreakBase struct {
	UUID            string  `json:"uuid"`
	StartTime       string  `json:"start_time"`
	EndTime         *string `json:"end_time"`
	DurationMinutes int     `json:"duration_minutes"`
	// status is either "active" or "completed"
	Status string `json:"status"`
}

type BreakRead struct {
	BreakBase
	BreakUUID string `json:"break_uuid"`
}

type BreakReadAll struct {
	BreakBase
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BreakUpdate struct {
	WorkSessionUUID string `json:"work_session_uuid"`
	IsBreaking      *bool  `json:"is_breaking"`
}

type BreakCreate struct {
	WorkSessionID string `json:"work_session_id"`
	// status is either "active" or "completed"
	Status string `json:"status"`
}

type BreakUpdateResponse struct {
	Success   bool    `json:"success"`
	StartTime string  `json:"start_time"`
	EndTime   *string `json:"end_time"`
	Status    string  `json:"status"`
}
