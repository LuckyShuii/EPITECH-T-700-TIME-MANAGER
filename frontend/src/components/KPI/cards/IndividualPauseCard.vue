<script setup lang="ts">
import { computed } from 'vue'
import type { IndividualPauseData } from '@/types/Kpi'
import { useKpiStore } from '@/store/KpiStore'

interface Props {
  data: IndividualPauseData | null
  loading?: boolean
  weekLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  data: null,
  weekLabel: 'Semaine S-1'
})

const kpiStore = useKpiStore()

const hasValidData = computed(() => {
  return props.data !== null && props.data.average_break_time >= 0
})

const formattedPauseTime = computed(() => {
  if (!props.data) return '0 min'
  
  const minutes = props.data.average_break_time
  
  if (minutes < 60) {
    return `${minutes} min`
  }
  
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  
  if (mins === 0) {
    return `${hours}h`
  }
  
  return `${hours}h ${mins}m`
})

const handlePreviousWeek = () => {
  kpiStore.changeWeek('previous')
}

const handleCurrentWeek = () => {
  kpiStore.changeWeek('current')
}
</script>

<template>
  <div class="bg-white border-2 border-black p-6 h-full flex flex-col">
    <!-- Header avec flèche retour et bouton reset -->
    <div class="flex items-center justify-between mb-4">
      <button
        @click="handlePreviousWeek"
        class="text-2xl font-bold hover:opacity-70 transition"
        title="Semaine précédente"
      >
        &lt;
      </button>
      <h2 class="text-lg font-bold">Temps de pauses</h2>
      <button
        @click="handleCurrentWeek"
        class="px-3 py-1  hover:bg-black hover:text-white transition"
        title="Revenir à aujourd'hui"
      >
        ↻
      </button>
    </div>

    <!-- Label de la semaine -->
    <p class="text-sm mb-4 opacity-75">{{ weekLabel }}</p>

    <!-- Contenu -->
    <div v-if="loading" class="flex-1 flex justify-center items-center">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="!hasValidData" class="flex-1 flex flex-col justify-center items-center opacity-60">
      <p class="text-sm">Aucune donnée</p>
    </div>

    <div v-else class="flex-1 flex flex-col justify-center items-center">
      <p class="text-sm opacity-75 mb-2">Temps moyen</p>
      <p class="text-5xl font-bold font-mono">{{ formattedPauseTime }}</p>
    </div>
  </div>
</template>