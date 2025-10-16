package model

import (
	"time"

	"app/internal/app/user/model"

	"github.com/lib/pq"
)

type WorkSessionBase struct {
	UUID            string `json:"uuid"`
	ClockIn         string `json:"clock_in"`
	ClockOut        string `json:"clock_out"`
	DurationMinutes int    `json:"duration_minutes"`
	// status is either "active", "paused" or "completed"
	Status string `json:"status"`
}

type WorkSessionRead struct {
	WorkSessionBase

	WorkSessionUUID string `json:"work_session_uuid"`

	// User fields
	UserUUID    string         `json:"user_uuid"`
	Username    string         `json:"username"`
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	PhoneNumber *string        `json:"phone_number"`
	Roles       pq.StringArray `json:"roles" gorm:"type:text[]"`
}

type WorkSessionReadAll struct {
	WorkSessionBase
	User      model.UserBase `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type WorkSessionUpdate struct {
	UserUUID  string `json:"user_uuid"`
	IsClocked *bool  `json:"is_clocked"`
}

type WorkSessionCreate struct {
	UserID int `json:"user_id"`
	// status is either "active", "paused" or "completed"
	Status string `json:"status"`
}

type WorkSessionUpdateResponse struct {
	Success      bool    `json:"success"`
	ClockInTime  string  `json:"clock_in_time"`
	ClockOutTime *string `json:"clock_out_time"`
	Status       string  `json:"status"`
}
