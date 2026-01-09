import APIHandler from '../APIHandler'

export default {
  // Pause moyenne individuelle (average break time)
  getAverageBreakTime: (userUuid: string, startDate: string, endDate: string) =>
    APIHandler.get(`/kpi/average-break-time/${userUuid}/${startDate}/${endDate}`),

  // Moyenne par shift
  getAverageTimePerShift: (userUuid: string, startDate: string, endDate: string) =>
    APIHandler.get(`/kpi/average-time-per-shift/${userUuid}/${startDate}/${endDate}`),

  // Travail hebdomadaire individuel
  getWorkingTimeIndividual: (userUuid: string, startDate: string, endDate: string) =>
    APIHandler.get(`/kpi/work-session-user-weekly-total/${userUuid}/${startDate}/${endDate}`),

  // Travail hebdomadaire équipe - CORRIGÉ avec /kpi/
  getWorkingTimeTeam: (teamUuid: string, startDate: string, endDate: string) =>
    APIHandler.get(`/kpi/work-session-team-weekly-total/${teamUuid}/${startDate}/${endDate}`),

  // Taux de présence
  getPresenceRate: (userUuid: string, startDate: string, endDate: string) =>
    APIHandler.get(`/kpi/presence-rate/${userUuid}/${startDate}/${endDate}`)
}