import type { Team } from '@/types/Team'
import APIHandler from '../APIHandler'

const resource = 'teams'

export default {
  /**
   * Récupère toutes les équipes avec leurs membres
   * GET /teams
   */
  getAll() {
    return APIHandler.get<Team[]>(`${resource}`)
  },

  /**
   * Récupère une équipe spécifique avec ses membres
   * GET /teams/{uuid}
   */
  getTeam(uuid: string) {
    return APIHandler.get<Team>(`${resource}/${uuid}`)
  },

  /**
   * Crée une nouvelle équipe
   * POST /teams
   */
  createTeam(payload: {
    name: string
    description: string
    member_uuids: Array<{ user_uuid: string; is_manager: boolean }>
  }) {
    return APIHandler.post(`${resource}`, payload)
  },

  /**
   * Ajoute des membres à une équipe
   * POST /teams/add-users
   */
  addMembers(payload: {
    team_uuid: string
    member_uuids: Array<{ user_uuid: string; is_manager: boolean }>
  }) {
    return APIHandler.post(`${resource}/add-users`, payload)
  },

  /**
   * Supprime une équipe
   * DELETE /teams/{uuid}
   */
  deleteTeam(uuid: string) {
    return APIHandler.delete(`${resource}/${uuid}`)
  },

  /**
   * Retire un membre d'une équipe
   * DELETE /teams/users/{team_uuid}/{user_uuid}
   */
  removeMember(teamUuid: string, userUuid: string) {
    return APIHandler.delete(`${resource}/users/${teamUuid}/${userUuid}`)
  },
  /**
 * Modifie une équipe (nom et description)
 * PUT /teams/{uuid}
 */
  updateTeam(uuid: string, payload: { name: string; description: string }) {
    return APIHandler.put(`${resource}/edit/${uuid}`, payload)
  }

}