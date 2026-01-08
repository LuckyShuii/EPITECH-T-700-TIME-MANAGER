<script setup lang="ts">
import { computed } from 'vue'
import { GridLayout, GridItem } from 'vue3-grid-layout'
import { useLayoutStore } from '@/store/LayoutStore'
import { useEditModeStore } from '@/store/EditModeStore'
import type { Layout } from '@/store/LayoutStore'

interface Props {
  dashboardName: string
}

const props = defineProps<Props>()
const layoutStore = useLayoutStore()
const editModeStore = useEditModeStore()

const currentLayout = computed({
  get: () => {
    const layout = layoutStore.getLayout(props.dashboardName)
    console.log('📊 Layout chargé:', layout)
    return layout
  },
  set: (newLayout: Layout) => {
    layoutStore.saveLayout(props.dashboardName, newLayout)
  }
})

const isEditMode = computed(() => editModeStore.isEditMode)
</script>

<template>
  <div class="relative">
    <!-- Indicateur visuel du mode édition -->
    <div
      v-if="isEditMode"
      class="absolute -top-9 left-0 right-0 px-4 py-2 border-2 text-center font-black uppercase tracking-wider z-50 edit-mode-indicator"
    >
      MODE PERSONNALISATION - DÉPLACEZ VOS WIDGETS
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
/* Indicateur du mode édition */
.edit-mode-indicator {
  background-color: white;
  color: black;
  border-color: black;
}

@media (prefers-color-scheme: dark) {
  .edit-mode-indicator {
    background-color: black;
    color: white;
    border-color: white;
  }
}

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