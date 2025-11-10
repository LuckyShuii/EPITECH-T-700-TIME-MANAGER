package repository

import (
	"math"

	"gorm.io/gorm"
)

type KPIRepository interface {
	GetWeeklyRatesByUserIDAndDateRange(userID int, startDate string, endDate string) (int, error)
	GetUserPresenceRate(userID int, startDate, endDate string) (float64, float64, float64, error)
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

func (repo *kpiRepository) GetUserPresenceRate(userID int, startDate, endDate string) (float64, float64, float64, error) {
	// Get the user's weekly rate
	var weeklyRateDB float64
	err := repo.db.Raw(`
		SELECT COALESCE(wr.amount, 40)
		FROM users u
		LEFT JOIN weekly_rate wr ON wr.id = u.weekly_rate_id
		WHERE u.id = ?
	`, userID).Scan(&weeklyRateDB).Error
	if err != nil {
		return 0, 0, 0, err
	}

	// Minutes to hours for weekly rate
	var totalMinutes int
	err = repo.db.Raw(`
		SELECT COALESCE(SUM(duration_minutes), 0)
		FROM (
			SELECT duration_minutes FROM work_session_active WHERE user_id = ? AND clock_in BETWEEN ? AND ?
			UNION ALL
			SELECT duration_minutes FROM work_session_archived WHERE user_id = ? AND clock_in BETWEEN ? AND ?
		) AS all_sessions
	`, userID, startDate, endDate, userID, startDate, endDate).Scan(&totalMinutes).Error
	if err != nil {
		return 0, 0, 0, err
	}

	// Convert total minutes to hours
	doneHours := float64(totalMinutes) / 60

	var presenceRate float64
	if weeklyRateDB > 0 {
		presenceRate = (doneHours / weeklyRateDB) * 100
	}

	// Round to 2 decimal places
	presenceRate = math.Round(presenceRate*100) / 100
	weeklyRateDB = math.Round(weeklyRateDB*100) / 100
	doneHours = math.Round(doneHours*100) / 100

	return presenceRate, weeklyRateDB, doneHours, nil
}
