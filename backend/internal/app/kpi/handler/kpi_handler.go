package handler

import (
	"fmt"
	"net/http"
	"time"

	"app/internal/app/kpi/model"
	KPIService "app/internal/app/kpi/service"

	"github.com/gin-gonic/gin"
)

type KPIHandler struct {
	service KPIService.KPIService
}

func NewKPIHandler(service KPIService.KPIService) *KPIHandler {
	return &KPIHandler{service: service}
}

func (handler *KPIHandler) isValidISO8601(date string) bool {
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.999999",
	}

	for _, layout := range layouts {
		if _, err := time.Parse(layout, date); err == nil {
			return true
		}
	}
	return false
}

// GetWorkSessionUserWeeklyTotal handles the HTTP request to get the total work session time for a user within a date range.
//
// @Summary Get total work session time for a user within a date range
// @Description Retrieves the total work session time in minutes for a specified user UUID between the provided start and end dates. 🔒 Requires role: **any**
// @Tags KPI
// @Accept json
// @Produce json
// @Param start_date path string true "Start Date in ISO 8601 format"
// @Param end_date path string true "End Date in ISO 8601 format"
// @Param user_uuid path string true "User UUID"
// @Success 200 {object} model.KPIWorkSessionUserWeeklyTotalResponse
// @Router /kpi/work-session-user-weekly-total/{user_uuid}/{start_date}/{end_date} [get]
func (handler *KPIHandler) GetWorkSessionUserWeeklyTotal(c *gin.Context) {
	startDate := c.Param("start_date")
	endDate := c.Param("end_date")
	userUUID := c.Param("user_uuid")

	err := handler.validateDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range: " + err.Error()})
		return
	}

	weeklyRates, err := handler.service.GetWorkSessionUserWeeklyTotal(startDate, endDate, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve weekly rates"})
		return
	}

	kpiResponse := model.KPIWorkSessionUserWeeklyTotalResponse{
		TotalTime: weeklyRates,
		StartDate: startDate,
		EndDate:   endDate,
		UserUUID:  userUUID,
	}

	c.JSON(http.StatusOK, kpiResponse)
}

func (handler *KPIHandler) validateDateRange(startDate string, endDate string) error {
	if !handler.isValidISO8601(startDate) || !handler.isValidISO8601(endDate) {
		return &gin.Error{
			Err:  http.ErrNotSupported,
			Type: gin.ErrorTypeBind,
		}
	}

	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)
	now := time.Now().Format(time.RFC3339Nano)
	twoYearsAgo := time.Now().AddDate(-2, 0, 0).Format(time.RFC3339Nano)
	threeDaysBeforeFromNow := time.Now().AddDate(0, 0, -3).Format(time.RFC3339Nano)

	// check if start date is before end date
	if start.After(end) {
		return &gin.Error{
			Err:  http.ErrNotSupported,
			Type: gin.ErrorTypeBind,
		}
	}

	// check if start date is after end date
	if startDate > endDate {
		return fmt.Errorf("start_date cannot be after end_date")
	}

	// check if start date is before end date
	if endDate < startDate {
		return fmt.Errorf("end_date cannot be before start_date")
	}

	// check if end date is in the future
	if startDate < twoYearsAgo {
		return fmt.Errorf("date range cannot exceed 2 years from the current date")
	}

	// check if end date is in the future
	if endDate > now {
		return fmt.Errorf("end_date cannot be in the future")
	}

	// check if end date is minimum 3 days from today
	if endDate > threeDaysBeforeFromNow {
		return fmt.Errorf("end_date cannot be less than 3 days from today")
	}

	return nil
}
