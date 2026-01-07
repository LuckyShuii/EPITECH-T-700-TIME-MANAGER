<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useEditModeStore } from '@/store/editModeStore'
import { useKpiStore } from '@/store/KpiStore'
import { storeToRefs } from 'pinia'
import ManagerLayout from '@/components/layout/ManagerLayout.vue'
import CalendarWidget from '@/components/widget/CalendarWidget.vue'
import ClockWidget from '@/components/widget/ClockWidget.vue'
import ClockButton from '@/components/ClockButton.vue'
import TeamPresenceWidget from '@/components/widget/TeamPresenceWidget.vue'
import TeamManagementModal from '@/components/Modal/TeamManagementModal.vue'
import TeamWorkingTimeCard from '@/components/KPI/cards/TeamWorkingTimeCard.vue'
import { useAuthStore } from '@/store/AuthStore'

const authStore = useAuthStore()
const editModeStore = useEditModeStore()
const kpiStore = useKpiStore()

const { clockInTime, sessionStatus } = storeToRefs(authStore)
const { currentTeam, loading, weekDisplayLabel } = storeToRefs(kpiStore)

const isTeamViewModalOpen = ref(false)

onMounted(async () => {
  editModeStore.setCurrentDashboard('manager')
  
  try {
    await kpiStore.fetchManagerTeams()
    await kpiStore.fetchWorkingTimeTeam()
  } catch (error) {
    console.error('Erreur lors du chargement des KPI:', error)
  }
})

onUnmounted(() => {
  editModeStore.reset()
})
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

    <!-- Calendrier -->
    <template #calendar>
      <div class="bg-blue-100 p-6 rounded h-full">
        <CalendarWidget />
      </div>
    </template>

    <!-- KPI: Travail hebdomadaire équipe -->
    <template #team-view>
      <TeamWorkingTimeCard
        :data="currentTeam"
        :loading="kpiStore.loading['workingTimeTeam']"
        :weekLabel="weekDisplayLabel"
      />
    </template>

    <!-- Widget présence équipe temps réel -->
    <template #team-presence>
      <TeamPresenceWidget />
    </template>

    <!-- Bouton: Voir mon équipe -->
    <template #modal-team>
      <button 
        @click="isTeamViewModalOpen = true"
        class="brutal-btn brutal-btn-primary h-full w-full flex flex-col items-center justify-center gap-4"
      >
        <p class="font-bold text-base">Voir mon équipe</p>
      </button>
    </template>
  </ManagerLayout>

  <!-- Modal en dehors du layout -->
  <TeamManagementModal v-model="isTeamViewModalOpen" />
</template>