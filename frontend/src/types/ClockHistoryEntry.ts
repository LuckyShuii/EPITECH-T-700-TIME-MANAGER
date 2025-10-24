export interface ClockHistoryEntry {
  clock_in: string,
  clock_out: string | null,
  duration_minutes: number,
  breaks_duration_minutes: number,
  status: 'completed' | 'active' | 'paused',
  work_session_uuid: string,
  user_uuid: string,
  username: string
}