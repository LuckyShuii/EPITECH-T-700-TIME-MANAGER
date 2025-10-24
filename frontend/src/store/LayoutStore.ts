import { defineStore } from 'pinia'
import { ref } from 'vue'

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
    { i: 'team-presence', x: 3, y: 0, w: 1, h: 3, minW: 1, minH: 3, static: true },
    { i: 'widget-4', x: 3, y: 3, w: 1, h: 1, minW: 1, minH: 1 },
    { i: 'widget-6', x: 0, y: 2, w: 1, h: 1, minW: 1, minH: 1 },
    { i: 'extra-widget', x: 0, y: 3, w: 1, h: 1, minW: 1, minH: 1 },
    { i: 'calendar', x: 1, y: 0, w: 2, h: 2, minW: 2, minH: 2 },
    { i: 'widget-2', x: 1, y: 2, w: 2, h: 1, minW: 2, minH: 1 },
    { i: 'widget-3', x: 1, y: 3, w: 2, h: 1, minW: 2, minH: 1 }
  ],
  manager: [
    { i: 'clock', x: 0, y: 0, w: 1, h: 2, minW: 1, minH: 2 },
    { i: 'kpi-stats', x: 0, y: 2, w: 1, h: 1, minW: 1, minH: 1 },
    { i: 'team-view', x: 0, y: 3, w: 1, h: 1, minW: 1, minH: 1 },
    { i: 'calendar', x: 1, y: 0, w: 2, h: 2, minW: 2, minH: 2 },
    { i: 'report-button', x: 1, y: 2, w: 2, h: 2, minW: 2, minH: 2 },
    { i: 'team-presence', x: 3, y: 0, w: 1, h: 3, minW: 1, minH: 3, static: true },
    { i: 'modal-team', x: 3, y: 3, w: 1, h: 1, minW: 1, minH: 1 }
  ],
  admin: [
    { i: 'add-employee', x: 0, y: 0, w: 1, h: 1, minW: 1, minH: 1 },
    { i: 'staff-settings', x: 0, y: 1, w: 1, h: 1, minW: 1, minH: 1 },
    { i: 'kpi-monthly', x: 0, y: 2, w: 1, h: 1, minW: 1, minH: 1 },
    { i: 'kpi-history', x: 0, y: 3, w: 1, h: 1, minW: 1, minH: 1 },
    { i: 'calendar', x: 1, y: 0, w: 2, h: 2, minW: 2, minH: 2 },
    { i: 'widget-6', x: 1, y: 2, w: 2, h: 2, minW: 2, minH: 2 },
    { i: 'remote-absence', x: 3, y: 0, w: 1, h: 2, minW: 1, minH: 2 },
    { i: 'manager-report', x: 3, y: 2, w: 1, h: 2, minW: 1, minH: 2 }
  ]
}

export const useLayoutStore = defineStore('layout', () => {
  // État : layouts actuels pour chaque dashboard
  const layouts = ref<Record<string, Layout>>({})

  // Charge le layout depuis localStorage ou utilise le défaut
  function loadLayout(dashboardName: string): Layout {
    const savedLayout = localStorage.getItem(`layout_${dashboardName}`)
    
    if (savedLayout) {
      try {
        return JSON.parse(savedLayout)
      } catch (e) {
        console.error('Erreur lors du chargement du layout:', e)
      }
    }
    
    // Si pas de sauvegarde ou erreur, retourne le layout par défaut
    return defaultLayouts[dashboardName] || []
  }

  // Sauvegarde le layout (localStorage pour l'instant, API plus tard)
  function saveLayout(dashboardName: string, layout: Layout) {
    // Met à jour l'état local
    layouts.value[dashboardName] = layout
    
    // Sauvegarde localStorage (temporaire)
    localStorage.setItem(`layout_${dashboardName}`, JSON.stringify(layout))
    
    // TODO: Quand l'API sera prête, décommenter et adapter :
    /*
    try {
      await fetch('/api/user/layouts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({ dashboardName, layout })
      })
    } catch (e) {
      console.error('Erreur lors de la sauvegarde du layout via API:', e)
    }
    */
  }

  // Réinitialise le layout aux valeurs par défaut
  function resetLayout(dashboardName: string) {
    const defaultLayout = defaultLayouts[dashboardName]
    if (defaultLayout) {
      saveLayout(dashboardName, defaultLayout)
    }
  }

  // Récupère le layout actuel pour un dashboard
  function getLayout(dashboardName: string): Layout {
    if (!layouts.value[dashboardName]) {
      layouts.value[dashboardName] = loadLayout(dashboardName)
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