package model

import (
	"app/internal/app/common/model"
	"encoding/json"
	"fmt"
)

type TeamBase struct {
	UUID        string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();unique;not null"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
}

// swagger:model TeamMember
type TeamMember struct {
	UserUUID          string            `json:"user_uuid"`
	Roles             model.StringArray `json:"roles"`
	Status            string            `json:"status"`
	IsManager         bool              `json:"is_manager"`
	Username          string            `json:"username"`
	Email             string            `json:"email"`
	FirstName         string            `json:"first_name"`
	LastName          string            `json:"last_name"`
	PhoneNumber       string            `json:"phone_number"`
	WorkSessionStatus *string           `json:"work_session_status,omitempty"`
	WeeklyRate        int               `json:"weekly_rate"`
	WeeklyRateName    *string           `json:"weekly_rate_name,omitempty"`
	FirstDayOfWeek    *int              `json:"first_day_of_week,omitempty"`
}

type TeamMemberInfo struct {
	TeamUUID          string  `json:"team_uuid"`
	TeamName          string  `json:"team_name"`
	TeamDescription   *string `json:"team_description,omitempty"`
	IsManager         bool    `json:"is_manager"`
	Username          string  `json:"username"`
	Email             string  `json:"email"`
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	PhoneNumber       string  `json:"phone_number"`
	WorkSessionStatus *string `json:"work_session_status,omitempty"`
	FirstDayOfWeek    *int    `json:"first_day_of_week,omitempty"`
}

type TeamMembers []TeamMember

func (tm *TeamMembers) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for TeamMembers, got %T", value)
	}
	return json.Unmarshal(bytes, tm)
}

// swagger:model TeamReadAll
type TeamReadAll struct {
	TeamBase
	TeamMembers TeamMembers `json:"team_members" gorm:"type:jsonb"`
}

// swagger:model NewTeamMember
type NewTeamMember struct {
	UserUUID  string `json:"user_uuid" binding:"required,uuid"`
	IsManager bool   `json:"is_manager"`
}

// swagger:model TeamCreate
type TeamCreate struct {
	Name        string           `json:"name" binding:"required"`
	Description *string          `json:"description"`
	MemberUUIDs *[]NewTeamMember `json:"member_uuids" binding:"required,dive,required"`
}

type TeamMemberCreate struct {
	UserID    int  `json:"user_id"`
	IsManager bool `json:"is_manager"`
}

type TeamAddUsers struct {
	TeamUUID    string          `json:"team_uuid" binding:"required,uuid"`
	MemberUUIDs []NewTeamMember `json:"member_uuids" binding:"required,dive,required"`
}

// swagger:model TeamUpdate
type TeamUpdate struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// swagger:model TeamMemberLight
type TeamMemberLight struct {
	UserID    int    `json:"user_id"`
	UserUUID  string `json:"user_uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
