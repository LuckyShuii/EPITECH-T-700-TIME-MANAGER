export interface ClockData {

    clockInTime: string | null, 
    clockOutTime: string | null,
    status: 'NOT_CLOCKED' | 'CLOCKED_IN' | 'CLOCKED_OUT'
}
