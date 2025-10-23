<script setup lang="ts">
import { ref } from 'vue'
import ManagerLayout from '../layout/ManagerLayout.vue'
import ClockWidget from '@/components/widget/ClockWidget.vue'
import CalendarWidget from '@/components/widget/CalendarWidget.vue'
import ClockButton from '@/components/ClockButton.vue'
import { useAuthStore } from '@/store/AuthStore'
import { storeToRefs } from 'pinia'
import TeamManagementModal from '@/components/Modal/TeamManagementModal.vue'

import TeamPresenceWidget from '@/components/widget/TeamPresenceWidget.vue'

const authStore = useAuthStore()
const { clockInTime, sessionStatus } = storeToRefs(authStore)



const isTeamViewModalOpen = ref(false)

const TeamViewModal = () => {
  isTeamViewModalOpen.value = true
}

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

    <!-- KPI Stats -->
    <template #kpi-stats>
      <div class="bg-green-100 p-6 rounded h-full">
        <button class="btn btn-primary w-full">📊 Rapport</button>
      </div>
    </template>

    <template #team-view>
      <div class="bg-orange-100 p-6 rounded h-full"></div>

    </template>

    <!-- Calendrier -->
    <template #calendar>
      <div class="bg-blue-100 p-6 rounded h-full">
        <CalendarWidget />
      </div>
    </template>

    <!-- Présence équipe -->
    <template #modal-team>
      <button 
        @click="isTeamViewModalOpen = true"
        class="h-full w-full bg-gradient-to-br from-primary-500 to-secondary-500 hover:shadow-card-hover text-white rounded-3xl shadow-card transition-all duration-300 flex flex-col items-center justify-center gap-4 group cursor-pointer"
      >
        <div class="text-4xl group-hover:scale-110 transition-transform duration-300">👥</div>
        <p class="font-bold text-base">Voir mon équipe</p>
      </button>
    </template>

    <!-- KPI Carousel -->
    <template #team-presence>
      <TeamPresenceWidget />
    </template>

    <!-- Bouton rapport -->
    <template #report-button>
      <div class="bg-orange-100 p-6 rounded h-full">
      </div>
    </template>
  </ManagerLayout>

  <!-- Modal en dehors du layout -->
  <TeamManagementModal v-model="isTeamViewModalOpen" />
</template>