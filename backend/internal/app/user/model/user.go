package model

import (
	"app/internal/app/common/model"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type UserBase struct {
	UUID           string            `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();unique;not null"`
	Username       string            `json:"username" gorm:"unique;not null"`
	Email          string            `json:"email" gorm:"unique;not null"`
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	PhoneNumber    *string           `json:"phone_number,omitempty"`
	Roles          model.StringArray `json:"roles" gorm:"type:text[];default:'{employee}'"`
	FirstDayOfWeek *int              `json:"first_day_of_week,omitempty" gorm:"default:1"`
}

type JSONLayout []map[string]any

func (j *JSONLayout) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONLayout: %T", value)
	}

	// First try to unmarshal as an array
	if err := json.Unmarshal(data, j); err == nil {
		return nil
	}

	// If it's not an array, try to wrap it in an array with a single element
	var single map[string]any
	if err := json.Unmarshal(data, &single); err == nil {
		*j = []map[string]any{single}
		return nil
	}

	return fmt.Errorf("failed to parse JSONLayout: %s", string(data))
}

func (j JSONLayout) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

type DashboardLayout struct {
	I      string `json:"i"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
	MinW   int    `json:"minW"`
	MinH   int    `json:"minH"`
	Static bool   `json:"static"`
}

// swagger:model UserDashboardLayout
type UserDashboardLayout struct {
	DashboardLayout JSONLayout `json:"dashboard_layout" gorm:"type:jsonb;default:'[]'"`
}

// swagger:model DashboardLayoutResponse
type DashboardLayoutResponse struct {
	Layout []DashboardLayout `json:"layout"`
}

// swagger:model UserDashboardLayoutUpdate
type UserDashboardLayoutUpdate struct {
	Layout []DashboardLayout `json:"layout"`
}

// swagger:model UserTeamMemberInfo
type UserTeamMemberInfo struct {
	TeamUUID        string  `json:"team_uuid" example:"4bc3df44-491c-4073-9e89-682bb0acfca0"`
	TeamName        string  `json:"team_name" example:"Développement"`
	TeamDescription *string `json:"team_description,omitempty" example:"Équipe principale backend"`
	IsManager       bool    `json:"is_manager" example:"true"`
}

// swagger:model UserReadAll
type UserReadAll struct {
	UserBase
	Status            *string              `json:"status"`
	WorkSessionStatus *string              `json:"work_session_status"`
	TeamsRaw          string               `json:"-" gorm:"column:teams"`
	Teams             []UserTeamMemberInfo `json:"teams" gorm:"-"`
	WeeklyRate        int                  `json:"weekly_rate"`
	WeeklyRateName    *string              `json:"weekly_rate_name"`
}

type UserUpdateEntry struct {
	UUID           string             `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();unique;not null"`
	Username       *string            `json:"username,omitempty" gorm:"unique;not null"`
	Email          *string            `json:"email,omitempty" gorm:"unique;not null"`
	FirstName      *string            `json:"first_name,omitempty"`
	LastName       *string            `json:"last_name,omitempty"`
	PhoneNumber    *string            `json:"phone_number,omitempty"`
	Roles          *model.StringArray `json:"roles,omitempty" gorm:"type:text[];default:'{employee}'"`
	Status         *string            `json:"status,omitempty"`
	WeeklyRateUUID *string            `json:"weekly_rate_uuid,omitempty"`
	WeeklyRateID   *int               `json:"weekly_rate_id,omitempty"`
	FirstDayOfWeek *int               `json:"first_day_of_week,omitempty"`
}

// StringArray is a custom type representing an array of strings.
//
// swagger:model
type UserMeJWT struct {
	UserUUID    string            `json:"user_uuid"`
	Roles       model.StringArray `json:"roles"`
	Email       string            `json:"email"`
	Username    string            `json:"username"`
	FirstName   string            `json:"first_name"`
	LastName    string            `json:"last_name"`
	PhoneNumber *string           `json:"phone_number,omitempty"`
}

type UserAll struct {
	UserBase
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PasswordHash string    `json:"password_hash" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// UserRead represents the user data returned in responses.
//
// swagger:model
type UserRead struct {
	UserBase
	Status            *string   `json:"status"`
	WorkSessionStatus *string   `json:"work_session_status"`
	WeeklyRate        *int      `json:"weekly_rate"`
	WeeklyRateName    *string   `json:"weekly_rate_name"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type UserReadJWT struct {
	UserBase
	ID           *uint  `json:"id" gorm:"primaryKey;autoIncrement"`
	PasswordHash string `json:"password_hash" gorm:"not null"`
}

type UserUpdate struct {
	UserBase
	PasswordHash *string `json:"password_hash,omitempty"`
}

type UserDelete struct {
	UUID string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();unique;not null"`
}

// UserCreate represents the payload for creating a new user.
//
// swagger:model
type UserCreate struct {
	UserBase
	WeeklyRateUUID *string `json:"weekly_rate_uuid,omitempty"`
	WeeklyRateID   *int    `json:"weekly_rate_id,omitempty"`
	PasswordHash   string  `json:"password_hash" gorm:"not null"`
	Password       string  `json:"password" gorm:"-:all"` // Ignored by GORM, used only for input
}

type UserLogin struct {
	Username *string `json:"username" gorm:"unique;not null"`
	Email    *string `json:"email" gorm:"unique;not null"`
	Password string  `json:"password" gorm:"not null"`
}

// UserUUIDPayload represents the payload for identifying a user by UUID.
//
// swagger:model
type UserUUIDPayload struct {
	UserUUID string `json:"user_uuid" example:"e1234abc-5678-90de-f123-4567890abcde"`
}

// UserStatusUpdatePayload represents the payload for updating a user's status.
//
// swagger:model
type UserStatusUpdatePayload struct {
	UserUUID string `json:"user_uuid" example:"e1234abc-5678-90de-f123-4567890abcde"`
	Status   string `json:"status" example:"active"`
}
