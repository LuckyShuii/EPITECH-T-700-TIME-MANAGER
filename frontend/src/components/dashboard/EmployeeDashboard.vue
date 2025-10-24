<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useEditModeStore } from '@/store/EditModeStore'  
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

const editModeStore = useEditModeStore()

// Enregistre que ce dashboard est actif
onMounted(() => {
  editModeStore.setCurrentDashboard('employee')
  console.log('✅ Dashboard enregistré:', editModeStore.currentDashboard)  // ← AJOUTE
})

// Nettoie quand on quitte le dashboard
onUnmounted(() => {
  editModeStore.reset()
})


</script>

<template>
  <!-- RETIRE le ref="layoutRef" ici ↓ -->
  <EmployeeLayout>
    <!-- Widget presence -->
    <template #team-presence>
      <TeamPresenceWidget />
    </template>
    
    <!-- Clock -->
    <template #clock>
      <ClockWidget :clockInTime="clockInTime" :status="sessionStatus" />
      <ClockButton />
    </template>
    
    <!-- Calendar -->
    <template #calendar>
      <div class="bg-blue-100 p-6 rounded h-full">
        <CalendarWidget />
      </div>
    </template>
    
    <!-- Widget 2 -->
    <template #widget-2>
      <div class="bg-green-100 p-6 rounded h-full">Widget 2</div>
    </template>
    
    <!-- Widget 3 -->
    <template #widget-3>
      <div class="bg-yellow-100 p-6 rounded h-full">Widget 3</div>
    </template>
    
    <!-- Widget 4 - Bouton équipe -->
    <template #widget-4>
      <button @click="isTeamViewModalOpen = true"
        class="h-full w-full bg-gradient-to-br from-primary-500 to-secondary-500 hover:shadow-card-hover text-white rounded-3xl shadow-card transition-all duration-300 flex flex-col items-center justify-center gap-4 group cursor-pointer">
        <div class="text-4xl group-hover:scale-110 transition-transform duration-300">👥</div>
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
      <div class="bg-green-100 p-6 rounded h-full">Widget 2</div>
    </template>
  </EmployeeLayout>
  
  <TeamManagementModal v-model="isTeamViewModalOpen" />
</template>