<script setup lang="ts">
import { computed } from 'vue'
import { GridLayout, GridItem } from 'vue3-grid-layout'
import { useLayoutStore } from '@/store/LayoutStore'
import { useEditModeStore } from '@/store/EditModeStore'  // ← NOUVEAU
import type { Layout } from '@/store/LayoutStore'

// Props du composant
interface Props {
  dashboardName: string // 'employee', 'manager' ou 'admin'
}

const props = defineProps<Props>()

// Stores
const layoutStore = useLayoutStore()
const editModeStore = useEditModeStore()  // ← NOUVEAU

// Layout actuel (réactif)
const currentLayout = computed({
  get: () => {
    const layout = layoutStore.getLayout(props.dashboardName)
    console.log('📊 Layout chargé:', layout)  // ← AJOUTE
    return layout
  },
  set: (newLayout: Layout) => {
    layoutStore.saveLayout(props.dashboardName, newLayout)
  }
})

// Le mode édition vient maintenant du store
const isEditMode = computed(() => editModeStore.isEditMode)  // ← MODIFIÉ

// SUPPRIME ces fonctions (elles sont maintenant dans le store) :
// function toggleEditMode() { ... }
// function resetLayout() { ... }

// SUPPRIME le defineExpose (on n'en a plus besoin)
</script>

<template>
  <div class="relative">
    <!-- Indicateur visuel du mode édition -->
    <div
      v-if="isEditMode"
      class="absolute -top-12 left-0 right-0 bg-blue-500 text-white px-4 py-2 rounded-t-lg text-center font-medium z-50"
    >
      Mode personnalisation activé - Déplacez vos widgets
    </div>

    <!-- La grille draggable -->
    <GridLayout
      v-model:layout="currentLayout"
      :col-num="4"
      :row-height="200"
      :is-draggable="isEditMode"
      :is-resizable="false"
      :vertical-compact="true"
      :use-css-transforms="true"
      :margin="[12, 12]"
      :prevent-collision="false"
      class="min-h-screen"
    >
      <!-- Génération des GridItem pour chaque widget -->
      <GridItem
        v-for="item in currentLayout"
        :key="item.i"
        :x="item.x"
        :y="item.y"
        :w="item.w"
        :h="item.h"
        :i="item.i"
        :min-w="item.minW"
        :min-h="item.minH"
        :static="item.static"
        :class="{ 
          'cursor-move': isEditMode && !item.static,
          'ring-2 ring-blue-400': isEditMode && !item.static,
          'opacity-75': isEditMode && item.static
          }"
        class="transition-all"
      >
        <!-- Le slot correspondant au widget -->
        <slot :name="item.i"></slot>
      </GridItem>
    </GridLayout>
  </div>
</template>

<style scoped>
/* Styles pour vue-grid-layout */
:deep(.vue-grid-item) {
  transition: all 0.2s ease;
}

:deep(.vue-grid-item.vue-grid-placeholder) {
  background: rgb(59, 130, 246, 0.3);
  border-radius: 0.75rem;
  border: 2px dashed rgb(59, 130, 246);
}
</style>