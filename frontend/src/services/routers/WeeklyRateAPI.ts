import type { WeeklyRate, WeeklyRateCreateData, WeeklyRateUpdateData } from '@/types/WeeklyRate'
import APIHandler from '../APIHandler'

const resource = 'users/weekly-rates'

export default {
  getAll() {
    return APIHandler.get(resource)
  },

  create(payload: WeeklyRateCreateData) {
    return APIHandler.post(`${resource}/create`, payload)
  },

  update(uuid: string, payload: WeeklyRateUpdateData) {
    return APIHandler.put(`${resource}/${uuid}/update`, payload)
  },

  delete(uuid: string) {
    return APIHandler.delete(`${resource}/${uuid}/delete`)
  },

  assignToUser(weeklyRateUuid: string, userUuid: string) {
    return APIHandler.post(`${resource}/${weeklyRateUuid}/assign-to-user/${userUuid}`)
  }
}