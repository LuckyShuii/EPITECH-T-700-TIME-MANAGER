<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'  // ← AJOUTE onMounted, onUnmounted
import { useEditModeStore } from '@/store/EditModeStore'  // ← AJOUTE
import AdminLayout from '@/components/layout/AdminLayout.vue'
import CalendarWidget from '@/components/widget/CalendarWidget.vue'
import BaseModal from '@/components/Modal/BaseModal.vue'
import RegisterForm from '@/components/RegisterForm.vue'
import StaffSettingsModal from '@/components/Modal/StaffSettingsModal.vue'
import TeamManagementAdminModal from '@/components/Modal/TeamManagementAdminModal.vue'

// AJOUTE CES LIGNES ↓
const editModeStore = useEditModeStore()

// Enregistre que ce dashboard est actif
onMounted(() => {
  editModeStore.setCurrentDashboard('admin')
})

// Nettoie quand on quitte le dashboard
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
</script>

<template>
  <AdminLayout>
    <!-- Le reste du template reste identique -->
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

    <template #kpi-history>
      <div class="bg-orange-100 p-4 rounded h-full flex items-center justify-center">
        <p class="text-sm font-medium">📈 KPI historique</p>
      </div>
    </template>

    <template #calendar>
      <div class="bg-white p-6 rounded h-full">
        <CalendarWidget />
      </div>
    </template>

    <template #widget-6>
      <div class="bg-gray-100 p-6 rounded h-full flex items-center justify-center">
        <p class="text-gray-500">Widget 6 - À définir</p>
      </div>
    </template>

    <template #remote-absence>
      <div class="bg-indigo-100 p-4 rounded h-full flex items-center justify-center">
        <p class="text-sm font-medium">🏠 TT / Absent</p>
      </div>
    </template>

    <template #manager-report>
      <div class="bg-yellow-100 p-4 rounded h-full flex items-center justify-center">
        <p class="text-sm font-medium">👔 Rapport manager</p>
      </div>
    </template>
  </AdminLayout>

  <!-- Modal d'ajout d'employé -->
  <BaseModal v-model="isAddEmployeeModalOpen" title="Créer un nouvel employé" size="lg">
    <RegisterForm @success="handleEmployeeCreated" @cancel="closeAddEmployeeModal" />
  </BaseModal>
  
  <!-- Modal paramétrage effectif -->
  <StaffSettingsModal v-model="isStaffSettingsModalOpen" />

  <!-- Modal gestion des équipes -->
  <TeamManagementAdminModal v-model="isTeamManagementModalOpen" />
</template>