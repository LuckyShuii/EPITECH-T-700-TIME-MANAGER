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
	StartDate string                                `json:"start_date"`
	EndDate   string                                `json:"end_date"`
	Members   []KPIWorkSessionTeamMemberWeeklyTotal `json:"members"`
}

// swagger:model KPIWorkSessionTeamMemberWeeklyTotal
type KPIWorkSessionTeamMemberWeeklyTotal struct {
	UserUUID  string `json:"user_uuid"`
	TotalTime int    `json:"total_time"`
}

// swagger:model MemberWeeklyRate
type MemberWeeklyRate struct {
	UserUUID  string `json:"user_uuid"`
	TotalTime int    `json:"total_time"`
}

// swagger:model KPIPresenceRateResponse
type KPIPresenceRateResponse struct {
	UserUUID           string  `json:"user_uuid"`
	PresenceRate       float64 `json:"presence_rate"`
	WeeklyRateExpected float64 `json:"weekly_rate_expected"`
	WeeklyTimeDone     float64 `json:"weekly_time_done"`
}
