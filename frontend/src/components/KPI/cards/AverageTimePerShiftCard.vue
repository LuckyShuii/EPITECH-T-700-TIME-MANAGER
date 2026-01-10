<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useKpiStore } from '@/store/KpiStore'
import { storeToRefs } from 'pinia'

const kpiStore = useKpiStore()
const { currentTeam, loading, weekDisplayLabel, averageTimePerShift } = storeToRefs(kpiStore)

const selectedEmployeeUuid = ref<string>('')

const employees = computed(() => {
  return currentTeam.value?.team_members ?? []
})

const hasValidData = computed(() => {
  return averageTimePerShift.value && averageTimePerShift.value.user_uuid
})

const averageHours = computed(() => {
  if (!hasValidData.value) return 0
  return Math.floor(averageTimePerShift.value.average_time / 60)
})

const averageMinutes = computed(() => {
  if (!hasValidData.value) return 0
  return averageTimePerShift.value.average_time % 60
})

const formattedAverage = computed(() => {
  if (averageMinutes.value === 0) {
    return `${averageHours.value}h`
  }
  return `${averageHours.value}h ${averageMinutes.value}m`
})

const totalHours = computed(() => {
  if (!hasValidData.value) return 0
  return Math.floor(averageTimePerShift.value.total_time / 60)
})

const totalMinutes = computed(() => {
  if (!hasValidData.value) return 0
  return averageTimePerShift.value.total_time % 60
})

const formattedTotal = computed(() => {
  if (totalMinutes.value === 0) {
    return `${totalHours.value}h`
  }
  return `${totalHours.value}h ${totalMinutes.value}m`
})

const handleEmployeeSelect = async () => {
  if (!selectedEmployeeUuid.value) return
  
  try {
    await kpiStore.fetchAverageTimePerShift(selectedEmployeeUuid.value, true)
  } catch (error) {
    console.error('Erreur lors du chargement des données:', error)
  }
}

const handlePreviousWeek = () => {
  kpiStore.changeWeek('previous')
}

const handleNextWeek = () => {
  kpiStore.changeWeek('next')
}

const handleCurrentWeek = () => {
  kpiStore.changeWeek('current')
}

onMounted(async () => {
  if (employees.value.length > 0) {
    selectedEmployeeUuid.value = employees.value[0].user_uuid
    await handleEmployeeSelect()
  }
})
</script>

<template>
  <div class="bg-white border-2 border-black p-4 h-full flex flex-col dark:bg-white">
    <div class="flex items-center justify-between mb-4">
      <button
        @click="handlePreviousWeek"
        class="text-lg font-bold hover:opacity-70 transition dark:text-gray-950"
        title="Semaine précédente"
      >
        &lt;
      </button>
      <div class="text-center flex-1">
        <h2 class="text-lg font-bold dark:text-gray-950">MOYENNE PAR SHIFT</h2>
        <p class="text-xs opacity-75 mt-1 dark:text-gray-500">{{ weekDisplayLabel }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="handleCurrentWeek"
          class="text-lg hover:opacity-70 transition dark:text-gray-950"
          title="Revenir à aujourd'hui"
        >
          ↻
        </button>
        <button
          @click="handleNextWeek"
          class="text-lg font-bold hover:opacity-70 transition dark:text-gray-950"
          title="Semaine suivante"
        >
          &gt;
        </button>
      </div>
    </div>

    <div class="mb-4">
      <label class="text-xs font-bold opacity-75 dark:text-gray-950 block mb-2">EMPLOYÉ</label>
      <select
        v-model="selectedEmployeeUuid"
        @change="handleEmployeeSelect"
        class="w-full border-2 border-black p-2 font-bold dark:text-gray-950 dark:bg-white"
      >
        <option v-for="employee in employees" :key="employee.user_uuid" :value="employee.user_uuid">
          {{ employee.first_name }} {{ employee.last_name }}
        </option>
      </select>
    </div>

    <div v-if="loading['averageTimePerShift']" class="flex-1 flex justify-center items-center">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="!hasValidData" class="flex-1 flex flex-col justify-center items-center opacity-60 dark:text-gray-950">
      <div class="text-center space-y-2">
        <p class="text-xs font-bold">Aucune donnée</p>
      </div>
    </div>

    <div v-else class="flex-1 flex flex-col justify-between">
      <div class="space-y-2">
        <div class="border-2 border-black p-2">
          <p class="text-xs opacity-75 mb-1 dark:text-gray-500">MOYENNE PAR SHIFT</p>
          <p class="text-3xl font-bold font-mono dark:text-gray-950">{{ formattedAverage }}</p>
        </div>

        <div class="border border-black p-2">
          <p class="text-xs opacity-75 mb-1 dark:text-gray-500">Nombre de shifts</p>
          <p class="text-xl font-bold font-mono dark:text-gray-950">{{ averageTimePerShift?.shift_count ?? 0 }}</p>
        </div>

        <div class="border border-black p-2">
          <p class="text-xs opacity-75 mb-1 dark:text-gray-500">Total heures</p>
          <p class="text-xl font-bold font-mono dark:text-gray-950">{{ formattedTotal }}</p>
        </div>
      </div>

      <div class="text-xs opacity-60 text-center mt-2 pt-2 dark:text-gray-500">
        <p v-if="averageTimePerShift?.date_range">Du {{ averageTimePerShift.date_range.start_date }} au {{ averageTimePerShift.date_range.end_date }}</p>
      </div>
    </div>
  </div>
</template>