import { defineStore } from 'pinia'
import { ref } from 'vue'
import LayoutService from '@/services/routers/LayoutService'

// Type pour un item de la grille
export interface GridItem {
  i: string // identifiant unique (nom du slot)
  x: number // position horizontale (0-3 pour une grille 4x4)
  y: number // position verticale (0-3 pour une grille 4x4)
  w: number // largeur en cases
  h: number // hauteur en cases
  minW?: number // largeur minimale (optionnel)
  minH?: number // hauteur minimale (optionnel)
  static?: boolean
}

// Type pour un layout complet
export type Layout = GridItem[]

// Layouts par défaut pour chaque dashboard
const defaultLayouts: Record<string, Layout> = {
  employee: [
  { i: 'clock', x: 0, y: 0, w: 1, h: 2, minW: 1, minH: 2 },
  { i: 'team-presence', x: 3, y: 0, w: 1, h: 4, minW: 1, minH: 4, static: true },
  { i: 'calendar', x: 1, y: 0, w: 2, h: 2, minW: 2, minH: 2 },
  { i: 'widget-2', x: 1, y: 2, w: 1, h: 2, minW: 1, minH: 2 },
  { i: 'widget-3', x: 2, y: 2, w: 1, h: 2, minW: 1, minH: 2 },
  { i: 'widget-4', x: 3, y: 3, w: 1, h: 1, minW: 1, minH: 1 }
],
 manager: [
  { i: 'clock', x: 0, y: 0, w: 1, h: 2, minW: 1, minH: 2 },
  { i: 'kpi-stats', x: 0, y: 2, w: 1, h: 2, minW: 1, minH: 2 },  // ← Agrandi à h: 2
  { i: 'calendar', x: 1, y: 0, w: 2, h: 2, minW: 2, minH: 2 },
  { i: 'team-view', x: 1, y: 2, w: 2, h: 2, minW: 2, minH: 2 },
  { i: 'team-presence', x: 3, y: 0, w: 1, h: 3, minW: 1, minH: 3, static: true },
  { i: 'modal-team', x: 3, y: 3, w: 1, h: 1, minW: 1, minH: 1 }
],
admin: [
  { i: 'add-employee', x: 0, y: 0, w: 1, h: 1, minW: 1, minH: 1 },
  { i: 'staff-settings', x: 0, y: 1, w: 1, h: 1, minW: 1, minH: 1 },
  { i: 'kpi-monthly', x: 0, y: 2, w: 1, h: 1, minW: 1, minH: 1 },
  { i: 'export-button', x: 0, y: 3, w: 1, h: 1, minW: 1, minH: 1 },
  { i: 'calendar', x: 1, y: 0, w: 2, h: 2, minW: 2, minH: 2 },
  { i: 'kpi-history', x: 1, y: 2, w: 2, h: 2, minW: 2, minH: 2 },
  { i: 'presence-rate', x: 3, y: 0, w: 1, h: 4, minW: 1, minH: 4 }
]
}

export const useLayoutStore = defineStore('layout', () => {
  const layouts = ref<Record<string, Layout>>({})

  // Charge le layout depuis l'API ou utilise le défaut
  async function loadLayout(dashboardName: string): Promise<Layout> {
    try {
      const response = await LayoutService.getLayout()
      
      if (response.data?.layout && response.data.layout.length > 0) {
        return response.data.layout
      }
    } catch (error) {
      console.error('Erreur lors du chargement du layout:', error)
    }
    
    // Si pas de sauvegarde ou erreur, retourne le layout par défaut
    return defaultLayouts[dashboardName] || []
  }

  // Sauvegarde le layout via l'API
  async function saveLayout(dashboardName: string, layout: Layout) {
    layouts.value[dashboardName] = layout
    
    try {
      await LayoutService.saveLayout(layout)
    } catch (error) {
      console.error('Erreur lors de la sauvegarde du layout:', error)
    }
  }

  // Réinitialise le layout aux valeurs par défaut
  async function resetLayout(dashboardName: string) {
    const defaultLayout = defaultLayouts[dashboardName]
    if (defaultLayout) {
      try {
        await LayoutService.deleteLayout()
        layouts.value[dashboardName] = defaultLayout
      } catch (error) {
        console.error('Erreur lors de la réinitialisation du layout:', error)
      }
    }
  }

  // Récupère le layout actuel pour un dashboard
  function getLayout(dashboardName: string): Layout {
    if (!layouts.value[dashboardName]) {
      loadLayout(dashboardName).then(layout => {
        layouts.value[dashboardName] = layout
      })
      return defaultLayouts[dashboardName] || []
    }
    return layouts.value[dashboardName]
  }

  return {
    layouts,
    loadLayout,
    saveLayout,
    resetLayout,
    getLayout
  }
})