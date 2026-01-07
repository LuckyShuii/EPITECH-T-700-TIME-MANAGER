<script setup lang="ts">
import { computed } from 'vue'
import type { WorkingTimeIndividualDisplay } from '@/types/Kpi'
import { useKpiStore } from '@/store/KpiStore'

interface Props {
  data: WorkingTimeIndividualDisplay | null
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
  return props.data !== null && props.data.total_time >= 0
})

// Convertir les minutes en heures
const totalHours = computed(() => {
  if (!props.data) return 0
  return Math.floor(props.data.total_time / 60)
})

const totalMinutes = computed(() => {
  if (!props.data) return 0
  return props.data.total_time % 60
})

const formattedTotal = computed(() => {
  if (totalMinutes.value === 0) {
    return `${totalHours.value}h`
  }
  return `${totalHours.value}h ${totalMinutes.value}m`
})

// Moyenne par jour (5 jours)
const averagePerDay = computed(() => {
  if (!props.data) return 0
  return (props.data.total_time / 60 / 5).toFixed(1)
})

// Différence avec la semaine précédente
const differenceHours = computed(() => {
  if (!props.data || props.data.difference === undefined) return 0
  return Math.floor(props.data.difference / 60)
})

const differenceMinutes = computed(() => {
  if (!props.data || props.data.difference === undefined) return 0
  return props.data.difference % 60
})

const formattedDifference = computed(() => {
  if (differenceHours.value === 0 && differenceMinutes.value === 0) {
    return 'Identique'
  }
  const sign = (props.data?.difference ?? 0) >= 0 ? '+' : ''
  if (differenceMinutes.value === 0) {
    return `${sign}${differenceHours.value}h`
  }
  return `${sign}${differenceHours.value}h ${differenceMinutes.value}m`
})

const trendColor = computed(() => {
  if (!props.data?.difference) return 'opacity-60'
  if (props.data.difference > 0) return 'text-green-700'
  if (props.data.difference < 0) return 'text-red-700'
  return 'opacity-60'
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
    <!-- Header avec navigation -->
    <div class="flex items-center justify-between mb-6 pb-4">
      <button
        @click="handlePreviousWeek"
        class="text-2xl font-bold hover:opacity-70 transition"
        title="Semaine précédente"
      >
        &lt;
      </button>
      <h2 class="text-lg font-bold">Heures travaillées</h2>
      <button
        @click="handleCurrentWeek"
        class="text-lg hover:opacity-70 transition"
        title="Revenir à aujourd'hui"
      >
        ↻
      </button>
    </div>

    <!-- Label semaine -->
    <p class="text-sm opacity-75 mb-6">{{ weekLabel }}</p>

    <!-- Contenu -->
    <div v-if="loading" class="flex-1 flex justify-center items-center">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="!hasValidData" class="flex-1 flex flex-col justify-center items-center opacity-60">
      <p class="text-sm">Aucune donnée</p>
    </div>

    <div v-else class="flex-1 flex flex-col justify-between">
      <!-- Total principal -->
      <div class="space-y-4">
        <div class="border-2 border-black p-4">
          <p class="text-xs opacity-75 mb-2">TOTAL CETTE SEMAINE</p>
          <p class="text-5xl font-bold font-mono">{{ formattedTotal }}</p>
        </div>

        <!-- Stats -->
        <div class="grid grid-cols-2 gap-4">
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

      <!-- Footer dates -->
      <div class="text-xs opacity-60 text-center mt-4 pt-4">
        <p>Du {{ data.start_date }} au {{ data.end_date }}</p>
      </div>
    </div>
  </div>
</template>