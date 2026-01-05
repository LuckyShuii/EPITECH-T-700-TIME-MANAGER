<script setup lang="ts">
import { computed } from 'vue'
import type { IndividualPauseData } from '@/types/kpi'

interface Props {
  data: IndividualPauseData | null
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  data: null
})

const emit = defineEmits<{
  viewDetails: [data: IndividualPauseData]
}>()

const handleClick = () => {
  if (props.data) {
    emit('viewDetails', props.data)
  }
}

const hasValidData = computed(() => {
  return props.data !== null && props.data.byDay.length > 0
})
</script>

<template>
  <div class="card bg-gradient-to-br from-cyan-500 to-blue-600 shadow-xl h-full cursor-pointer hover:shadow-2xl transition-all" @click="handleClick">
    <div class="card-body p-6 flex flex-col h-full text-white">
      <div class="mb-4">
        <h2 class="card-title text-xl font-bold">Temps de pause</h2>
        <p class="text-sm opacity-80">Semaine S-1</p>
      </div>
      
      <div v-if="loading" class="flex-1 flex justify-center items-center">
        <span class="loading loading-spinner loading-lg"></span>
      </div>

      <div v-else-if="!hasValidData" class="flex-1 flex flex-col justify-center items-center opacity-70">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-center mt-3">Aucune donnée disponible</p>
      </div>

      <div v-else-if="data" class="flex-1 flex flex-col justify-between">
        <div class="grid grid-cols-2 gap-4 mb-4">
          <div class="text-center">
            <div class="text-sm opacity-90 mb-1">Moyenne/jour</div>
            <div class="text-4xl font-bold">{{ data.averagePausePerDay }}min</div>
          </div>
          
          <div class="text-center">
            <div class="text-sm opacity-90 mb-1">Total semaine</div>
            <div class="text-4xl font-bold">{{ Math.floor(data.totalPauseWeek / 60) }}h{{ data.totalPauseWeek % 60 }}m</div>
          </div>
        </div>

        <div class="space-y-1 max-h-32 overflow-y-auto">
          <div class="text-xs opacity-80 mb-2">Détail par jour</div>
          <div 
            v-for="day in data.byDay" 
            :key="day.day"
            class="flex justify-between items-center px-3 py-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors"
          >
            <span class="text-sm font-medium capitalize">{{ day.day }}</span>
            <span class="text-sm font-bold">{{ day.minutes }}min</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>