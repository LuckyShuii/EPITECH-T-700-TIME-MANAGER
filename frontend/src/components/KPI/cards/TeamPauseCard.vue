<script setup lang="ts">
import { computed } from 'vue'
import type { TeamPauseData } from '@/types/kpi'
import type { ChartConfiguration } from 'chart.js'
import BarChart from '../charts/BarChart.vue'

interface Props {
  data: TeamPauseData | null
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  data: null
})

const emit = defineEmits<{
  viewDetails: [data: TeamPauseData]
}>()

const handleClick = () => {
  if (props.data) {
    emit('viewDetails', props.data)
  }
}

const hasValidData = computed(() => {
  return props.data !== null && props.data.values.length > 0
})

const chartData = computed<ChartConfiguration['data']>(() => {
  if (!props.data) {
    return { labels: [], datasets: [] }
  }

  return {
    labels: props.data.labels,
    datasets: [{
      label: 'Minutes de pause moyenne',
      data: props.data.values,
      backgroundColor: 'rgba(14, 165, 233, 0.8)',
      borderColor: 'rgb(14, 165, 233)',
      borderWidth: 1
    }]
  }
})

const chartOptions = computed<ChartConfiguration['options']>(() => ({
  responsive: true,
  maintainAspectRatio: true,
  plugins: {
    legend: {
      display: false
    },
    tooltip: {
      callbacks: {
        label: (context) => `${context.parsed.y} min`
      }
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: {
        callback: (value) => `${value}min`
      }
    }
  }
}))
</script>

<template>
  <div class="card bg-base-100 shadow-xl h-full cursor-pointer hover:shadow-2xl transition-shadow" @click="handleClick">
    <div class="card-body p-4 flex flex-col h-full">
      <div class="mb-3">
        <h2 class="card-title text-lg">Pause moyenne équipe</h2>
        <p class="text-xs text-base-content/70">Par membre - Semaine S-1</p>
      </div>
      
      <div v-if="loading" class="flex-1 flex justify-center items-center">
        <span class="loading loading-spinner loading-lg"></span>
      </div>

      <div v-else-if="!hasValidData" class="flex-1 flex flex-col justify-center items-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-base-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-center text-sm text-base-content/70 mt-3">Aucune donnée disponible</p>
      </div>

      <div v-else class="flex-1 flex flex-col overflow-hidden">
        <BarChart :data="chartData" :options="chartOptions" />
      </div>
    </div>
  </div>
</template>