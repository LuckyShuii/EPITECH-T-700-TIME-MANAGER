export interface RegisterFormData {
  first_name: string
  last_name: string
  email: string
  username: string
  phone_number: string
  roles: string[]
}

export interface RegisterFormErrors {
  first_name?: string
  last_name?: string
  email?: string
  username?: string
  phone_number?: string
  roles?: string
  password?: string
  general?: string  // Pour les erreurs globales de l'API
}