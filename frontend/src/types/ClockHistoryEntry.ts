export interface ClockHistoryEntry {
  clockInTime: string
  clockOutTime: string | null
  totalHours: number
  status: 'active' | 'completed' | 'paused'
}