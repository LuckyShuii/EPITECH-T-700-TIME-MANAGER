<script setup lang="ts">
import { ref } from 'vue'
import EmployeeLayout from '../layout/EmployeeLayout.vue'
import ClockWidget from '@/components/widget/ClockWidget.vue'
import CalendarWidget from '@/components/widget/CalendarWidget.vue'
import ClockButton from '@/components/ClockButton.vue'
import TeamManagementModal from '@/components/Modal/TeamManagementModal.vue'
import { useAuthStore } from '@/store/AuthStore'
import { storeToRefs } from 'pinia'

import TeamPresenceWidget from '@/components/widget/TeamPresenceWidget.vue'


const authStore = useAuthStore()
const { clockInTime, sessionStatus } = storeToRefs(authStore)

const isTeamViewModalOpen = ref(false)
</script>

<template>
  <EmployeeLayout>
    <!-- Position 1 : Haut gauche -->
    <template #widget-1>
      <TeamPresenceWidget />
    </template>

    <!-- Position 3 : Clock (bas gauche) -->
    <template #clock>
      <ClockWidget :clockInTime="clockInTime" :status="sessionStatus" />
      <ClockButton />
    </template>

    <!-- Position 4 : Calendar (haut centre - 2 colonnes) -->
    <template #calendar>
      <div class="bg-blue-100 p-6 rounded h-full">
        <CalendarWidget />
      </div>
    </template>

    <!-- Position 5 : Widget 2 (bas centre gauche) -->
    <template #widget-2>
      <div class="bg-green-100 p-6 rounded h-full">Widget 2</div>
    </template>

    <!-- Position 6 : Widget 3 (bas centre droite) -->
    <template #widget-3>
      <div class="bg-yellow-100 p-6 rounded h-full">Widget 3</div>
    </template>

    <!-- Position 7 : Voir mon équipe (haut droite) -->
    <template #widget-4>
      <button @click="isTeamViewModalOpen = true"
        class="h-full w-full bg-gradient-to-br from-primary-500 to-secondary-500 hover:shadow-card-hover text-white rounded-3xl shadow-card transition-all duration-300 flex flex-col items-center justify-center gap-4 group cursor-pointer">
        <div class="text-4xl group-hover:scale-110 transition-transform duration-300">👥</div>
        <p class="font-bold text-base">Voir mon équipe</p>
      </button>
    </template>

    <!-- Position 8 : Widget 5 (bas droite) -->
    <template #widget-5>
      <div class="bg-orange-100 p-6 rounded h-full">Widget 5</div>
    </template>
  </EmployeeLayout>

  

  <!-- Modal en dehors du layout -->
  <TeamManagementModal v-model="isTeamViewModalOpen" />
</template>
