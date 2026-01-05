import type {
  WorkingTimeData,
  WeeklyProgressData,
  ShiftAverageData,
  TeamWorkingTimeData,
  IndividualPauseData,
  PresenceRateData,
  ClockingTimeData,
  TeamPauseData,
  ShiftPauseData
} from '@/types/kpi'

export const mockWorkingTimeIndividual: WorkingTimeData = {
  totalWeek: 38.5,
  byDay: [
    { day: 'lundi', hours: 8.2 },
    { day: 'mardi', hours: 7.5 },
    { day: 'mercredi', hours: 8.0 },
    { day: 'jeudi', hours: 7.3 },
    { day: 'vendredi', hours: 7.5 }
  ]
}

export const mockWeeklyProgress: WeeklyProgressData = {
  weekStartDate: '2025-10-20',
  contractHours: 35,
  workedHours: 28.5,
  remainingHours: 6.5,
  percentageComplete: 81.4
}

export const mockShiftAverage: ShiftAverageData = {
  userId: 123,
  labels: ['S-4', 'S-3', 'S-2', 'S-1'],
  values: [7.5, 8.2, 7.8, 8.0]
}



// Mock pour les cas "pas de données"
export const mockEmptyWorkingTime: WorkingTimeData = {
  totalWeek: 0,
  byDay: []
}

export const mockTeamWorkingTime: TeamWorkingTimeData = {
  totalTeam: 156.0,
  members: [
    { userId: 1, name: 'Jean Dupont', hours: 38.5 },
    { userId: 2, name: 'Marie Martin', hours: 40.0 },
    { userId: 3, name: 'Paul Bernard', hours: 35.5 },
    { userId: 4, name: 'Sophie Dubois', hours: 42.0 }
  ]
}
// ... les anciennes mock data ...

export const mockIndividualPause: IndividualPauseData = {
  averagePausePerDay: 45,
  totalPauseWeek: 225,
  byDay: [
    { day: 'lundi', minutes: 50 },
    { day: 'mardi', minutes: 40 },
    { day: 'mercredi', minutes: 45 },
    { day: 'jeudi', minutes: 48 },
    { day: 'vendredi', minutes: 42 }
  ]
}

export const mockPresenceRate: PresenceRateData = {
  userId: 123,
  presenceRate: 95,
  daysPresent: 19,
  daysExpected: 20
}

export const mockClockingTime: ClockingTimeData = {
  labels: ['8h00', '8h10', '8h20', '8h30', '8h40', '8h50', '9h00'],
  values: [3, 5, 2, 8, 4, 1, 6]
}

export const mockTeamPause: TeamPauseData = {
  labels: ['Jean', 'Marie', 'Paul', 'Sophie'],
  values: [45, 38, 52, 42]
}

export const mockShiftPause: ShiftPauseData = {
  effectiveWork: 420,
  pause: 60
}

export const mockEmptyShiftAverage: ShiftAverageData | null = null