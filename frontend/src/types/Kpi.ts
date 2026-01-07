// Types pour "Temps de pause moyenne individu"
export interface IndividualPauseData {
  average_break_time: number
  start_date: string
  end_date: string
}

// Types pour "Travail hebdomadaire individuel"
export interface WorkingTimeIndividualData {
  total_time: number
  start_date: string
  end_date: string
  user_uuid: string
}

// Type enrichi pour affichage avec comparaison
export interface WorkingTimeIndividualDisplay extends WorkingTimeIndividualData {
  previousTotal?: number
  difference?: number
}

// Types pour "Travail hebdomadaire par équipe"
export interface TeamMemberWorkingTime {
  first_name: string
  last_name: string
  total_time: number
  user_uuid: string
}

export interface WorkingTimeTeamData {
  total_time: number
  members: TeamMemberWorkingTime[]
  start_date: string
  end_date: string
  team_name: string
  team_uuid: string
}

// Types pour "Taux de présence"
export interface PresenceRateData {
  presence_rate: number
  weekly_time_done: number
  weekly_rate_expected: number
  user_uuid: string
  first_name: string
  last_name: string
}