<script setup lang="ts">
import { computed, ref } from 'vue'
import { Square3Stack3DIcon, ArrowPathIcon } from '@heroicons/vue/24/outline'
import { useEditModeStore } from '@/store/EditModeStore'
import { useLayoutStore } from '@/store/LayoutStore'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

const editModeStore = useEditModeStore()
const layoutStore = useLayoutStore()

const isEditMode = computed(() => editModeStore.isEditMode)
const showResetConfirm = ref(false)

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
    class="w-10 h-10 flex items-center justify-center font-bold hover:bg-black hover:text-white transition-none rounded-none"
    title="Réinitialiser le dashboard"
  >
    <ArrowPathIcon class="w-5 h-5" />
  </button>

  <!-- Bouton personnalisation -->
  <button
    @click="toggleEditMode"
    :class="[
      'w-10 h-10 flex items-center justify-center font-bold hover:bg-black hover:text-white transition-none rounded-none',
      isEditMode ? 'bg-black text-white' : ''
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