<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/store/AuthStore'
import type { Team } from '@/types/Team'
import API from '@/services/API'

const authStore = useAuthStore()

const teams = ref<Team[]>([])
const currentTeamIndex = ref(0)
const isLoading = ref(false)

let refreshInterval: NodeJS.Timeout | null = null

const currentTeam = computed(() => {
    return teams.value[currentTeamIndex.value] || null
})

const showCarousel = computed(() => {
    return teams.value.length > 1
})

const loadTeams = async () => {
    if (!authStore.user?.user_uuid) return

    isLoading.value = true
    try {
        const userResponse = await API.userAPI.getUserSpecific(authStore.user.user_uuid)
        const userTeams = userResponse.data.teams
        const teamsToLoad = authStore.isManager
            ? userTeams.filter((t: any) => t.is_manager)
            : userTeams

        if (teamsToLoad.length === 0) {
            teams.value = []
            return
        }

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

const getStatusClass = (status: string) => {
    switch (status) {
        case 'active':
            return 'bg-green-600'
        case 'paused':
            return 'bg-yellow-600'
        case 'no_active_session':
        default:
            return 'bg-red-600'
    }
}

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

onMounted(() => {
  if (authStore.isClockedIn) {
    loadTeams()
    refreshInterval = setInterval(() => {
      if (authStore.isClockedIn) {
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
    class="h-full flex flex-col items-center justify-center border-2 border-black bg-white p-6"
  >
    <div class="text-center">
      <p class="font-bold uppercase text-xs tracking-widest mb-2">ACCÈS RESTREINT</p>
      <p class="text-xs font-bold">Pointez pour voir votre équipe</p>
    </div>
  </div>

  <!-- Contenu normal si pointé -->
  <div v-else class="h-full flex flex-col border-2 border-black bg-white p-6">
    <!-- Header avec titre et carrousel -->
    <div class="flex items-center justify-between mb-6 pb-4">
      <!-- Flèche gauche (si plusieurs équipes) -->
      <button 
        v-if="showCarousel"
        @click="goToPreviousTeam"
        class="w-8 h-8 flex items-center justify-center border-2 border-black font-bold hover:bg-black hover:text-white transition-none"
      >
        ←
      </button>
      
      <!-- Titre : nom de l'équipe -->
      <div class="flex-1 text-center">
        <h3 class="font-black uppercase tracking-wider">
          {{ currentTeam?.name || 'Mon équipe' }}
        </h3>
        <p v-if="showCarousel" class="text-xs font-bold mt-1">
          {{ currentTeamIndex + 1 }} / {{ teams.length }}
        </p>
      </div>

      <!-- Flèche droite (si plusieurs équipes) -->
      <button 
        v-if="showCarousel"
        @click="goToNextTeam"
        class="w-8 h-8 flex items-center justify-center border-2 border-black font-bold hover:bg-black hover:text-white transition-none"
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
        class="flex items-center gap-3 p-3 border-2 border-gray-400 hover:border-black hover:bg-gray-100 transition-none"
      >
        <!-- Indicateur de statut (carré coloré) -->
        <div 
          class="w-3 h-3 flex-shrink-0 border border-black"
          :class="getStatusClass(member.work_session_status)"
          :title="getStatusText(member.work_session_status)"
        ></div>

        <!-- Info membre -->
        <div class="flex-1 min-w-0">
          <p class="font-bold text-sm truncate">
            {{ member.first_name }} {{ member.last_name }}
          </p>
          <p class="text-xs font-bold">
            {{ getStatusText(member.work_session_status) }}
          </p>
        </div>

        <!-- Badge Manager -->
        <span 
          v-if="member.is_manager" 
          class="border-2 border-black px-2 py-1 text-xs font-bold flex-shrink-0"
        >
          MANAGER
        </span>
      </div>

      <!-- Message si aucun membre -->
      <div v-if="currentTeam.team_members.length === 0" class="text-center py-8 opacity-50">
        <p class="text-sm font-bold">Aucun membre</p>
      </div>
    </div>

    <!-- Message si aucune équipe -->
    <div v-else class="flex-1 flex items-center justify-center opacity-50">
      <p class="text-sm font-bold">
        {{ authStore.isManager 
          ? 'Vous n\'êtes manager d\'aucune équipe' 
          : 'Vous n\'êtes dans aucune équipe' 
        }}
      </p>
    </div>
  </div>
</template>