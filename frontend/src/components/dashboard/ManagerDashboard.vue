<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import ManagerLayout from '../layout/ManagerLayout.vue'
import { useEditModeStore } from '@/store/EditModeStore'
import ClockWidget from '@/components/widget/ClockWidget.vue'
import CalendarWidget from '@/components/widget/CalendarWidget.vue'
import ClockButton from '@/components/ClockButton.vue'
import { useAuthStore } from '@/store/AuthStore'
import { storeToRefs } from 'pinia'
import TeamManagementModal from '@/components/Modal/TeamManagementModal.vue'
import TeamPresenceWidget from '@/components/widget/TeamPresenceWidget.vue'

// Import des KPI cards PERTINENTS pour Manager
import TeamWorkingTimeCard from '@/components/kpi/cards/TeamWorkingTimeCard.vue'
import ClockingTimeCard from '@/components/kpi/cards/ClockingTimeCard.vue'

// Import des mock data
import {
  mockTeamWorkingTime,
  mockClockingTime
} from '@/mocks/kpiMockData'

const authStore = useAuthStore()
const { clockInTime, sessionStatus } = storeToRefs(authStore)
const editModeStore = useEditModeStore()
const isTeamViewModalOpen = ref(false)

onMounted(() => {
  editModeStore.setCurrentDashboard('manager')
})

onUnmounted(() => {
  editModeStore.reset()
})

const TeamViewModal = () => {
  isTeamViewModalOpen.value = true
}

// Handler pour les KPI
const handleKpiDetails = (data: any) => {
  console.log('KPI détails:', data)
}

// État de chargement des KPI
const kpiLoading = ref(false)
</script>

<template>
  <ManagerLayout>
    <!-- Clock + Boutons -->
    <template #clock>
      <ClockWidget
        :clockInTime="clockInTime"
        :status="sessionStatus"
      />
      <ClockButton />
    </template>

    <!-- KPI: Heures de pointage équipe -->
    <template #kpi-stats>
      <ClockingTimeCard 
        :data="mockClockingTime"
        :loading="kpiLoading"
        @view-details="handleKpiDetails"
      />
    </template>

    <!-- Calendrier -->
    <template #calendar>
      <div class="bg-blue-100 p-6 rounded h-full">
        <CalendarWidget />
      </div>
    </template>

    <!-- KPI: Travail hebdomadaire équipe -->
    <template #team-view>
      <TeamWorkingTimeCard 
        :data="mockTeamWorkingTime"
        :loading="kpiLoading"
        @view-details="handleKpiDetails"
      />
    </template>

    <!-- Widget présence équipe temps réel -->
    <template #team-presence>
      <TeamPresenceWidget />
    </template>

    <!-- Bouton: Voir mon équipe -->
    <template #modal-team>
      <button @click="isTeamViewModalOpen = true"
        class="brutal-btn brutal-btn-primary h-full w-full flex flex-col items-center justify-center gap-4">
        <p class="font-bold text-base">Voir mon équipe</p>
      </button>
    </template>
  </ManagerLayout>

  <!-- Modal en dehors du layout -->
  <TeamManagementModal v-model="isTeamViewModalOpen" />
</template>