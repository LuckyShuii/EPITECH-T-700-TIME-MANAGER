import type { RegisterFormData } from '@/types/RegisterForm'
import APIHandler from '../APIHandler'

const resource = 'users'

export default {
  register(payload: RegisterFormData) {
    return APIHandler.post(`${resource}/register`, payload)
  }
}