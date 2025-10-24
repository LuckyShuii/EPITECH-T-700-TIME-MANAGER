<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/store/AuthStore'
import type { Team } from '@/types/Team'
import API from '@/services/API'

const authStore = useAuthStore()

// État : équipes du manager
const teams = ref<Team[]>([])
const currentTeamIndex = ref(0)
const isLoading = ref(false)

// Intervalle de refresh
let refreshInterval: NodeJS.Timeout | null = null

// Computed : équipe actuellement affichée
const currentTeam = computed(() => {
    return teams.value[currentTeamIndex.value] || null
})

// Computed : afficher le carrousel (si plusieurs équipes)
const showCarousel = computed(() => {
    return teams.value.length > 1
})

// Charger les équipes du manager
// Charger les équipes (adaptée selon manager ou employé)
const loadTeams = async () => {
    if (!authStore.user?.user_uuid) return

    isLoading.value = true
    try {
        // Récupérer les infos du user
        const userResponse = await API.userAPI.getUserSpecific(authStore.user.user_uuid)

        // Si manager : filtrer les équipes où is_manager = true
        // Si employé : prendre toutes ses équipes
        const userTeams = userResponse.data.teams
        const teamsToLoad = authStore.isManager
            ? userTeams.filter((t: any) => t.is_manager)
            : userTeams

        // Si aucune équipe, on arrête
        if (teamsToLoad.length === 0) {
            teams.value = []
            return
        }

        // Pour chaque équipe, récupérer les détails avec les membres et leur work_session_status
        const teamPromises = teamsToLoad.map((teamInfo: any) =>
            API.teamAPI.getTeam(teamInfo.team_uuid)
        )

        const teamResponses = await Promise.all(teamPromises)
        teams.value = teamResponses.map(response => response.data)

    } catch (error) {
        console.error('Erreur lors du chargement des équipes:', error)
    } finally {
        isLoading.value = false
    }
}
// Navigation carrousel
const goToPreviousTeam = () => {
    if (currentTeamIndex.value > 0) {
        currentTeamIndex.value--
    } else {
        currentTeamIndex.value = teams.value.length - 1
    }
}

const goToNextTeam = () => {
    if (currentTeamIndex.value < teams.value.length - 1) {
        currentTeamIndex.value++
    } else {
        currentTeamIndex.value = 0
    }
}

// Fonction pour obtenir la classe de couleur selon le statut
const getStatusClass = (status: string) => {
    switch (status) {
        case 'active':
            return 'bg-success'
        case 'paused':
            return 'bg-warning'
        case 'no_active_session':
        default:
            return 'bg-error'
    }
}

// Fonction pour obtenir le texte du statut
const getStatusText = (status: string) => {
    switch (status) {
        case 'active':
            return 'Présent'
        case 'paused':
            return 'En pause'
        case 'no_active_session':
        default:
            return 'Absent'
    }
}

// Lifecycle
onMounted(() => {
  // Ne charger que si pointé
  if (authStore.isClockedIn) {
    loadTeams()
    
    // Refresh automatique toutes les 30 secondes
    refreshInterval = setInterval(() => {
      if (authStore.isClockedIn) {  // Vérifier à chaque refresh
        loadTeams()
      }
    }, 30000)
  }
})

onUnmounted(() => {
    if (refreshInterval) {
        clearInterval(refreshInterval)
    }
})
</script>

<template>
  <!-- Message si pas pointé -->
  <div 
    v-if="!authStore.isClockedIn" 
    class="h-full flex flex-col items-center justify-center bg-base-100 rounded-3xl shadow-card p-6"
  >
    <div class="text-center opacity-70">
      <p class="font-semibold text-sm mb-2">Accès restreint</p>
      <p class="text-xs">Pointez pour voir les personnes présentes de votre équipe</p>
    </div>
  </div>

  <!-- Contenu normal si pointé -->
  <div v-else class="h-full flex flex-col bg-base-100 rounded-3xl shadow-card p-6">
    <!-- Header avec titre et carrousel -->
    <div class="flex items-center justify-between mb-4">
      <!-- Flèche gauche (si plusieurs équipes) -->
      <button 
        v-if="showCarousel"
        @click="goToPreviousTeam"
        class="btn btn-ghost btn-sm btn-circle"
      >
        ←
      </button>
      
      <!-- Titre : nom de l'équipe -->
      <div class="flex-1 text-center">
        <h3 class="font-bold text-lg">
          {{ currentTeam?.name || 'Mon équipe' }}
        </h3>
        <p v-if="showCarousel" class="text-xs opacity-70">
          {{ currentTeamIndex + 1 }} / {{ teams.length }}
        </p>
      </div>

      <!-- Flèche droite (si plusieurs équipes) -->
      <button 
        v-if="showCarousel"
        @click="goToNextTeam"
        class="btn btn-ghost btn-sm btn-circle"
      >
        →
      </button>
    </div>

    <!-- Indicateur de chargement -->
    <div v-if="isLoading && !currentTeam" class="flex-1 flex items-center justify-center">
      <span class="loading loading-spinner loading-md"></span>
    </div>

    <!-- Liste des membres -->
    <div v-else-if="currentTeam" class="flex-1 overflow-y-auto space-y-2">
      <div 
        v-for="member in currentTeam.team_members" 
        :key="member.user_uuid"
        class="flex items-center gap-3 p-3 bg-base-200 rounded-lg"
      >
        <!-- Indicateur de statut (pastille colorée) -->
        <div 
          class="w-3 h-3 rounded-full flex-shrink-0"
          :class="getStatusClass(member.work_session_status)"
          :title="getStatusText(member.work_session_status)"
        ></div>

        <!-- Info membre -->
        <div class="flex-1 min-w-0">
          <p class="font-medium text-sm truncate">
            {{ member.first_name }} {{ member.last_name }}
          </p>
          <p class="text-xs opacity-70">
            {{ getStatusText(member.work_session_status) }}
          </p>
        </div>

        <!-- Badge Manager -->
        <span 
          v-if="member.is_manager" 
          class="badge badge-xs badge-primary flex-shrink-0"
        >
          Manager
        </span>
      </div>

      <!-- Message si aucun membre -->
      <div v-if="currentTeam.team_members.length === 0" class="text-center py-8 opacity-50">
        <p class="text-sm">Aucun membre dans cette équipe</p>
      </div>
    </div>

    <!-- Message si aucune équipe -->
    <div v-else class="flex-1 flex items-center justify-center opacity-50">
      <p class="text-sm">
        {{ authStore.isManager 
          ? 'Vous n\'êtes manager d\'aucune équipe' 
          : 'Vous n\'êtes dans aucune équipe' 
        }}
      </p>
    </div>
  </div>
</template>