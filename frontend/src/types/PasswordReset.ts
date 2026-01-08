// Types pour la réinitialisation de mot de passe
export interface PasswordResetPayload {
  user_email: string
}

export interface PasswordResetResponse {
  message: string
}

export interface PasswordChangePayload {
  new_password: string
  user_uuid: string
}

export interface PasswordChangeResponse {
  message: string
}