<script setup lang="ts">
import { computed } from 'vue'
import { useKpiStore } from '@/store/KpiStore'
import { storeToRefs } from 'pinia'

const kpiStore = useKpiStore()
const { presenceRate, loading, weekDisplayLabel } = storeToRefs(kpiStore)

const hasValidData = computed(() => {
  return presenceRate.value && Array.isArray(presenceRate.value) && presenceRate.value.length > 0
})

const getPresenceColor = (rate: number) => {
  if (rate >= 90) return 'bg-green-600'
  if (rate >= 75) return 'bg-yellow-500'
  if (rate >= 60) return 'bg-orange-500'
  return 'bg-red-600'
};

const getPresenceStatusText = (rate: number) => {
  if (rate >= 90) return 'Excellent'
  if (rate >= 75) return 'Bon'
  if (rate >= 60) return 'Moyen'
  return 'Faible'
};
</script>

<template>
  <div class="bg-white border-2 border-black p-4 h-full flex flex-col dark:bg-white">
    <!-- Header -->
    <div class="border-b-2 border-black pb-3 mb-4">
      <h2 class="text-lg font-bold dark:text-gray-950">TAUX DE PRÉSENCE</h2>
      <p class="text-xs opacity-75 mt-1 dark:text-gray-500">{{ weekDisplayLabel }}</p>
    </div>

    <!-- Contenu -->
    <div v-if="loading['presenceRate']" class="flex-1 flex justify-center items-center">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="!hasValidData" class="flex-1 flex flex-col justify-center items-center opacity-60 dark:text-gray-950">
      <div class="text-center space-y-2">
        <p class="text-xs font-bold">Aucune donnée</p>
      </div>
    </div>

    <div v-else class="flex-1 overflow-y-auto space-y-3 pr-2">
      <div
        v-for="user in presenceRate"
        :key="user.user_uuid"
        class="border border-black p-3 space-y-2 dark:border-gray-700 dark:bg-gray-50"
      >
        <!-- Nom et pourcentage -->
        <div class="flex items-center justify-between">
          <p class="text-xs font-bold dark:text-gray-950">{{ user.first_name }} {{ user.last_name }}</p>
          <p class="text-sm font-mono font-bold dark:text-gray-950">{{ Math.round(user.presence_rate) }}%</p>
        </div>

        <!-- Barre de progression -->
        <div class="w-full bg-gray-300 border border-black h-3 dark:bg-gray-300">
          <div
            :class="getPresenceColor(user.presence_rate)"
            :style="{ width: `${Math.min(user.presence_rate, 100)}%` }"
            class="h-full border-r border-black transition-all duration-300"
          ></div>
        </div>

        <!-- Stats heures -->
        <div class="grid grid-cols-2 gap-2 text-xs">
          <div>
            <p class="opacity-75 dark:text-gray-600">Attendu</p>
            <p class="font-mono font-bold dark:text-gray-950">{{ Math.round(user.weekly_rate_expected * 10) / 10 }}h</p>
          </div>
          <div>
            <p class="opacity-75 dark:text-gray-600">Effectué</p>
            <p class="font-mono font-bold dark:text-gray-950">{{ Math.round(user.weekly_time_done * 10) / 10 }}h</p>
          </div>
        </div>

        <!-- Status badge -->
        <div class="text-center">
          <span class="inline-block border border-black px-2 py-1 text-xs font-bold dark:bg-gray-200 dark:text-gray-950">
            {{ getPresenceStatusText(user.presence_rate) }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>