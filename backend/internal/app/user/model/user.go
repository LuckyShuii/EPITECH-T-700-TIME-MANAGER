package model

import (
	"time"

	"github.com/lib/pq"
)

type UserBase struct {
	UUID        string         `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();unique;not null"`
	Username    string         `json:"username" gorm:"unique;not null"`
	Email       string         `json:"email" gorm:"unique;not null"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	PhoneNumber *string        `json:"phone_number,omitempty"`
	Roles       pq.StringArray `json:"roles" gorm:"type:text[];default:'{employee}'"`
}

type UserAll struct {
	UserBase
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PasswordHash string    `json:"password_hash" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type UserRead struct {
	UserBase
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
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

type UserCreate struct {
	UserBase
	PasswordHash string `json:"password_hash" gorm:"not null"`
	Password     string `json:"password" gorm:"-:all"` // Ignored by GORM, used only for input
}

type UserLogin struct {
	Username *string `json:"username" gorm:"unique;not null"`
	Email    *string `json:"email" gorm:"unique;not null"`
	Password string  `json:"password" gorm:"not null"`
}
