package model

// swagger:model WeeklyRate
type WeeklyRate struct {
	UUID     *string `gorm:"not null" json:"uuid"`
	RateName string  `gorm:"not null" json:"rate_name"`
	Amount   float64 `gorm:"not null" json:"amount"`
}

// swagger:model CreateWeeklyRate
type CreateWeeklyRate struct {
	RateName string  `gorm:"not null" json:"rate_name"`
	Amount   float64 `gorm:"not null" json:"amount"`
}

// swagger:model UpdateWeeklyRate
type UpdateWeeklyRate struct {
	RateName string  `gorm:"not null" json:"rate_name"`
	Amount   float64 `gorm:"not null" json:"amount"`
}
