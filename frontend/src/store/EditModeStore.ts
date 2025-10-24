import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useLayoutStore } from './LayoutStore'

export const useEditModeStore = defineStore('editMode', () => {
  // État : est-ce qu'on est en mode édition ?
  const isEditMode = ref(false)
  
  // Le dashboard actuellement actif
  const currentDashboard = ref<string | null>(null)

  // Active/désactive le mode édition
  function toggleEditMode() {
    // Si on SORT du mode édition, on sauvegarde
    if (isEditMode.value && currentDashboard.value) {
      const layoutStore = useLayoutStore()
      const currentLayout = layoutStore.getLayout(currentDashboard.value)
      layoutStore.saveLayout(currentDashboard.value, currentLayout)
    }
    
    isEditMode.value = !isEditMode.value
  }

  // Définit le dashboard actif
  function setCurrentDashboard(dashboardName: string) {
    currentDashboard.value = dashboardName
  }

  // Réinitialise tout (utile lors de la déconnexion ou changement de page)
  function reset() {
    isEditMode.value = false
    currentDashboard.value = null
  }

  return {
    isEditMode,
    currentDashboard,
    toggleEditMode,
    setCurrentDashboard,
    reset
  }
})