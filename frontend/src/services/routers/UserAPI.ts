import type { RegisterFormData } from '@/types/RegisterForm'
import APIHandler from '../APIHandler'


const resource = 'users'

export default {
  register(payload: RegisterFormData) {
    return APIHandler.post(`${resource}/register`, payload)
  },

  getAll() {
    return APIHandler.get(resource)
  },

  getUser(userUuid: string) {
    return APIHandler.get(`${resource}/${userUuid}`)
  },

  deleteUser(userUuid: string) {
    return APIHandler.delete(`${resource}/delete/${userUuid}`)
  },

  updateStatus(userUuid: string, status: 'active' | 'inactive') {
    return APIHandler.post(`${resource}/update-status`, {
      user_uuid: userUuid,
      status: status
    })
  },

  update(userUuid: string, payload: any) {
    return APIHandler.put(`${resource}`, {
      ...payload,
      uuid: userUuid
    })
  },

  getUserSpecific(userUuid: string) {
    return APIHandler.get(`${resource}/specific/${userUuid}`)
  }
}