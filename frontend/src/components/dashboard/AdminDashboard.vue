<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useEditModeStore } from '@/store/EditModeStore'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import CalendarWidget from '@/components/widget/CalendarWidget.vue'
import BaseModal from '@/components/Modal/BaseModal.vue'
import RegisterForm from '@/components/RegisterForm.vue'
import StaffSettingsModal from '@/components/Modal/StaffSettingsModal.vue'
import TeamManagementAdminModal from '@/components/Modal/TeamManagementAdminModal.vue'

import ClockingTimeCard from '@/components/kpi/cards/ClockingTimeCard.vue'
import WeeklyProgressCard from '@/components/kpi/cards/WeeklyProgressCard.vue'
import ShiftAverageCard from '@/components/kpi/cards/ShiftAverageCard.vue'
import TeamPauseCard from '@/components/kpi/cards/TeamPauseCard.vue'
import ShiftPauseCard from '@/components/kpi/cards/ShiftPauseCard.vue'

import {
  mockWeeklyProgress,
  mockShiftAverage,
  mockClockingTime,
  mockTeamPause,
  mockShiftPause
} from '@/mocks/kpiMockData'

const editModeStore = useEditModeStore()

onMounted(() => {
  editModeStore.setCurrentDashboard('admin')
})

onUnmounted(() => {
  editModeStore.reset()
})

const isAddEmployeeModalOpen = ref(false)
const openAddEmployeeModal = () => {
  isAddEmployeeModalOpen.value = true
}

const closeAddEmployeeModal = () => {
  isAddEmployeeModalOpen.value = false
}

const handleEmployeeCreated = () => {
  closeAddEmployeeModal()
}

const isStaffSettingsModalOpen = ref(false)
const openStaffSettingsModal = () => {
  isStaffSettingsModalOpen.value = true
}

const isTeamManagementModalOpen = ref(false)
const openTeamManagementModal = () => {
  isTeamManagementModalOpen.value = true
}

const handleKpiDetails = (data: any) => {
  console.log('KPI détails:', data)
}

const kpiLoading = ref(false)
</script>

<template>
  <AdminLayout>
    <!-- Bouton: Nouvel employé -->
    <template #add-employee>
      <button @click="openAddEmployeeModal" class="brutal-btn brutal-btn-primary w-full h-full flex flex-col items-center justify-center gap-4">
        <div class="text-2xl">+</div>
        <p class="font-bold">Nouvel employé</p>
      </button>
    </template>

    <!-- Bouton: Paramétrage effectifs -->
    <template #staff-settings>
      <button @click="openStaffSettingsModal" class="brutal-btn brutal-btn-primary w-full h-full flex flex-col items-center justify-center gap-4">
        <div class="text-2xl">!</div>
        <p class="font-bold">Paramétrage effectifs</p>
      </button>
    </template>

    <!-- Bouton: Gestion des équipes -->
    <template #kpi-monthly>
      <button @click="openTeamManagementModal" class="brutal-btn brutal-btn-primary w-full h-full flex flex-col items-center justify-center gap-4">
        <div class="text-2xl">=</div>
        <p class="font-bold">Gestion des équipes</p>
      </button>
    </template>

    <!-- KPI: Temps de pointage moyen hebdo -->
    <template #kpi-history>
      <ClockingTimeCard 
        :data="mockClockingTime"
        :loading="kpiLoading"
        @view-details="handleKpiDetails"
      />
    </template>

    <!-- Calendrier -->
    <template #calendar>
      <div class="bg-white p-6 rounded h-full">
        <CalendarWidget />
      </div>
    </template>

    <!-- KPI: Pause moyenne équipe -->
    <template #widget-6>
      <TeamPauseCard 
        :data="mockTeamPause"
        :loading="kpiLoading"
        @view-details="handleKpiDetails"
      />
    </template>

    <!-- KPI: Répartition shift (pause) -->
    <template #widget-7>
      <ShiftPauseCard 
        :data="mockShiftPause"
        :loading="kpiLoading"
        @view-details="handleKpiDetails"
      />
    </template>

    <!-- KPI: Moyenne par shift -->
    <template #remote-absence>
      <ShiftAverageCard 
        :data="mockShiftAverage"
        :loading="kpiLoading"
        @view-details="handleKpiDetails"
      />
    </template>

    <!-- KPI: Progression hebdo -->
    <template #manager-report>
      <WeeklyProgressCard 
        :data="mockWeeklyProgress"
        :loading="kpiLoading"
        @view-details="handleKpiDetails"
      />
    </template>
  </AdminLayout>

  <!-- Modals -->
  <BaseModal v-model="isAddEmployeeModalOpen" title="Créer un nouvel employé">
    <RegisterForm @success="handleEmployeeCreated" @cancel="closeAddEmployeeModal" />
  </BaseModal>
  
  <StaffSettingsModal v-model="isStaffSettingsModalOpen" />

  <TeamManagementAdminModal v-model="isTeamManagementModalOpen" />
</template>