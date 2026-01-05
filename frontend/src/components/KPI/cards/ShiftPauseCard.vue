<script setup lang="ts">
import { computed } from 'vue'
import type { ShiftPauseData } from '@/types/kpi'
import type { ChartConfiguration } from 'chart.js'
import PieChart from '../charts/PieChart.vue'

interface Props {
  data: ShiftPauseData | null
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  data: null
})

const emit = defineEmits<{
  viewDetails: [data: ShiftPauseData]
}>()

const handleClick = () => {
  if (props.data) {
    emit('viewDetails', props.data)
  }
}

const hasValidData = computed(() => {
  return props.data !== null
})

const chartData = computed<ChartConfiguration['data']>(() => {
  if (!props.data) {
    return { labels: [], datasets: [] }
  }

  return {
    labels: ['Travail effectif', 'Pause'],
    datasets: [{
      data: [props.data.effectiveWork, props.data.pause],
      backgroundColor: [
        'rgb(34, 197, 94)',  // green-500
        'rgb(251, 146, 60)'  // orange-400
      ],
      borderWidth: 2,
      borderColor: '#fff'
    }]
  }
})

const chartOptions = computed<ChartConfiguration['options']>(() => ({
  responsive: true,
  maintainAspectRatio: true,
  plugins: {
    legend: {
      position: 'bottom'
    },
    tooltip: {
      callbacks: {
        label: (context) => {
          const minutes = context.parsed
          const hours = Math.floor(minutes / 60)
          const mins = minutes % 60
          return `${context.label}: ${hours}h${mins}min`
        }
      }
    }
  }
}))

const formatMinutes = (minutes: number) => {
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return `${h}h${m}m`
}
</script>

<template>
  <div class="card bg-base-100 shadow-xl h-full cursor-pointer hover:shadow-2xl transition-shadow" @click="handleClick">
    <div class="card-body p-4 flex flex-col h-full">
      <div class="mb-3">
        <h2 class="card-title text-lg">Répartition shift</h2>
        <p class="text-xs text-base-content/70">Jour J-1</p>
      </div>
      
      <div v-if="loading" class="flex-1 flex justify-center items-center">
        <span class="loading loading-spinner loading-lg"></span>
      </div>

      <div v-else-if="!hasValidData" class="flex-1 flex flex-col justify-center items-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-base-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z M20.488 9H15V3.512A9.025 9.025 0 0120.488 9z" />
        </svg>
        <p class="text-center text-sm text-base-content/70 mt-3">Aucune donnée disponible</p>
      </div>

      <div v-else-if="data" class="flex-1 flex flex-col">
        <div class="flex-1 flex items-center justify-center">
          <div class="w-full max-w-[200px]">
            <PieChart :data="chartData" :options="chartOptions" />
          </div>
        </div>

        <div class="mt-3 grid grid-cols-2 gap-2">
          <div class="stat bg-success/10 rounded-lg p-2">
            <div class="stat-title text-xs">Travail</div>
            <div class="stat-value text-sm text-success">{{ formatMinutes(data.effectiveWork) }}</div>
          </div>
          
          <div class="stat bg-orange-100 rounded-lg p-2">
            <div class="stat-title text-xs">Pause</div>
            <div class="stat-value text-sm text-orange-600">{{ formatMinutes(data.pause) }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>