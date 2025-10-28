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
	TotalTime int    `json:"total_time"`
	TeamUUID  string `json:"team_uuid"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
