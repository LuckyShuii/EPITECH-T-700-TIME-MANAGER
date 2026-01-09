package handler

import (
	"fmt"
	"net/http"
	"time"

	AuthService "app/internal/app/auth/service"
	"app/internal/app/kpi/model"
	KPIService "app/internal/app/kpi/service"
	Config "app/internal/config"

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
		"2006-01-02",
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

func stripToDate(input string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05.999999", input)
	if err == nil {
		return t.Truncate(24 * time.Hour), nil
	}

	t, err = time.Parse("2006-01-02", input)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("invalid date format")
}

// GetWorkSessionUserWeeklyTotal handles the HTTP request to get the total work session time for a user within a date range.
//
// @Summary Get total work session time for a user within a date range
// @Description Retrieves the total work session time in minutes for a specified user UUID between the provided start and end dates. 🔒 Requires role: **any**
// @Tags KPI
// @Accept json
// @Security     BearerAuth
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

// GetWorkSessionTeamWeeklyTotal handles the HTTP request to get the total work session time for a team within a date range.
//
// @Summary Get total work session time for a team within a date range
// @Description Retrieves the total work session time in minutes for a specified team UUID between the provided start and end dates. 🔒 Requires role: **manager**
// @Tags KPI
// @Accept json
// @Security     BearerAuth
// @Produce json
// @Param start_date path string true "Start Date in ISO 8601 format"
// @Param end_date path string true "End Date in ISO 8601 format"
// @Param team_uuid path string true "Team UUID"
// @Success 200 {object} model.KPIWorkSessionTeamWeeklyTotalResponse
// @Router /kpi/work-session-team-weekly-total/{team_uuid}/{start_date}/{end_date} [get]
func (handler *KPIHandler) GetWorkSessionTeamWeeklyTotal(c *gin.Context) {
	startDate := c.Param("start_date")
	endDate := c.Param("end_date")
	teamUUID := c.Param("team_uuid")

	err := handler.validateDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range: " + err.Error()})
		return
	}

	weeklyRates, err := handler.service.GetWorkSessionTeamWeeklyTotal(startDate, endDate, teamUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve weekly rates: " + err.Error()})
		return
	}

	kpiResponse := model.KPIWorkSessionTeamWeeklyTotalResponse{
		TotalTime: weeklyRates.TotalTime,
		StartDate: startDate,
		EndDate:   endDate,
		TeamUUID:  teamUUID,
		Members:   weeklyRates.Members,
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

	return nil
}

// GetPresenceRate handles the HTTP request to get the presence rate for a user within a date range.
//
// @Summary Get presence rate for a user within a date range
// @Description Retrieves the presence rate percentage for a specified user UUID between the provided start and end dates. 🔒 Requires role: **any**
// @Tags KPI
// @Accept json
// @Security     BearerAuth
// @Produce json
// @Param start_date path string true "Start Date in ISO 8601 format"
// @Param end_date path string true "End Date in ISO 8601 format"
// @Param user_uuid path string true "User UUID"
// @Success 200 {object} model.KPIPresenceRateResponse
// @Router /kpi/presence-rate/{user_uuid}/{start_date}/{end_date} [get]
func (handler *KPIHandler) GetPresenceRate(c *gin.Context) {
	startDate := c.Param("start_date")
	endDate := c.Param("end_date")
	userUUID := c.Param("user_uuid")

	err := handler.validateDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date ranges: " + err.Error()})
		return
	}

	response, err := handler.service.GetPresenceRate(startDate, endDate, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve presence rate: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ExportKPIData handles the HTTP request to export KPI data within a date range.
//
// @Summary Export KPI data within a date range
// @Description Exports KPI data for the specified date range. 🔒 Requires role: **manager, admin**
// @Tags KPI
// @Accept json
// @Security     BearerAuth
// @Produce json
// @Param kpi_export_request body model.KPIExportRequest true "KPI Export Request"
// @Success 200 {object} model.KPIExportResponse
// @Router /kpi/export [post]
func (handler *KPIHandler) ExportKPIData(c *gin.Context) {
	var exportRequest model.KPIExportRequest
	if err := c.ShouldBindJSON(&exportRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": Config.ErrorMessages()["NO_CLAIMS"]})
		return
	}

	authClaims := claims.(*AuthService.Claims)

	err := handler.validateDateRange(exportRequest.StartDate, exportRequest.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date ranges: " + err.Error()})
		return
	}

	exportResponse, err := handler.service.ExportKPIData(exportRequest.StartDate, exportRequest.EndDate, authClaims.UUID, exportRequest.KPIType, exportRequest.UUIDToSearch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export KPI data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, exportResponse)
}

// GetAverageBreakTime handles the HTTP request to get the average break time for a user within a date range.
//
// @Summary Get average break time for a user within a date range
// @Description Retrieves the average break time in minutes for a specified user UUID between the provided start and end dates. These dates must be 5 days! (work days) Not less, not more 🔒 Requires role: **manager, admin**
// @Tags KPI
// @Accept json
// @Security     BearerAuth
// @Produce json
// @Param start_date path string true "Start Date in ISO 8601 format"
// @Param end_date path string true "End Date in ISO 8601 format"
// @Param user_uuid path string true "User UUID"
// @Success 200 {object} model.KPIAverageBreakTimeResponse
// @Router /kpi/average-break-time/{user_uuid}/{start_date}/{end_date} [get]
func (handler *KPIHandler) GetAverageBreakTime(c *gin.Context) {
	startDate := c.Param("start_date")
	endDate := c.Param("end_date")
	userUUID := c.Param("user_uuid")

	err := handler.validateDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range: " + err.Error()})
		return
	}

	start, err := stripToDate(startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
		return
	}

	end, err := stripToDate(endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
		return
	}

	diffDays := int(end.Sub(start).Hours() / 24)
	if diffDays != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date range must be exactly 5 days apart"})
		return
	}

	averageBreakTime, err := handler.service.GetAverageBreakTime(startDate, endDate, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve average break time"})
		return
	}

	kpiResponse := model.KPIAverageBreakTimeResponse{
		AverageBreakTime: averageBreakTime.AverageBreakTime,
		StartDate:        startDate,
		EndDate:          endDate,
		UserUUID:         userUUID,
		FirstName:        averageBreakTime.FirstName,
		LastName:         averageBreakTime.LastName,
	}

	c.JSON(http.StatusOK, kpiResponse)
}

// GetAverageTimePerShift handles the HTTP request to get the average time per shift for a user within a date range.
//
// @Summary Get average time per shift for a user within a date range
// @Description Retrieves the average time per shift in minutes for a specified user UUID between the provided start and end dates. 🔒 Requires role: **manager, admin**
// @Tags KPI
// @Accept json
// @Security     BearerAuth
// @Produce json
// @Param start_date path string true "Start Date in ISO 8601 format"
// @Param end_date path string true "End Date in ISO 8601 format"
// @Param user_uuid path string true "User UUID"
// @Success 200 {object} model.KPIAverageTimePerShiftResponse
// @Router /kpi/average-time-per-shift/{user_uuid}/{start_date}/{end_date} [get]
func (handler *KPIHandler) GetAverageTimePerShift(c *gin.Context) {
	startDate := c.Param("start_date")
	endDate := c.Param("end_date")
	userUUID := c.Param("user_uuid")

	err := handler.validateDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range: " + err.Error()})
		return
	}

	averageTimePerShift, err := handler.service.GetAverageTimePerShift(startDate, endDate, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve average time per shift: " + err.Error()})
		return
	}

	kpiResponse := model.KPIAverageTimePerShiftResponse{
		AverageTimePerShift: averageTimePerShift.AverageTimePerShift,
		TotalShifts:         averageTimePerShift.TotalShifts,
		TotalTime:           averageTimePerShift.TotalTime,
		StartDate:           startDate,
		EndDate:             endDate,
		UserUUID:            userUUID,
		FirstName:           averageTimePerShift.FirstName,
		LastName:            averageTimePerShift.LastName,
	}

	c.JSON(http.StatusOK, kpiResponse)
}
