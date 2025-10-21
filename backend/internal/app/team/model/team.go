package model

import (
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
	UserUUID  string   `json:"user_uuid"`
	Roles     []string `json:"roles"`
	Status    string   `json:"status"`
	IsManager bool     `json:"is_manager"`
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
	TeamMembers TeamMembers `json:"team_members" gorm:"column:team_members"`
}

// swagger:model NewTeamMember
type NewTeamMember struct {
	UserUUID  string `json:"user_uuid" binding:"required,uuid4"`
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
