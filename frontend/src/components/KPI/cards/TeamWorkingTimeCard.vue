<script setup lang="ts">
import { computed } from 'vue'
import type { TeamWorkingTimeData } from '@/types/kpi'

interface Props {
  data: TeamWorkingTimeData | null
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  data: null
})

const emit = defineEmits<{
  viewDetails: [data: TeamWorkingTimeData]
}>()

const handleClick = () => {
  if (props.data) {
    emit('viewDetails', props.data)
  }
}

// Vérifier si on a des données valides
const hasValidData = computed(() => {
  return props.data !== null && 
         props.data.members.length > 0
})

// Calculer la moyenne d'heures par membre
const averageHours = computed(() => {
  if (!props.data || props.data.members.length === 0) return 0
  return props.data.totalTeam / props.data.members.length
})

// Trouver le membre avec le plus d'heures
const topMember = computed(() => {
  if (!props.data || props.data.members.length === 0) return null
  return props.data.members.reduce((max, member) => 
    member.hours > max.hours ? member : max
  )
})
</script>

<template>
  <div class="card bg-base-100 shadow-xl h-full cursor-pointer hover:shadow-2xl transition-shadow" @click="handleClick">
    <div class="card-body p-4 flex flex-col h-full">
      <!-- Header -->
      <div class="mb-3">
        <h2 class="card-title text-lg">Travail hebdomadaire équipe</h2>
        <p class="text-xs text-base-content/70">Semaine S-1</p>
      </div>
      
      <!-- État de chargement -->
      <div v-if="loading" class="flex-1 flex justify-center items-center">
        <span class="loading loading-spinner loading-lg"></span>
      </div>

      <!-- Pas de données disponibles -->
      <div v-else-if="!hasValidData" class="flex-1 flex flex-col justify-center items-center space-y-3">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-base-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        <p class="text-center text-sm text-base-content/70">
          Aucune donnée disponible
        </p>
      </div>

      <!-- Données disponibles -->
      <div v-else-if="data" class="flex-1 flex flex-col space-y-3 overflow-hidden">
        <!-- Statistiques compactes -->
        <div class="grid grid-cols-2 gap-2">
          <div class="stat bg-base-200 rounded-lg p-3">
            <div class="stat-title text-xs">Total équipe</div>
            <div class="stat-value text-2xl text-primary">{{ data.totalTeam }}h</div>
          </div>
          
          <div class="stat bg-base-200 rounded-lg p-3">
            <div class="stat-title text-xs">Moyenne</div>
            <div class="stat-value text-2xl text-secondary">{{ averageHours.toFixed(1) }}h</div>
          </div>
        </div>

        <!-- Top performer compact -->
        <div v-if="topMember" class="alert alert-success py-2 px-3">
          <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-5 w-5" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="text-xs">
            <strong>{{ topMember.name }}</strong> : {{ topMember.hours }}h
          </span>
        </div>

        <!-- Liste scrollable des membres -->
        <div class="flex-1 overflow-hidden flex flex-col">
          <div class="text-xs font-semibold mb-2 text-base-content/70">
            {{ data.members.length }} membre(s)
          </div>
          
          <div class="flex-1 overflow-y-auto space-y-1 pr-1">
            <div 
              v-for="member in data.members" 
              :key="member.userId"
              class="flex justify-between items-center p-2 bg-base-200 hover:bg-base-300 rounded-lg transition-colors text-sm"
            >
              <span class="font-medium truncate mr-2">{{ member.name }}</span>
              <span class="badge badge-ghost badge-sm shrink-0">{{ member.hours }}h</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>