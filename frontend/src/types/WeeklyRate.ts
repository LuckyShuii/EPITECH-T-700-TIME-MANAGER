export interface WeeklyRate {
  uuid: string
  rate_name: string
  amount: number
}

export interface WeeklyRateCreateData {
  rate_name: string
  amount: number
}

export interface WeeklyRateUpdateData {
  rate_name: string
  amount: number
}
