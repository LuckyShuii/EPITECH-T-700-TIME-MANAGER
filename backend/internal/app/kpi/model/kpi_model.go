package model

// swagger:model KPIWorkSessionUserWeeklyTotalResponse
type KPIWorkSessionUserWeeklyTotalResponse struct {
	TotalTime int    `json:"total_time"`
	UserUUID  string `json:"user_uuid"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// swagger:model KPIWorkSessionTeamWeeklyTotalResponse
type KPIWorkSessionTeamWeeklyTotalResponse struct {
	TotalTime int                                   `json:"total_time"`
	TeamUUID  string                                `json:"team_uuid"`
	TeamName  string                                `json:"team_name"`
	StartDate string                                `json:"start_date"`
	EndDate   string                                `json:"end_date"`
	Members   []KPIWorkSessionTeamMemberWeeklyTotal `json:"members"`
}

// swagger:model KPIWorkSessionTeamMemberWeeklyTotal
type KPIWorkSessionTeamMemberWeeklyTotal struct {
	UserUUID  string `json:"user_uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	TotalTime int    `json:"total_time"`
}

// swagger:model MemberWeeklyRate
type MemberWeeklyRate struct {
	UserUUID  string `json:"user_uuid"`
	TotalTime int    `json:"total_time"`
}

// swagger:model KPIPresenceRateResponse
type KPIPresenceRateResponse struct {
	FirstName          string  `json:"first_name"`
	LastName           string  `json:"last_name"`
	UserUUID           string  `json:"user_uuid"`
	PresenceRate       float64 `json:"presence_rate"`
	WeeklyRateExpected float64 `json:"weekly_rate_expected"`
	WeeklyTimeDone     float64 `json:"weekly_time_done"`
}

// swagger:model KPIExportRequest
type KPIExportRequest struct {
	KPIType      string `json:"kpi_type" binding:"required,oneof=work_session_user_weekly_total work_session_team_weekly_total presence_rate weekly_average_break_time average_time_per_shift"`
	StartDate    string `json:"start_date" binding:"required"`
	EndDate      string `json:"end_date" binding:"required"`
	UUIDToSearch string `json:"uuid_to_search"`
}

// swagger:model KPIExportResponse
type KPIExportResponse struct {
	File string `json:"file"`
	URL  string `json:"url"`
}

// swagger:model KPIAverageBreakTimeResponse
type KPIAverageBreakTimeResponse struct {
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	UserUUID         string  `json:"user_uuid"`
	AverageBreakTime float64 `json:"average_break_time"`
	StartDate        string  `json:"start_date"`
	EndDate          string  `json:"end_date"`
}

// swagger:model KPIAverageTimePerShiftResponse
type KPIAverageTimePerShiftResponse struct {
	FirstName           string  `json:"first_name"`
	LastName            string  `json:"last_name"`
	UserUUID            string  `json:"user_uuid"`
	AverageTimePerShift float64 `json:"average_time_per_shift"`
	TotalShifts         int     `json:"total_shifts"`
	TotalTime           int     `json:"total_time"`
	StartDate           string  `json:"start_date"`
	EndDate             string  `json:"end_date"`
}
