<script setup lang="ts">
import { ref, watch, computed, defineModel } from 'vue'
import { useAuthStore } from '@/store/AuthStore'
import type { Team, TeamMember, TeamInfo } from '@/types/Team'
import API from '@/services/API'
import { useNotificationsStore } from '@/store/NotificationsStore'

const modelValue = defineModel<boolean>()
const authStore = useAuthStore()
const notificationsStore = useNotificationsStore()

const isLoading = ref(false)
const userTeams = ref<TeamInfo[]>([])
const selectedTeam = ref<Team | null>(null)
const selectedMember = ref<TeamMember | null>(null)

const modalTitle = computed(() => {
    return userTeams.value.length > 1 ? 'Mes équipes' : 'Mon équipe'
})

const showTeamDropdown = computed(() => {
    return userTeams.value.length > 1
})

const loadUserTeams = async () => {
    if (!authStore.user?.user_uuid) {
        console.error('Aucun UUID utilisateur disponible')
        return
    }

    isLoading.value = true
    try {
        const response = await API.userAPI.getUserSpecific(authStore.user.user_uuid)
        userTeams.value = response.data.teams || []

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

const loadTeamMembers = async (teamUuid: string) => {
    isLoading.value = true
    try {
        const response = await API.teamAPI.getTeam(teamUuid)
        selectedTeam.value = response.data
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

watch(() => modelValue.value, (newValue) => {
    if (newValue) {
        loadUserTeams()
    }
})

const selectMember = (member: TeamMember) => {
    selectedMember.value = member
}

const closeModal = () => {
    modelValue.value = false
    selectedTeam.value = null
    selectedMember.value = null
    userTeams.value = []
}
</script>

<template>
    <Transition name="modal">
        <div v-if="modelValue" class="modal modal-open" @click.self="closeModal">
            <div class="modal-box border-4 border-black rounded-none bg-base-100 max-w-6xl h-[80vh]">
                <!-- Header -->
                <div class="flex justify-between items-center mb-6 pb-4 border-b-2 border-black">
                    <h3 class="font-black text-2xl uppercase tracking-wider">{{ modalTitle }}</h3>
                    <button class="border-2 border-black w-8 h-8 flex items-center justify-center font-bold hover:bg-black hover:text-white transition-none" @click="closeModal">
                        X
                    </button>
                </div>

                <!-- Indicateur de chargement global -->
                <div v-if="isLoading && !selectedTeam" class="flex justify-center items-center h-full">
                    <span class="loading loading-spinner loading-lg"></span>
                </div>

                <!-- Contenu principal -->
                <div v-else class="grid grid-cols-[300px_1fr] gap-6" style="height: calc(80vh - 120px);">
                    <!-- Gauche : Sélection d'équipe + Liste des membres -->
                    <div class="border-r-2 border-black pr-4 overflow-y-auto">
                        <!-- Dropdown équipes (si plusieurs) -->
                        <div v-if="showTeamDropdown" class="mb-4">
                            <select class="select select-bordered w-full border-2"
                                @change="loadTeamMembers(($event.target as HTMLSelectElement).value)"
                                :value="selectedTeam?.uuid">
                                <option v-for="team in userTeams" :key="team.team_uuid" :value="team.team_uuid">
                                    {{ team.team_name }}
                                </option>
                            </select>
                        </div>

                        <!-- Titre équipe (si une seule) -->
                        <div v-else-if="selectedTeam" class="mb-4">
                            <h4 class="font-bold uppercase text-sm tracking-wider">{{ selectedTeam.name }}</h4>
                            <p class="text-xs font-bold mt-1">{{ selectedTeam.description }}</p>
                        </div>

                        <!-- Sous-titre membres -->
                        <h4 class="font-bold uppercase text-xs mb-3 tracking-widest">
                            Membres ({{ selectedTeam?.team_members.length || 0 }})
                        </h4>

                        <!-- Liste des membres -->
                        <div v-if="selectedTeam" class="space-y-1">
                            <div v-for="member in selectedTeam.team_members" :key="member.user_uuid"
                                @click="selectMember(member)" class="p-2 border-2 border-gray-400 cursor-pointer" 
                                :class="{
                                    'bg-black text-white border-black': selectedMember?.user_uuid === member.user_uuid,
                                    'hover:bg-gray-100': selectedMember?.user_uuid !== member.user_uuid,
                                    'opacity-50': member.status === 'inactive'
                                }">
                                <div class="flex items-center justify-between">
                                    <div>
                                        <p class="font-bold text-sm">
                                            {{ member.first_name }} {{ member.last_name }}
                                        </p>
                                        <div class="flex gap-1 mt-1">
                                            <!-- Badge Manager -->
                                            <span v-if="member.is_manager" class="text-xs font-bold border-2 border-black px-2 py-1">
                                                MANAGER
                                            </span>

                                            <!-- Badge Status -->
                                            <span class="text-xs font-bold border-2 px-2 py-1" :class="{
                                                'border-green-700 text-green-700': member.status === 'active',
                                                'border-red-700 text-red-700': member.status === 'inactive',
                                                'border-yellow-700 text-yellow-700': member.status === 'pending'
                                            }">
                                                {{ member.status === 'active' ? 'ACTIF' : member.status === 'inactive' ?
                                                    'INACTIF' : 'ATTENTE' }}
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
                        <h4 class="font-bold text-lg uppercase tracking-wider mb-6">
                            {{ selectedMember.first_name }} {{ selectedMember.last_name }}
                        </h4>

                        <div class="space-y-4">
                            <!-- Username -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-center border-b-2 border-black pb-4">
                                <span class="font-bold uppercase text-xs tracking-widest">Username</span>
                                <span class="font-bold">{{ selectedMember.username }}</span>
                            </div>

                            <!-- Email -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-center border-b-2 border-black pb-4">
                                <span class="font-bold uppercase text-xs tracking-widest">Email</span>
                                <span class="font-bold">{{ selectedMember.email }}</span>
                            </div>

                            <!-- Téléphone -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-center border-b-2 border-black pb-4">
                                <span class="font-bold uppercase text-xs tracking-widest">Téléphone</span>
                                <span class="font-bold">{{ selectedMember.phone_number }}</span>
                            </div>

                            <!-- Rôles -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-start pt-4">
                                <span class="font-bold uppercase text-xs tracking-widest">Rôles</span>
                                <div class="flex flex-wrap gap-2">
                                    <span v-for="role in selectedMember.roles" :key="role" class="border-2 border-black px-2 py-1 font-bold uppercase text-xs">
                                        {{ role }}
                                    </span>
                                </div>
                            </div>

                            <!-- Statut -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-center pt-4">
                                <span class="font-bold uppercase text-xs tracking-widest">Statut</span>
                                <span class="font-bold border-2 px-2 py-1" :class="{
                                    'border-green-700 text-green-700': selectedMember.status === 'active',
                                    'border-red-700 text-red-700': selectedMember.status === 'inactive',
                                    'border-yellow-700 text-yellow-700': selectedMember.status === 'pending'
                                }">
                                    {{ selectedMember.status === 'active' ? 'ACTIF' : selectedMember.status ===
                                        'inactive' ? 'INACTIF' : 'ATTENTE' }}
                                </span>
                            </div>

                            <!-- Manager de cette équipe -->
                            <div class="grid grid-cols-[150px_1fr] gap-2 items-center pt-4">
                                <span class="font-bold uppercase text-xs tracking-widest">Manager</span>
                                <span class="font-bold">
                                    {{ selectedMember.is_manager ? 'OUI' : 'NON' }}
                                </span>
                            </div>

                            <!-- Volume horaire -->
                            <div v-if="authStore.user?.roles.includes('manager') || authStore.user?.roles.includes('admin')"
                                class="grid grid-cols-[150px_1fr] gap-2 items-center pt-4">
                                <span class="font-bold uppercase text-xs tracking-widest">Volume</span>
                                <div class="flex items-center gap-2">
                                    <span class="font-bold">{{ selectedMember.weekly_rate }}h/sem</span>
                                    <span class="border-2 border-black px-2 py-1 text-xs font-bold">{{ selectedMember.weekly_rate_name }}</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Placeholder si aucun membre sélectionné -->
                    <div v-else class="flex items-center justify-center opacity-50 h-full">
                        <p class="font-bold">Sélectionnez un membre</p>
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