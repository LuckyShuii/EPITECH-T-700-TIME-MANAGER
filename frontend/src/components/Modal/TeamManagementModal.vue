<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useAuthStore } from '@/store/AuthStore'
import type { Team, TeamMember, TeamInfo } from '@/types/Team'
import API from '@/services/API'
import { useNotificationsStore } from '@/store/NotificationsStore'

interface Props {
    modelValue: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
    'update:modelValue': [value: boolean]
}>()

const authStore = useAuthStore()
const notificationsStore = useNotificationsStore()

// État : chargement
const isLoading = ref(false)

// État : équipes du user
const userTeams = ref<TeamInfo[]>([])

// État : équipe sélectionnée (Team complète avec membres)
const selectedTeam = ref<Team | null>(null)

// État : membre sélectionné pour affichage détails
const selectedMember = ref<TeamMember | null>(null)

// Computed : titre dynamique
const modalTitle = computed(() => {
    return userTeams.value.length > 1 ? 'Mes équipes' : 'Mon équipe'
})

// Computed : afficher le dropdown ou non
const showTeamDropdown = computed(() => {
    return userTeams.value.length > 1
})

// Charger les équipes du user connecté
const loadUserTeams = async () => {
    if (!authStore.user?.user_uuid) {
        console.error('Aucun UUID utilisateur disponible')
        return
    }

    isLoading.value = true
    try {
        const response = await API.userAPI.getUserSpecific(authStore.user.user_uuid)
        console.log('👤 User data:', response.data)
        userTeams.value = response.data.teams || []
        console.log('📋 Teams chargées:', userTeams.value)

        // Sélectionner automatiquement la première équipe
        if (userTeams.value.length > 0) {
            await loadTeamMembers(userTeams.value[0].team_uuid)

        }
    } catch (error) {
        console.error('Erreur lors du chargement des équipes:', error)
        notificationsStore.addNotification({
            status: 'error',
            title: 'Erreur de chargement',
            description: 'Impossible de charger vos équipes'
        })
    }
}

// Charger les membres d'une équipe spécifique
const loadTeamMembers = async (teamUuid: string) => {
    isLoading.value = true
    try {
        const response = await API.teamAPI.getTeam(teamUuid)
        console.log('👥 Team data:', response.data)
        selectedTeam.value = response.data

        // Réinitialiser la sélection du membre
        selectedMember.value = null
    } catch (error) {
        console.error('Erreur lors du chargement des membres:', error)
        notificationsStore.addNotification({
            status: 'error',
            title: 'Erreur de chargement',
            description: 'Impossible de charger les membres de l\'équipe'
        })
    }
}

// Watcher : charger les équipes à l'ouverture de la modale
watch(() => props.modelValue, (newValue) => {
    if (newValue) {
        loadUserTeams()
    }
})

// Sélectionner un membre pour afficher ses détails
const selectMember = (member: TeamMember) => {
    selectedMember.value = member
}

// Fermer la modale
const closeModal = () => {
    emit('update:modelValue', false)
    selectedTeam.value = null
    selectedMember.value = null
    userTeams.value = []
}

</script>

<template>
    <Transition name="modal">
        <div v-if="modelValue" class="modal modal-open" @click.self="closeModal">
            <div class="modal-box max-w-6xl h-[80vh]">
                <!-- Header -->
                <div class="flex justify-between items-center mb-6">
                    <h3 class="font-bold text-2xl">{{ modalTitle }}</h3>
                    <button class="btn btn-sm btn-circle btn-ghost" @click="closeModal">
                        ✕
                    </button>
                </div>

                <!-- Indicateur de chargement global -->
                <div v-if="isLoading && !selectedTeam" class="flex justify-center items-center h-full">
                    <span class="loading loading-spinner loading-lg"></span>
                </div>

                <!-- Contenu principal -->
                <div v-else class="grid grid-cols-[300px_1fr] gap-6" style="height: calc(80vh - 120px);">
                    <!-- Gauche : Sélection d'équipe + Liste des membres -->
                    <div class="border-r border-base-300 pr-4 overflow-y-auto">
                        <!-- Dropdown équipes (si plusieurs) -->
                        <div v-if="showTeamDropdown" class="mb-4">
                            <select class="select select-bordered w-full"
                                @change="loadTeamMembers(($event.target as HTMLSelectElement).value)"
                                :value="selectedTeam?.uuid">
                                <option v-for="team in userTeams" :key="team.team_uuid" :value="team.team_uuid">
                                    {{ team.team_name }}
                                </option>
                            </select>
                        </div>

                        <!-- Titre équipe (si une seule) -->
                        <div v-else-if="selectedTeam" class="mb-4">
                            <h4 class="font-semibold text-lg">{{ selectedTeam.name }}</h4>
                            <p class="text-sm opacity-70">{{ selectedTeam.description }}</p>
                        </div>

                        <!-- Sous-titre membres -->
                        <h4 class="font-semibold text-sm mb-3 opacity-70">
                            Membres ({{ selectedTeam?.team_members.length || 0 }})
                        </h4>

                        <!-- Liste des membres -->
                        <div v-if="selectedTeam" class="space-y-1">
                            <div v-for="member in selectedTeam.team_members" :key="member.user_uuid"
                                @click="selectMember(member)" class="p-2 rounded-lg cursor-pointer transition-colors"
                                :class="{
                                    'bg-primary text-primary-content': selectedMember?.user_uuid === member.user_uuid,
                                    'bg-base-200 hover:bg-base-300': selectedMember?.user_uuid !== member.user_uuid,
                                    'opacity-50': member.status === 'inactive'
                                }">
                                <div class="flex items-center justify-between">
                                    <div>
                                        <p class="font-medium text-sm">
                                            {{ member.first_name }} {{ member.last_name }}
                                        </p>
                                        <div class="flex gap-1 mt-1">
                                            <!-- Badge Manager -->
                                            <span v-if="member.is_manager" class="badge badge-xs badge-primary">
                                                Manager
                                            </span>

                                            <!-- Badge Status -->
                                            <span class="badge badge-xs" :class="{
                                                'badge-success': member.status === 'active',
                                                'badge-error': member.status === 'inactive',
                                                'badge-warning': member.status === 'pending'
                                            }">
                                                {{ member.status === 'active' ? 'Actif' : member.status === 'inactive' ?
                                                    'Inactif' : 'En attente' }}
                                            </span>
                                        </div>
                                    </div>

                                    <!-- Icône hourglass pour pending -->
                                    <svg v-if="member.status === 'pending'" xmlns="http://www.w3.org/2000/svg"
                                        fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
                                        class="w-4 h-4 opacity-50">
                                        <path stroke-linecap="round" stroke-linejoin="round"
                                            d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                                    </svg>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Droite : Détails du membre sélectionné -->
                    <div v-if="selectedMember" class="overflow-y-auto">
                        <h4 class="font-semibold text-xl mb-6">
                            {{ selectedMember.first_name }} {{ selectedMember.last_name }}
                        </h4>

                        <div class="space-y-4">
                            <!-- Username -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                                <span class="opacity-70">Username</span>
                                <span class="font-medium">{{ selectedMember.username }}</span>
                            </div>

                            <!-- Email -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                                <span class="opacity-70">Email</span>
                                <span class="font-medium">{{ selectedMember.email }}</span>
                            </div>

                            <!-- Téléphone -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                                <span class="opacity-70">Téléphone</span>
                                <span class="font-medium">{{ selectedMember.phone_number }}</span>
                            </div>

                            <div class="divider"></div>

                            <!-- Rôles -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
                                <span class="opacity-70 pt-1">Rôles</span>
                                <div class="flex flex-wrap gap-2">
                                    <span v-for="role in selectedMember.roles" :key="role" class="badge badge-outline">
                                        {{ role }}
                                    </span>
                                </div>
                            </div>

                            <!-- Statut -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                                <span class="opacity-70">Statut</span>
                                <span class="badge" :class="{
                                    'badge-success': selectedMember.status === 'active',
                                    'badge-error': selectedMember.status === 'inactive',
                                    'badge-warning': selectedMember.status === 'pending'
                                }">
                                    {{ selectedMember.status === 'active' ? 'Actif' : selectedMember.status ===
                                        'inactive' ? 'Inactif' : 'En attente' }}
                                </span>
                            </div>

                            <!-- Manager de cette équipe -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                                <span class="opacity-70">Manager</span>
                                <span class="font-medium">
                                    {{ selectedMember.is_manager ? 'Oui' : 'Non' }}
                                </span>
                            </div>

                            <!-- Volume horaire (visible uniquement pour les managers/admins) -->
                            <div v-if="authStore.user?.roles.includes('manager') || authStore.user?.roles.includes('admin')"
                                class="grid grid-cols-[150px_1fr] gap-2 items-center">
                                <span class="opacity-70">Volume horaire</span>
                                <div class="flex items-center gap-2">
                                    <span class="font-medium">{{ selectedMember.weekly_rate }}h / semaine</span>
                                    <span class="badge badge-sm badge-outline">{{ selectedMember.weekly_rate_name }}</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Placeholder si aucun membre sélectionné -->
                    <div v-else class="flex items-center justify-center opacity-50 h-full">
                        <p>Sélectionnez un membre dans la liste</p>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
</template>
<style scoped>
.modal-enter-active,
.modal-leave-active {
    transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
    opacity: 0;
}
</style>