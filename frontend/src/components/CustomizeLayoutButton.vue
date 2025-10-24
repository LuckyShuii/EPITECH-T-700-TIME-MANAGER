<script setup lang="ts">
import { computed } from 'vue'
import { PencilSquareIcon } from '@heroicons/vue/24/outline'
import { useEditModeStore } from '@/store/EditModeStore'

// Store
const editModeStore = useEditModeStore()

// État du mode édition
const isEditMode = computed(() => editModeStore.isEditMode)

// Action
function toggleEditMode() {
  editModeStore.toggleEditMode()
}
</script>

<template>
  <button
    @click="toggleEditMode"
    :class="[
      'btn btn-ghost btn-circle',
      isEditMode ? 'text-green-500' : ''
    ]"
    :title="isEditMode ? 'Terminer la personnalisation' : 'Personnaliser le dashboard'"
  >
    <PencilSquareIcon class="w-5 h-5" />
  </button>
</template>