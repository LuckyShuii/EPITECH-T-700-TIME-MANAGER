// Types pour la réinitialisation de mot de passe
export interface PasswordResetPayload {
  user_email: string
}

export interface PasswordResetResponse {
  message: string
}

export interface PasswordChangePayload {
  token: string
  new_password: string
}

export interface PasswordChangeResponse {
  message: string
}