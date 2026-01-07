import APIHandler from '../APIHandler'

const resource = 'kpi'

export default {
  // Temps de pause moyenne individu
  getAverageBreakTime(userUuid: string, startDate: string, endDate: string) {
    return APIHandler.get(`${resource}/average-break-time/${userUuid}/${startDate}/${endDate}`)
  },

  // Travail hebdomadaire individuel
  getWorkingTimeIndividual(userUuid: string, startDate: string, endDate: string) {
    return APIHandler.get(`${resource}/work-session-user-weekly-total/${userUuid}/${startDate}/${endDate}`)
  },

  // Travail hebdomadaire par équipe
  getWorkingTimeTeam(teamUuid: string, startDate: string, endDate: string) {
    return APIHandler.get(`${resource}/work-session-team-weekly-total/${teamUuid}/${startDate}/${endDate}`)
  },

  // Taux de présence
  getPresenceRate(userUuid: string, startDate: string, endDate: string) {
    return APIHandler.get(`${resource}/presence-rate/${userUuid}/${startDate}/${endDate}`)
  }
}