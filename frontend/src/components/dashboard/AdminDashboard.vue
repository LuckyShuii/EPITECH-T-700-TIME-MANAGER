<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useEditModeStore } from '@/store/editModeStore'
import { useKpiStore } from '@/store/KpiStore'
import { storeToRefs } from 'pinia'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import CalendarWidget from '@/components/widget/CalendarWidget.vue'
import BaseModal from '@/components/Modal/BaseModal.vue'
import RegisterForm from '@/components/RegisterForm.vue'
import StaffSettingsModal from '@/components/Modal/StaffSettingsModal.vue'
import TeamManagementAdminModal from '@/components/Modal/TeamManagementAdminModal.vue'
import TeamWorkingTimeCard from '@/components/KPI/cards/TeamWorkingTimeCard.vue'
import PresenceRateCard from '@/components/KPI/cards/PresenceRateCard.vue'
import { ArrowDownTrayIcon } from '@heroicons/vue/24/solid'

const editModeStore = useEditModeStore()
const kpiStore = useKpiStore()

const { currentTeam, weekDisplayLabel, presenceRate } = storeToRefs(kpiStore)

onMounted(async () => {
  editModeStore.setCurrentDashboard('admin')
  
  try {
    await kpiStore.fetchManagerTeams()
    await kpiStore.fetchWorkingTimeTeam()
    await kpiStore.fetchPresenceRate()
  } catch (error) {
    console.error('Erreur lors du chargement des KPI:', error)
  }
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
    <template #team-gestion>
      <button @click="openTeamManagementModal" class="brutal-btn brutal-btn-primary w-full h-full flex flex-col items-center justify-center gap-4">
        <div class="text-2xl">=</div>
        <p class="font-bold">Gestion des équipes</p>
      </button>
    </template>

    <!-- KPI: Travail hebdomadaire équipe -->
    <template #kpi-history>
      <TeamWorkingTimeCard
        :data="currentTeam"
        :loading="kpiStore.loading['workingTimeTeam']"
        :weekLabel="weekDisplayLabel"
      />
    </template>

    <!-- Calendrier -->
    <template #calendar>
      <div class="bg-white p-6 rounded h-full">
        <CalendarWidget />
      </div>
    </template>

    <!-- KPI: Taux de présence -->
    <template #presence-rate>
      <PresenceRateCard
        :data="presenceRate"
        :loading="kpiStore.loading['presenceRate']"
        :weekLabel="weekDisplayLabel"
      />
    </template>


    <!-- Export: Statistique des KPI-->
<template #export-button>
  <button @click="openTeamManagementModal" class="brutal-btn brutal-btn-primary w-full h-full flex flex-col items-center justify-center gap-4">
    <ArrowDownTrayIcon class="w-5 h-5" />
    <p class="font-bold">Export des statistiques</p>
  </button>
</template>
  </AdminLayout>

  <!-- Modals -->
  <BaseModal v-model="isAddEmployeeModalOpen" title="Créer un nouvel employé">
    <RegisterForm @success="handleEmployeeCreated" @cancel="closeAddEmployeeModal" />
  </BaseModal>
  <StaffSettingsModal v-model="isStaffSettingsModalOpen" />
  <TeamManagementAdminModal v-model="isTeamManagementModalOpen" />
</template>