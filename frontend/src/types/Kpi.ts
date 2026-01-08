export interface PresenceRateUser {
  user_uuid: string
  first_name: string
  last_name: string
  presence_rate: number
  weekly_rate_expected: number
  weekly_time_done: number
}

export type PresenceRateData = PresenceRateUser[]

export interface IndividualPauseData {
  average_break_time: number
  total_breaks: number
  user_uuid: string
}

export interface WorkingTimeIndividualData {
  user_uuid: string
  total_time: number
  start_date: string
  end_date: string
}

export interface WorkingTimeIndividualDisplay extends WorkingTimeIndividualData {
  previousTotal?: number
  difference?: number
}

export interface WorkingTimeTeamMember {
  user_uuid: string
  first_name: string
  last_name: string
  total_time: number
}

export interface WorkingTimeTeamData {
  team_uuid: string
  team_name: string
  total_time: number
  start_date: string
  end_date: string
  members: WorkingTimeTeamMember[]
  previousTotal?: number
  difference?: number
}

export interface AverageTimePerShiftData {
  user_uuid: string
  average_time: number
  shift_count: number
  date_range: {
    start_date: string
    end_date: string
  }
}