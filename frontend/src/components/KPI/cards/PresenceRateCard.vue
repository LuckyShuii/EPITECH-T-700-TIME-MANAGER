<script setup lang="ts">
import { computed } from 'vue'
import type { PresenceRateData } from '@/types/kpi'
import type { ChartConfiguration } from 'chart.js'
import DoughnutChart from '../charts/DoughnutChart.vue'

interface Props {
  data: PresenceRateData | null
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  data: null
})

const emit = defineEmits<{
  viewDetails: [data: PresenceRateData]
}>()

const handleClick = () => {
  if (props.data) {
    emit('viewDetails', props.data)
  }
}

const hasValidData = computed(() => {
  return props.data !== null && props.data.daysExpected > 0
})

const chartData = computed<ChartConfiguration['data']>(() => {
  if (!props.data) {
    return { labels: [], datasets: [] }
  }

  const absent = props.data.daysExpected - props.data.daysPresent

  return {
    labels: ['Présent', 'Absent'],
    datasets: [{
      data: [props.data.daysPresent, absent],
      backgroundColor: [
        'rgb(34, 197, 94)',  // green-500
        'rgb(239, 68, 68)'   // red-500
      ],
      borderWidth: 0
    }]
  }
})

const chartOptions = computed<ChartConfiguration['options']>(() => ({
  responsive: true,
  maintainAspectRatio: true,
  cutout: '70%',
  plugins: {
    legend: {
      display: false
    },
    tooltip: {
      callbacks: {
        label: (context) => {
          return `${context.label}: ${context.parsed} jour(s)`
        }
      }
    }
  }
}))
</script>

<template>
  <div class="card bg-base-100 shadow-xl h-full cursor-pointer hover:shadow-2xl transition-shadow" @click="handleClick">
    <div class="card-body p-6 flex flex-col h-full">
      <div class="mb-3">
        <h2 class="card-title text-lg">Taux de présence</h2>
        <p class="text-xs text-base-content/70">Semaine S-1</p>
      </div>
      
      <div v-if="loading" class="flex-1 flex justify-center items-center">
        <span class="loading loading-spinner loading-lg"></span>
      </div>

      <div v-else-if="!hasValidData" class="flex-1 flex flex-col justify-center items-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-base-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-center text-sm text-base-content/70 mt-3">Aucune donnée disponible</p>
      </div>

      <div v-else-if="data" class="flex-1 flex flex-col">
        <div class="relative flex-1 flex items-center justify-center">
          <div class="w-full max-w-[200px]">
            <DoughnutChart :data="chartData" :options="chartOptions" />
          </div>
          
          <div class="absolute inset-0 flex flex-col items-center justify-center">
            <div class="text-4xl font-bold text-primary">{{ data.presenceRate }}%</div>
            <div class="text-xs text-base-content/70">présence</div>
          </div>
        </div>

        <div class="mt-4 grid grid-cols-2 gap-2">
          <div class="stat bg-success/10 rounded-lg p-2">
            <div class="stat-title text-xs">Jours présents</div>
            <div class="stat-value text-lg text-success">{{ data.daysPresent }}</div>
          </div>
          
          <div class="stat bg-base-200 rounded-lg p-2">
            <div class="stat-title text-xs">Jours attendus</div>
            <div class="stat-value text-lg">{{ data.daysExpected }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>