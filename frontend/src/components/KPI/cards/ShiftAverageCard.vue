<script setup lang="ts">
import { computed } from 'vue'
import type { ShiftAverageData } from '@/types/kpi'
import type { ChartConfiguration } from 'chart.js'
import LineChart from '../charts/LineChart.vue'

interface Props {
  data: ShiftAverageData | null
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  data: null
})

const emit = defineEmits<{
  viewDetails: [data: ShiftAverageData]
}>()

const handleClick = () => {
  if (props.data) {
    emit('viewDetails', props.data)
  }
}

// Vérifier si on a des données valides
const hasValidData = computed(() => {
  return props.data !== null && 
         props.data.labels.length > 0 && 
         props.data.values.length > 0
})

// Configuration du graphique Chart.js
const chartData = computed<ChartConfiguration['data']>(() => {
  if (!props.data) {
    return {
      labels: [],
      datasets: []
    }
  }

  return {
    labels: props.data.labels,
    datasets: [
      {
        label: 'Heures moyennes par shift',
        data: props.data.values,
        borderColor: 'rgb(59, 130, 246)', // Bleu primary
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        tension: 0.3, // Courbe légèrement arrondie
        fill: true,
        pointRadius: 5,
        pointHoverRadius: 7
      }
    ]
  }
})

const chartOptions = computed<ChartConfiguration['options']>(() => ({
  responsive: true,
  maintainAspectRatio: true,
  plugins: {
    legend: {
      display: false // On cache la légende car elle est redondante avec le titre
    },
    tooltip: {
      callbacks: {
        label: (context) => {
          return `${context.parsed.y.toFixed(1)}h par shift`
        }
      }
    }
  },
  scales: {
    y: {
      beginAtZero: false, // Ne pas forcer à 0 pour mieux voir les variations
      ticks: {
        callback: (value) => `${value}h`
      }
    }
  }
}))

// Calculer la moyenne globale
const globalAverage = computed(() => {
  if (!props.data || props.data.values.length === 0) return 0
  const sum = props.data.values.reduce((acc, val) => acc + val, 0)
  return sum / props.data.values.length
})

// Calculer la tendance (dernière valeur vs première valeur)
const trend = computed(() => {
  if (!props.data || props.data.values.length < 2) return 0
  const first = props.data.values[0]
  const last = props.data.values[props.data.values.length - 1]
  return last - first
})

const trendLabel = computed(() => {
  if (trend.value > 0) return 'En hausse'
  if (trend.value < 0) return 'En baisse'
  return 'Stable'
})

const trendColor = computed(() => {
  if (trend.value > 0) return 'text-success'
  if (trend.value < 0) return 'text-error'
  return 'text-base-content'
})
</script>

<template>
  <div class="card bg-base-100 shadow-xl cursor-pointer hover:shadow-2xl transition-shadow" @click="handleClick">
    <div class="card-body">
      <h2 class="card-title">Moyenne par shift</h2>
      <p class="text-sm text-base-content/70">Évolution sur 4 semaines</p>
      
      <!-- État de chargement -->
      <div v-if="loading" class="flex justify-center items-center h-64">
        <span class="loading loading-spinner loading-lg"></span>
      </div>

      <!-- Pas de données disponibles -->
      <div v-else-if="!hasValidData" class="flex flex-col justify-center items-center h-64 space-y-4">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-base-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <p class="text-center text-base-content/70">
          Aucune donnée disponible pour le moment.
        </p>
      </div>

      <!-- Données disponibles -->
      <div v-else-if="data" class="space-y-4">
        <!-- Statistiques en haut -->
        <div class="stats shadow w-full">
          <div class="stat">
            <div class="stat-title">Moyenne globale</div>
            <div class="stat-value text-2xl text-primary">{{ globalAverage.toFixed(1) }}h</div>
          </div>
          
          <div class="stat">
            <div class="stat-title">Tendance</div>
            <div class="stat-value text-2xl" :class="trendColor">
              {{ trendLabel }}
            </div>
            <div class="stat-desc" :class="trendColor">
              {{ trend > 0 ? '+' : '' }}{{ trend.toFixed(1) }}h
            </div>
          </div>
        </div>

        <!-- Graphique -->
        <div class="w-full h-64">
          <LineChart :data="chartData" :options="chartOptions" />
        </div>
      </div>
    </div>
  </div>
</template>