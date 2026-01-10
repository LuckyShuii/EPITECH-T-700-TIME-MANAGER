<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useEditModeStore } from '@/store/EditModeStore'
import { useKpiStore } from '@/store/KpiStore'
import EmployeeLayout from '../layout/EmployeeLayout.vue'
import ClockWidget from '@/components/widget/ClockWidget.vue'
import CalendarWidget from '@/components/widget/CalendarWidget.vue'
import ClockButton from '@/components/ClockButton.vue'
import TeamManagementModal from '@/components/Modal/TeamManagementModal.vue'
import IndividualPauseWidget from '@/components/KPI/cards/IndividualPauseCard.vue'
import WorkingTimeIndividualCard from '@/components/KPI/cards/WorkingTimeIndividualCard.vue'
import { useAuthStore } from '@/store/AuthStore'
import { storeToRefs } from 'pinia'
import TeamPresenceWidget from '@/components/widget/TeamPresenceWidget.vue'

const authStore = useAuthStore()
const kpiStore = useKpiStore()

const { clockInTime, sessionStatus } = storeToRefs(authStore)
const { individualPause, workingTimeIndividual, loading, weekDisplayLabel } = storeToRefs(kpiStore)

const isTeamViewModalOpen = ref(false)
const editModeStore = useEditModeStore()

onMounted(async () => {
  editModeStore.setCurrentDashboard('employee')
  console.log('Dashboard enregistré:', editModeStore.currentDashboard)
  
  try {
    await kpiStore.fetchIndividualPause()
    await kpiStore.fetchWorkingTimeIndividual()
  } catch (error) {
    console.error('Erreur lors du chargement des KPI:', error)
  }
})

onUnmounted(() => {
  editModeStore.reset()
})
</script>

<template>
  <EmployeeLayout>
    <!-- Widget presence -->
    <template #team-presence>
      <TeamPresenceWidget />
    </template>

    <!-- Clock -->
    <template #clock>
      <!-- center horizontally -->
      <div class="flex justify-center flex-col items-center mb-4 h-full border-2 border-black p-4">
        <ClockWidget :clockInTime="clockInTime" :status="sessionStatus" />
        <ClockButton />
      </div>
    </template>

    <!-- Calendar -->
    <template #calendar>
        <CalendarWidget />
    </template>

    <!-- Widget 2 - Individual Pause -->
    <template #widget-2>
      <IndividualPauseWidget
        :data="individualPause"
        :loading="loading['individualPause']"
        :weekLabel="weekDisplayLabel"
      />
    </template>

    <!-- Widget 3 - Working Time Individual -->
    <template #widget-3>
      <WorkingTimeIndividualCard
        :data="workingTimeIndividual"
        :loading="loading['workingTimeIndividual']"
        :weekLabel="weekDisplayLabel"
      />
    </template>

    <!-- Widget 4 - Bouton équipe -->
    <template #widget-4>
      <button
        @click="isTeamViewModalOpen = true"
        class="brutal-btn brutal-btn-primary h-full w-full flex flex-col items-center justify-center gap-4"
      >
        <p class="font-bold text-base">Voir mon équipe</p>
      </button>
    </template>

    <!-- Extra widget -->
    <template #extra-widget>
      <div class="bg-purple-100 p-6 rounded h-full">Extra Widget</div>
    </template>

    <!-- Widget 5 -->
    <template #widget-5>
      <div class="bg-orange-100 p-6 rounded h-full">Widget 5</div>
    </template>

    <!-- Widget 6 -->
    <template #widget-6>
      <div class="bg-green-100 p-6 rounded h-full">Widget 6</div>
    </template>
  </EmployeeLayout>

  <TeamManagementModal v-model="isTeamViewModalOpen" />
</template>