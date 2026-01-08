<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useNotificationsStore } from '@/store/NotificationsStore'
import API from '@/services/API'
import type { KpiExportRequest } from '@/services/routers/ExportAPI'

const notificationsStore = useNotificationsStore()

const isLoading = ref(false)
const isLoadingData = ref(false)

const users = ref<any[]>([])
const teams = ref<any[]>([])

const formData = ref<KpiExportRequest>({
  start_date: '',
  end_date: '',
  kpi_type: 'work_session_user_weekly_total',
  uuid_to_search: ''
})

const kpiOptions = [
  { value: 'work_session_user_weekly_total', label: 'Total heures hebdomadaires par utilisateur' },
  { value: 'average_break_time', label: 'Temps de pause moyen' },
  { value: 'average_time_per_shift', label: 'Temps moyen par shift' },
  { value: 'presence_rate', label: 'Taux de présence' },
  { value: 'work_session_team_weekly_total', label: 'Total heures hebdomadaires par équipe' }
]

const needsTeam = computed(() => {
  return formData.value.kpi_type === 'work_session_team_weekly_total'
})

const loadData = async () => {
  isLoadingData.value = true
  try {
    const usersResponse = await API.userAPI.getAll()
    users.value = usersResponse.data || []
    
    const teamsResponse = await API.teamAPI.getAll()
    teams.value = teamsResponse.data || []
  } catch (error) {
    console.error('Erreur lors du chargement des données:', error)
    notificationsStore.addNotification({
      status: 'error',
      title: 'Erreur de chargement',
      description: 'Impossible de charger les utilisateurs et équipes'
    })
  } finally {
    isLoadingData.value = false
  }
}

onMounted(() => {
  loadData()
})

const handleExport = async () => {
  if (!formData.value.start_date || !formData.value.end_date || !formData.value.uuid_to_search) {
    notificationsStore.addNotification({
      status: 'error',
      title: 'Formulaire incomplet',
      description: 'Veuillez remplir tous les champs'
    })
    return
  }

  isLoading.value = true
  try {
    const response = await API.exportAPI.exportKpiData(formData.value)
    
    window.open(response.url, '_blank')
    
    notificationsStore.addNotification({
      status: 'success',
      title: 'Export réussi',
      description: 'Le fichier a été téléchargé'
    })
  } catch (error) {
    console.error('Erreur lors de l\'export:', error)
    notificationsStore.addNotification({
      status: 'error',
      title: 'Erreur d\'export',
      description: 'Impossible d\'exporter les données'
    })
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="space-y-6">
    <!-- Loading initial -->
    <div v-if="isLoadingData" class="flex justify-center items-center py-8">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <template v-else>
      <!-- Dates -->
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="label pt-0">
            <span class="label-text font-bold uppercase text-xs tracking-widest">Date de début</span>
          </label>
          <input v-model="formData.start_date" type="date" class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100" />
        </div>
        <div>
          <label class="label pt-0">
            <span class="label-text font-bold uppercase text-xs tracking-widest">Date de fin</span>
          </label>
          <input v-model="formData.end_date" type="date" class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100" />
        </div>
      </div>

      <!-- Type de KPI -->
      <div>
        <label class="label pt-0">
          <span class="label-text font-bold uppercase text-xs tracking-widest">Type de KPI</span>
        </label>
        <select v-model="formData.kpi_type" class="select select-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100">
          <option v-for="kpi in kpiOptions" :key="kpi.value" :value="kpi.value">
            {{ kpi.label }}
          </option>
        </select>
      </div>

      <!-- Dropdown Utilisateurs (caché si team est sélectionné) -->
      <div v-if="!needsTeam">
        <label class="label pt-0">
          <span class="label-text font-bold uppercase text-xs tracking-widest">Sélectionner un utilisateur</span>
        </label>
        <select v-model="formData.uuid_to_search" class="select select-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100">
          <option value="">-- Choisir un utilisateur --</option>
          <option v-for="user in users" :key="user.uuid" :value="user.uuid">
            {{ user.first_name }} {{ user.last_name }} ({{ user.username }})
          </option>
        </select>
      </div>

      <!-- Dropdown Équipes (caché si user est sélectionné) -->
      <div v-else>
        <label class="label pt-0">
          <span class="label-text font-bold uppercase text-xs tracking-widest">Sélectionner une équipe</span>
        </label>
        <select v-model="formData.uuid_to_search" class="select select-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100">
          <option value="">-- Choisir une équipe --</option>
          <option v-for="team in teams" :key="team.uuid" :value="team.uuid">
            {{ team.name }}
          </option>
        </select>
      </div>

      <!-- Boutons -->
      <div class="flex gap-2 pt-4 border-t-2 border-black">
        <button @click="handleExport" :disabled="isLoading" class="brutal-btn brutal-btn-primary flex-1">
          {{ isLoading ? 'Chargement...' : 'Exporter' }}
        </button>
      </div>
    </template>
  </div>
</template>