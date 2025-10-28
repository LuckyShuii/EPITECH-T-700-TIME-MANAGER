package model

// swagger:model KPIWorkSessionUserWeeklyTotalResponse
type KPIWorkSessionUserWeeklyTotalResponse struct {
	TotalTime int    `json:"total_time"`
	UserUUID  string `json:"user_uuid"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
