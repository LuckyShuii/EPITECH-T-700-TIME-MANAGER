package repository

import (
	"gorm.io/gorm"
)

type KPIRepository interface {
	GetWeeklyRatesByUserIDAndDateRange(userID int, startDate string, endDate string) (int, error)
}

type kpiRepository struct {
	db *gorm.DB
}

func NewKPIRepository(db *gorm.DB) KPIRepository {
	return &kpiRepository{db}
}

func (repo *kpiRepository) GetWeeklyRatesByUserIDAndDateRange(userID int, startDate string, endDate string) (int, error) {
	var weeklyRates struct {
		TotalDurationMinutes int `gorm:"column:total_duration_minutes"`
	}
	err := repo.db.Raw(
		`SELECT
			SUM(ws.duration_minutes) AS total_duration_minutes
		FROM users AS u
		INNER JOIN (
			SELECT
				user_id,
				clock_in,
				duration_minutes
			FROM work_session_active
			UNION ALL
			SELECT
				user_id,
				clock_in,
				duration_minutes
			FROM work_session_archived
		) AS ws(user_id, clock_in, duration_minutes)
			ON u.id = ws.user_id
		WHERE
			u.id = ?
			AND ws.clock_in BETWEEN ? AND ?`,
		userID, startDate, endDate,
	).Scan(&weeklyRates).Error

	if err != nil {
		return 0, err
	}

	return weeklyRates.TotalDurationMinutes, nil
}
