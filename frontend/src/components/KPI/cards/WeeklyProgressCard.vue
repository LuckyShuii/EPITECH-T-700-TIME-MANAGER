<script setup lang="ts">
import type { WeeklyProgressData } from '@/types/kpi'

interface Props {
  data: WeeklyProgressData | null  // ← CHANGEMENT ICI
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  data: null  // ← AJOUT ICI
})

const emit = defineEmits<{
  viewDetails: [data: WeeklyProgressData]
}>()

const handleClick = () => {
  if (props.data) {
    emit('viewDetails', props.data)
  }
}

const formatWeekStart = (date: string): string => {
  const d = new Date(date)
  return d.toLocaleDateString('fr-FR', { 
    day: '2-digit', 
    month: 'long',
    year: 'numeric'
  })
}

// Helper pour vérifier si on a des données valides
const hasValidData = computed(() => {
  return props.data !== null && props.data.contractHours > 0
})
</script>

<template>
  <div class="card bg-base-100 shadow-xl cursor-pointer hover:shadow-2xl transition-shadow" @click="handleClick">
    <div class="card-body">
      <h2 class="card-title">Progression hebdomadaire</h2>
      <p v-if="data" class="text-sm text-base-content/70">Semaine du {{ formatWeekStart(data.weekStartDate) }}</p>
      
      <!-- État de chargement -->
      <div v-if="loading" class="flex justify-center items-center h-32">
        <span class="loading loading-spinner loading-lg"></span>
      </div>

      <!-- Pas de données disponibles -->
      <div v-else-if="!hasValidData" class="flex flex-col justify-center items-center h-32 space-y-4">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-base-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <p class="text-center text-base-content/70">
          Aucune donnée disponible pour cette semaine.<br>
          <span class="text-sm">Les données seront disponibles dès demain.</span>
        </p>
      </div>

      <!-- Données disponibles -->
      <div v-else class="space-y-4">
        <!-- Barre de progression -->
        <div class="w-full">
          <div class="flex justify-between text-sm mb-2">
            <span class="font-medium">{{ data.workedHours }}h / {{ data.contractHours }}h</span>
            <span class="font-bold text-primary">{{ data.percentageComplete.toFixed(1) }}%</span>
          </div>
          <progress 
            class="progress progress-primary w-full" 
            :value="data.percentageComplete" 
            max="100"
          ></progress>
        </div>

        <!-- Statistiques -->
        <div class="stats stats-vertical lg:stats-horizontal shadow w-full">
          <div class="stat">
            <div class="stat-title">Heures travaillées</div>
            <div class="stat-value text-2xl text-success">{{ data.workedHours }}h</div>
          </div>
          
          <div class="stat">
            <div class="stat-title">Heures restantes</div>
            <div class="stat-value text-2xl" :class="data.remainingHours > 0 ? 'text-warning' : 'text-success'">
              {{ data.remainingHours }}h
            </div>
          </div>
        </div>

        <!-- Info contexte -->
        <div class="alert alert-info">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <span class="text-xs">Données mises à jour quotidiennement (J+1)</span>
        </div>
      </div>
    </div>
  </div>
</template>