import type { RegisterFormData } from '@/types/RegisterForm'
import type { PasswordResetPayload, PasswordResetResponse, PasswordChangePayload } from '@/types/PasswordReset'


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

  updateStatus(userUuid: string, status: 'active' | 'disabled') {
    return APIHandler.put(`${resource}/update-status`, {
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
    return APIHandler.get(`${resource}/${userUuid}`)
  },

  sendPasswordResetEmail(userEmail: string): Promise<PasswordResetResponse> {
    const payload: PasswordResetPayload = {
      user_email: userEmail
    }
    return APIHandler.post(`${resource}/reset-password`, payload)
  },
  changePassword(payload: PasswordChangePayload) {
    return APIHandler.post(`${resource}/update-password`, payload)
  }
}