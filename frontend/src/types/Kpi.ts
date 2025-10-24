// Types pour "Travail hebdomadaire individuel"
export interface DayWork {
  day: string
  hours: number
}

export interface WorkingTimeData {
  totalWeek: number
  byDay: DayWork[]
}

// Types pour "Travail hebdomadaire par équipe"
export interface TeamMember {
  userId: number
  name: string
  hours: number
}

export interface TeamWorkingTimeData {
  totalTeam: number
  members: TeamMember[]
}

// Types pour "Moyenne par shift par individu"
export interface ShiftAverageData {
  userId: number
  labels: string[]
  values: number[]
}

// Types pour "Taux de présence"
export interface PresenceRateData {
  userId: number
  presenceRate: number
  daysPresent: number
  daysExpected: number
}

// Types pour "Heure de pointage moyen hebdo"
export interface ClockingTimeData {
  labels: string[]
  values: number[]
}

// Types pour "Temps de pause moyenne individu"
export interface DayPause {
  day: string
  minutes: number
}

export interface IndividualPauseData {
  averagePausePerDay: number
  totalPauseWeek: number
  byDay: DayPause[]
}

// Types pour "Temps de pause moyenne équipe"
export interface TeamPauseData {
  labels: string[]
  values: number[]
}

// Types pour "Progression hebdomadaire"
export interface WeeklyProgressData {
  weekStartDate: string // Format ISO: '2025-10-21'
  contractHours: number
  workedHours: number
  remainingHours: number
  percentageComplete: number
}

// Types pour "Temps de pause par shift"
export interface ShiftPauseData {
  effectiveWork: number
  pause: number
}