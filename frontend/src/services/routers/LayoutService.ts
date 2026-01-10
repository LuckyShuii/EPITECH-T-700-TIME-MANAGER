import type { Layout } from '@/store/LayoutStore'
import APIHandler from '../APIHandler'
const resource = 'users/current-user-dashboard-layout'

export default {
  // Récupère le layout sauvegardé
  getLayout() {
    return APIHandler.get(resource)
  },

  // Sauvegarde le layout
  saveLayout(layout: Layout) {
    return APIHandler.put(`${resource}/edit`, { layout })
  },

  // Supprime le layout (réinitialisation)
  deleteLayout() {
    return APIHandler.delete(`${resource}/delete`)
  }
}