<script setup lang="ts">
import { computed, ref } from 'vue'
import { Square3Stack3DIcon, ArrowPathIcon } from '@heroicons/vue/24/outline'
import { useEditModeStore } from '@/store/editModeStore'
import { useLayoutStore } from '@/store/layoutStore'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

// Stores
const editModeStore = useEditModeStore()
const layoutStore = useLayoutStore()



// État
const isEditMode = computed(() => editModeStore.isEditMode)
const showResetConfirm = ref(false)

// Actions
function toggleEditMode() {
  editModeStore.toggleEditMode()
}

function confirmReset() {
  showResetConfirm.value = true
}

function resetLayout() {
  const dashboardName = editModeStore.currentDashboard
  if (dashboardName) {
    layoutStore.resetLayout(dashboardName)
    window.location.reload()
  }
}
</script>

<template>
  <!-- Bouton reset -->
  <button
    @click="confirmReset"
    class="btn btn-ghost btn-circle"
    title="Réinitialiser le dashboard"
  >
    <ArrowPathIcon class="w-5 h-5" />
  </button>

  <!-- Bouton personnalisation -->
  <button
    @click="toggleEditMode"
    :class="[
      'btn btn-ghost btn-circle',
      isEditMode ? 'text-green-500' : ''
    ]"
    :title="isEditMode ? 'Terminer la personnalisation' : 'Personnaliser le dashboard'"
  >
    <Square3Stack3DIcon class="w-5 h-5" />
  </button>

  <!-- Modal de confirmation -->
  <ConfirmDialog
    v-model="showResetConfirm"
    title="Réinitialiser le dashboard"
    message="Voulez-vous réinitialiser votre disposition par défaut ?<br><span class='text-sm opacity-70'>Cette action est irréversible et rechargera la page.</span>"
    confirm-text="Réinitialiser"
    cancel-text="Annuler"
    variant="warning"
    @confirm="resetLayout"
  />
</template>