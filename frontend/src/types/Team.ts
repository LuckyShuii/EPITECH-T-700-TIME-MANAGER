/**
 * Représente un membre d'équipe avec toutes ses informations
 * Retourné par l'API dans /teams/{uuid} et /teams
 */

export interface TeamMember {
  user_uuid: string
  username: string
  first_name: string
  last_name: string
  email: string
  phone_number: string
  roles: string[]
  status: 'active' | 'inactive' | 'pending'
  is_manager: boolean
  work_session_status: 'active' | 'paused' | 'no_active_session'  // ← Type plus strict
}

/**
 * Représente une équipe complète avec tous ses membres
 * Retourné par l'API dans /teams/{uuid} et /teams
 */
export interface Team {
  uuid: string
  name: string
  description: string
  team_members: TeamMember[]
}

/**
 * Représente les informations basiques d'une équipe
 * Trouvé dans la réponse de /users/{uuid} dans le champ teams[]
 */
export interface TeamInfo {
  team_uuid: string
  team_name: string
  team_description: string
  is_manager: boolean
}