import type {
  WorkingTimeData,
  WeeklyProgressData,
  ShiftAverageData
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

export const mockEmptyWeeklyProgress: WeeklyProgressData | null = null

export const mockEmptyShiftAverage: ShiftAverageData | null = null