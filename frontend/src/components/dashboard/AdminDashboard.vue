<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useEditModeStore } from '@/store/EditModeStore'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import CalendarWidget from '@/components/widget/CalendarWidget.vue'
import BaseModal from '@/components/Modal/BaseModal.vue'
import RegisterForm from '@/components/RegisterForm.vue'
import StaffSettingsModal from '@/components/Modal/StaffSettingsModal.vue'
import TeamManagementAdminModal from '@/components/Modal/TeamManagementAdminModal.vue'

// Import des KPI cards
import WorkingTimeCard from '@/components/kpi/cards/WorkingTimeCard.vue'
import WeeklyProgressCard from '@/components/kpi/cards/WeeklyProgressCard.vue'
import ShiftAverageCard from '@/components/kpi/cards/ShiftAverageCard.vue'

// Import des mock data
import {
  mockWorkingTimeIndividual,
  mockWeeklyProgress,
  mockShiftAverage
} from '@/mocks/kpiMockData'

const editModeStore = useEditModeStore()

onMounted(() => {
  editModeStore.setCurrentDashboard('admin')
})

onUnmounted(() => {
  editModeStore.reset()
})

// Contrôle du modal d'ajout employé
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

// Contrôle du modal paramétrage effectif
const isStaffSettingsModalOpen = ref(false)

const openStaffSettingsModal = () => {
  isStaffSettingsModalOpen.value = true
}

// Contrôle du modal gestion des équipes
const isTeamManagementModalOpen = ref(false)

const openTeamManagementModal = () => {
  isTeamManagementModalOpen.value = true
}

// Handler pour les KPI (placeholder)
const handleKpiDetails = (data: any) => {
  console.log('KPI détails:', data)
  // TODO: Ouvrir un modal ou naviguer
}

// État de chargement des KPI (pour test)
const kpiLoading = ref(false)
</script>

<template>
  <AdminLayout>
    <template #add-employee>
      <button @click="openAddEmployeeModal"
        class="h-full w-full bg-gradient-to-br from-primary-500 to-secondary-500 hover:shadow-card-hover text-white rounded-3xl shadow-card transition-all duration-300 flex flex-col items-center justify-center gap-4 group cursor-pointer">
        <div class="text-4xl group-hover:scale-110 transition-transform duration-300">➕</div>
        <p class="font-bold text-base">Nouvel employé</p>
      </button>
    </template>

    <template #staff-settings>
      <button @click="openStaffSettingsModal"
        class="h-full w-full bg-gradient-to-br from-purple-500 to-indigo-600 hover:shadow-card-hover text-white rounded-3xl shadow-card transition-all duration-300 flex flex-col items-center justify-center gap-4 group cursor-pointer">
        <div class="text-4xl group-hover:scale-110 transition-transform duration-300">⚙️</div>
        <p class="font-bold text-base">Paramétrage effectifs</p>
      </button>
    </template>

    <template #kpi-monthly>
      <button @click="openTeamManagementModal"
        class="h-full w-full bg-gradient-to-br from-green-500 to-teal-600 hover:shadow-card-hover text-white rounded-3xl shadow-card transition-all duration-300 flex flex-col items-center justify-center gap-4 group cursor-pointer">
        <div class="text-4xl group-hover:scale-110 transition-transform duration-300">👥</div>
        <p class="font-bold text-base">Gestion des équipes</p>
      </button>
    </template>

    <!-- KPI Cards dans les slots disponibles -->
    <template #widget-6>
      <WorkingTimeCard 
        :data="mockWorkingTimeIndividual"
        :loading="kpiLoading"
        @view-details="handleKpiDetails"
      />
    </template>

    <!-- CALENDAR RESTE ICI -->
    <template #calendar>
      <div class="bg-white p-6 rounded h-full">
        <CalendarWidget />
      </div>
    </template>

    <template #manager-report>
      <WeeklyProgressCard 
        :data="mockWeeklyProgress"
        :loading="kpiLoading"
        @view-details="handleKpiDetails"
      />
    </template>

    <template #remote-absence>
      <ShiftAverageCard 
        :data="mockShiftAverage"
        :loading="kpiLoading"
        @view-details="handleKpiDetails"
      />
    </template>

    <template #kpi-history>
      <div class="bg-yellow-100 p-4 rounded h-full flex items-center justify-center">
        <p class="text-sm font-medium">👔 Rapport manager</p>
      </div>
    </template>
  </AdminLayout>

  <!-- Modals -->
  <BaseModal v-model="isAddEmployeeModalOpen" title="Créer un nouvel employé" size="lg">
    <RegisterForm @success="handleEmployeeCreated" @cancel="closeAddEmployeeModal" />
  </BaseModal>
  
  <StaffSettingsModal v-model="isStaffSettingsModalOpen" />

  <TeamManagementAdminModal v-model="isTeamManagementModalOpen" />
</template>