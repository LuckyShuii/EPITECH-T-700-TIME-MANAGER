<script setup lang="ts">
import { computed } from 'vue'
import { useKpiStore } from '@/store/KpiStore'
import { storeToRefs } from 'pinia'

const kpiStore = useKpiStore()
const { currentTeam, loading, weekDisplayLabel } = storeToRefs(kpiStore)

const hasValidData = computed(() => {
  return currentTeam.value !== null && currentTeam.value.total_time >= 0
})

// Convertir les minutes en heures
const totalHours = computed(() => {
  if (!currentTeam.value) return 0
  return Math.floor(currentTeam.value.total_time / 60)
})

const totalMinutes = computed(() => {
  if (!currentTeam.value) return 0
  return currentTeam.value.total_time % 60
})

const formattedTotal = computed(() => {
  if (totalMinutes.value === 0) {
    return `${totalHours.value}h`
  }
  return `${totalHours.value}h ${totalMinutes.value}m`
})

// Moyenne par jour (5 jours)
const averagePerDay = computed(() => {
  if (!currentTeam.value) return 0
  return (currentTeam.value.total_time / 60 / 5).toFixed(1)
})

// Différence avec la semaine précédente
const differenceHours = computed(() => {
  if (!currentTeam.value || !(currentTeam.value as any).difference === undefined) return 0
  return Math.floor((currentTeam.value as any).difference / 60)
})

const differenceMinutes = computed(() => {
  if (!currentTeam.value || !(currentTeam.value as any).difference === undefined) return 0
  return (currentTeam.value as any).difference % 60
})

const formattedDifference = computed(() => {
  if (differenceHours.value === 0 && differenceMinutes.value === 0) {
    return 'Identique'
  }
  const sign = ((currentTeam.value as any)?.difference ?? 0) >= 0 ? '+' : ''
  if (differenceMinutes.value === 0) {
    return `${sign}${differenceHours.value}h`
  }
  return `${sign}${differenceHours.value}h ${differenceMinutes.value}m`
})

const trendColor = computed(() => {
  if (!(currentTeam.value as any)?.difference) return 'opacity-60'
  if ((currentTeam.value as any).difference > 0) return 'text-green-700'
  if ((currentTeam.value as any).difference < 0) return 'text-red-700'
  return 'opacity-60'
})

const handlePreviousTeam = () => {
  kpiStore.goToPreviousTeam()
}

const handleNextTeam = () => {
  kpiStore.goToNextTeam()
}

const handlePreviousWeek = () => {
  kpiStore.changeWeek('previous')
}

const handleCurrentWeek = () => {
  kpiStore.changeWeek('current')
}
</script>

<template>
  <div class="bg-white border-2 border-black p-6 h-full flex flex-col">
    <!-- Header avec navigation équipe et semaine -->
    <div class="flex items-center justify-between mb-6 pb-4">
      <button
        @click="handlePreviousTeam"
        class="text-2xl font-bold hover:opacity-70 transition  dark:text-gray-500"
        title="Équipe précédente"
      >
        &lt;
      </button>
      <div class="text-center flex-1">
        <h2 class="text-lg font-bold dark:text-gray-950">{{ currentTeam?.team_name || 'Aucune équipe' }}</h2>
        <p class="text-xs opacity-75 mt-1 dark:text-gray-500">{{ weekDisplayLabel }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="handleCurrentWeek"
          class="text-lg hover:opacity-70 transition dark:text-gray-500"
          title="Revenir à aujourd'hui"
        >
          ↻
        </button>
        <button
          @click="handleNextTeam"
          class="text-2xl font-bold hover:opacity-70 transition dark:text-gray-500"
          title="Équipe suivante"
        >
          &gt;
        </button>
      </div>
    </div>

    <!-- Contenu -->
    <div v-if="loading['workingTimeTeam']" class="flex-1 flex justify-center items-center">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="!hasValidData" class="flex-1 flex flex-col justify-center items-center opacity-60  dark:text-gray-500">
      <p class="text-sm">Aucune donnée</p>
    </div>

    <div v-else class="flex-1 flex flex-col justify-between  dark:text-gray-500">
      <!-- Total principal -->
      <div class="space-y-4">
        <div class="border-2 border-black p-4">
          <p class="text-xs opacity-75 mb-2">TOTAL CETTE SEMAINE</p>
          <p class="text-5xl font-bold font-mono">{{ formattedTotal }}</p>
        </div>

        <!-- Stats -->
        <div class="grid grid-cols-2 gap-4  dark:text-gray-500">
          <div class="border border-black p-3">
            <p class="text-xs opacity-75 mb-1">Moyenne/jour</p>
            <p class="text-3xl font-bold font-mono">{{ averagePerDay }}h</p>
          </div>

          <div class="border border-black p-3" :class="trendColor">
            <p class="text-xs opacity-75 mb-1">vs Sem. précédente</p>
            <p class="text-2xl font-bold font-mono">{{ formattedDifference }}</p>
          </div>
        </div>
      </div>

      <!-- Liste des membres -->
      <div class="mt-4 pt-4 border-t-2 border-black">
        <p class="text-xs opacity-75 mb-2 font-bold">MEMBRES ({{ currentTeam?.members.length }})</p>
        <div class="space-y-1 max-h-32 overflow-y-auto">
          <div
            v-for="member in currentTeam?.members"
            :key="member.user_uuid"
            class="flex justify-between items-center text-xs p-2 bg-gray-100"
          >
            <span class="font-bold">{{ member.first_name }} {{ member.last_name }}</span>
            <span class="font-mono">{{ Math.floor(member.total_time / 60) }}h</span>
          </div>
        </div>
      </div>

      <!-- Footer dates -->
      <div class="text-xs opacity-60 text-center mt-4 pt-4">
        <p>Du {{ currentTeam?.start_date }} au {{ currentTeam?.end_date }}</p>
      </div>
    </div>
  </div>
</template>