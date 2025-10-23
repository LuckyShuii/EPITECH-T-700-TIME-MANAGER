<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import type { Team, TeamMember } from '@/types/Team'
import type { Employee } from '@/types/Employee'
import API from '@/services/API'
import { useNotificationsStore } from '@/store/NotificationsStore'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

interface Props {
    modelValue: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
    'update:modelValue': [value: boolean]
}>()

const notificationsStore = useNotificationsStore()

// État : chargement
const isLoading = ref(false)

// État : liste des équipes
const teams = ref<Team[]>([])

// État : équipe sélectionnée
const selectedTeam = ref<Team | null>(null)

// État : liste de tous les employés (pour l'ajout de membre)
const allEmployees = ref<Employee[]>([])

// État : modales
const showCreateTeamModal = ref(false)
const showAddMemberModal = ref(false)
const showDeleteTeamConfirm = ref(false)
const showRemoveMemberConfirm = ref(false)

// État : membre à retirer
const memberToRemove = ref<TeamMember | null>(null)

// Formulaire de création d'équipe
const createTeamForm = ref({
    name: '',
    description: '',
    manager_uuid: '' // Nouveau champ pour le manager initial
})

// Formulaire d'ajout de membre
const addMemberForm = ref({
    user_uuid: '',
    is_manager: false
})

// Charger toutes les équipes
const loadTeams = async () => {
    isLoading.value = true
    try {
        const response = await API.teamAPI.getAll()
        teams.value = response.data

        // Si une équipe était sélectionnée, la recharger
        if (selectedTeam.value) {
            const updatedTeam = teams.value.find(t => t.uuid === selectedTeam.value?.uuid)
            if (updatedTeam) {
                selectedTeam.value = updatedTeam
            }
        }
    } catch (error) {
        console.error('Erreur lors du chargement des équipes:', error)
        notificationsStore.addNotification({
            status: 'error',
            title: 'Erreur de chargement',
            description: 'Impossible de charger les équipes'
        })
    } finally {
        isLoading.value = false
    }
}

// Charger tous les employés (pour le dropdown)
const loadAllEmployees = async () => {
    try {
        const response = await API.userAPI.getAll()
        allEmployees.value = response.data
    } catch (error) {
        console.error('Erreur lors du chargement des employés:', error)
        notificationsStore.addNotification({
            status: 'error',
            title: 'Erreur de chargement',
            description: 'Impossible de charger la liste des employés'
        })
    }
}

// Watcher : charger à l'ouverture de la modale
watch(() => props.modelValue, (newValue) => {
    if (newValue) {
        loadTeams()
        loadAllEmployees()
    }
})

// Sélectionner une équipe
const selectTeam = (team: Team) => {
    selectedTeam.value = team
}

// Fermer la modale principale
const closeModal = () => {
    emit('update:modelValue', false)
    selectedTeam.value = null
    teams.value = []
    resetCreateTeamForm()
    resetAddMemberForm()
}

// Reset formulaire création équipe
const resetCreateTeamForm = () => {
    createTeamForm.value = {
        name: '',
        description: '',
        manager_uuid: ''
    }
}

// Reset formulaire ajout membre
const resetAddMemberForm = () => {
    addMemberForm.value = {
        user_uuid: '',
        is_manager: false
    }
}

// Créer une nouvelle équipe
const createTeam = async () => {
    // Validation
    if (!createTeamForm.value.name.trim()) {
        notificationsStore.addNotification({
            status: 'warning',
            title: 'Champ requis',
            description: 'Le nom de l\'équipe est obligatoire'
        })
        return
    }

    if (!createTeamForm.value.manager_uuid) {
        notificationsStore.addNotification({
            status: 'warning',
            title: 'Manager requis',
            description: 'Veuillez sélectionner un manager pour l\'équipe'
        })
        return
    }

    try {
        // Préparer le payload avec le manager
        const payload = {
            name: createTeamForm.value.name,
            description: createTeamForm.value.description,
            member_uuids: [{
                user_uuid: createTeamForm.value.manager_uuid,
                is_manager: true
            }]
        }

        await API.teamAPI.createTeam(payload)
        notificationsStore.addNotification({
            status: 'success',
            title: 'Équipe créée',
            description: `L'équipe "${createTeamForm.value.name}" a été créée avec succès`
        })
        showCreateTeamModal.value = false
        resetCreateTeamForm()
        await loadTeams()
    } catch (error) {
        console.error('Erreur lors de la création de l\'équipe:', error)
        notificationsStore.addNotification({
            status: 'error',
            title: 'Erreur de création',
            description: 'Impossible de créer l\'équipe'
        })
    }
}

// Supprimer une équipe
const deleteTeam = async () => {
    if (!selectedTeam.value) return

    try {
        await API.teamAPI.deleteTeam(selectedTeam.value.uuid)
        notificationsStore.addNotification({
            status: 'success',
            title: 'Équipe supprimée',
            description: `L'équipe "${selectedTeam.value.name}" a été supprimée`
        })
        selectedTeam.value = null
        await loadTeams()
    } catch (error) {
        console.error('Erreur lors de la suppression:', error)
        notificationsStore.addNotification({
            status: 'error',
            title: 'Erreur de suppression',
            description: 'Impossible de supprimer l\'équipe'
        })
    }
}

// Ajouter un membre à l'équipe
const addMember = async () => {
    if (!selectedTeam.value || !addMemberForm.value.user_uuid) {
        notificationsStore.addNotification({
            status: 'warning',
            title: 'Champ requis',
            description: 'Veuillez sélectionner un employé'
        })
        return
    }

    try {
        await API.teamAPI.addMembers({
            team_uuid: selectedTeam.value.uuid,
            member_uuids: [{
                user_uuid: addMemberForm.value.user_uuid,
                is_manager: addMemberForm.value.is_manager
            }]
        })
        notificationsStore.addNotification({
            status: 'success',
            title: 'Membre ajouté',
            description: 'Le membre a été ajouté à l\'équipe'
        })
        showAddMemberModal.value = false
        resetAddMemberForm()
        await loadTeams()
    } catch (error) {
        console.error('Erreur lors de l\'ajout du membre:', error)
        notificationsStore.addNotification({
            status: 'error',
            title: 'Erreur d\'ajout',
            description: 'Impossible d\'ajouter le membre'
        })
    }
}

// Préparer la suppression d'un membre
const confirmRemoveMember = (member: TeamMember) => {
    memberToRemove.value = member
    showRemoveMemberConfirm.value = true
}

// Retirer un membre de l'équipe
const removeMember = async () => {
    if (!selectedTeam.value || !memberToRemove.value) return

    try {
        await API.teamAPI.removeMember(selectedTeam.value.uuid, memberToRemove.value.user_uuid)
        notificationsStore.addNotification({
            status: 'success',
            title: 'Membre retiré',
            description: `${memberToRemove.value.first_name} ${memberToRemove.value.last_name} a été retiré de l'équipe`
        })
        memberToRemove.value = null
        await loadTeams()
    } catch (error) {
        console.error('Erreur lors du retrait du membre:', error)
        notificationsStore.addNotification({
            status: 'error',
            title: 'Erreur de retrait',
            description: 'Impossible de retirer le membre'
        })
    }
}

// Computed : employés disponibles (non déjà dans l'équipe)
const availableEmployees = computed(() => {
    if (!selectedTeam.value || !selectedTeam.value.team_members) {
        return allEmployees.value
    }

    const memberUuids = selectedTeam.value.team_members.map(m => m.user_uuid)
    return allEmployees.value.filter(emp => emp.uuid && !memberUuids.includes(emp.uuid))
})
</script>

<template>
    <Transition name="modal">
        <div v-if="modelValue" class="modal modal-open" @click.self="closeModal">
            <div class="modal-box max-w-6xl h-[80vh]">
                <!-- Header -->
                <div class="flex justify-between items-center mb-6">
                    <h3 class="font-bold text-2xl">Gestion des équipes</h3>
                    <div class="flex gap-2">
                        <button class="btn btn-primary btn-sm" @click="showCreateTeamModal = true">
                            + Créer une équipe
                        </button>
                        <button class="btn btn-sm btn-circle btn-ghost" @click="closeModal">
                            ✕
                        </button>
                    </div>
                </div>

                <!-- Indicateur de chargement global -->
                <div v-if="isLoading && teams.length === 0" class="flex justify-center items-center h-full">
                    <span class="loading loading-spinner loading-lg"></span>
                </div>

                <!-- Contenu principal -->
                <div v-else class="grid grid-cols-[300px_1fr] gap-6" style="height: calc(80vh - 120px);">
                    <!-- Gauche : Liste des équipes -->
                    <div class="border-r border-base-300 pr-4 overflow-y-auto">
                        <h4 class="font-semibold text-sm mb-3 opacity-70">
                            Équipes ({{ teams.length }})
                        </h4>

                        <!-- Liste des équipes -->
                        <div class="space-y-1">
                            <div v-for="team in teams" :key="team.uuid" @click="selectTeam(team)"
                                class="p-3 rounded-lg cursor-pointer transition-colors" :class="{
                                    'bg-primary text-primary-content': selectedTeam?.uuid === team.uuid,
                                    'bg-base-200 hover:bg-base-300': selectedTeam?.uuid !== team.uuid
                                }">
                                <p class="font-medium text-sm">{{ team.name }}</p>
                                <p class="text-xs opacity-70 mt-1">{{ team.team_members.length }} membre(s)</p>
                            </div>
                        </div>

                        <!-- Message si aucune équipe -->
                        <div v-if="teams.length === 0" class="text-center py-8 opacity-50">
                            <p class="text-sm">Aucune équipe créée</p>
                        </div>
                    </div>

                    <!-- Droite : Détails de l'équipe sélectionnée -->
                    <div v-if="selectedTeam" class="overflow-y-auto">
                        <!-- En-tête équipe -->
                        <div class="mb-6">
                            <h4 class="font-semibold text-xl">{{ selectedTeam.name }}</h4>
                            <p class="text-sm opacity-70 mt-1">{{ selectedTeam.description }}</p>
                        </div>

                        <!-- Liste des membres -->
                        <div class="mb-6">
                            <div class="flex justify-between items-center mb-3">
                                <h5 class="font-semibold text-sm opacity-70">
                                    Membres ({{ selectedTeam.team_members.length }})
                                </h5>
                                <button class="btn btn-sm btn-primary" @click="showAddMemberModal = true">
                                    + Ajouter un membre
                                </button>
                            </div>

                            <!-- Tableau des membres -->
                            <div class="space-y-2">
                                <div v-for="member in selectedTeam.team_members" :key="member.user_uuid"
                                    class="flex items-center justify-between p-3 bg-base-200 rounded-lg">
                                    <div class="flex-1">
                                        <p class="font-medium text-sm">
                                            {{ member.first_name }} {{ member.last_name }}
                                        </p>
                                        <p class="text-xs opacity-70">{{ member.email }}</p>
                                    </div>

                                    <div class="flex items-center gap-2">
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

                                        <!-- Bouton retirer -->
                                        <button class="btn btn-ghost btn-xs btn-circle"
                                            @click="confirmRemoveMember(member)">
                                            ✕
                                        </button>
                                    </div>
                                </div>
                            </div>

                            <!-- Message si aucun membre -->
                            <div v-if="selectedTeam.team_members.length === 0" class="text-center py-8 opacity-50">
                                <p class="text-sm">Aucun membre dans cette équipe</p>
                            </div>
                        </div>

                        <div class="divider"></div>

                        <!-- Actions équipe -->
                        <div class="flex justify-end">
                            <button class="btn btn-error btn-outline" @click="showDeleteTeamConfirm = true">
                                Supprimer l'équipe
                            </button>
                        </div>
                    </div>

                    <!-- Placeholder si aucune équipe sélectionnée -->
                    <div v-else class="flex items-center justify-center opacity-50 h-full">
                        <p>Sélectionnez une équipe dans la liste</p>
                    </div>
                </div>
            </div>
        </div>
    </Transition>

    <!-- Modal de création d'équipe -->
    <dialog :open="showCreateTeamModal" class="modal">
        <div class="modal-box">
            <h3 class="font-bold text-lg mb-4">Créer une nouvelle équipe</h3>

            <div class="space-y-4">
                <!-- Nom de l'équipe -->
                <div>
                    <label class="label">
                        <span class="label-text">Nom de l'équipe *</span>
                    </label>
                    <input v-model="createTeamForm.name" type="text" placeholder="Ex: Équipe développement"
                        class="input input-bordered w-full" />
                </div>

                <!-- Description -->
                <div>
                    <label class="label">
                        <span class="label-text">Description</span>
                    </label>
                    <textarea v-model="createTeamForm.description" placeholder="Description de l'équipe..."
                        class="textarea textarea-bordered w-full" rows="3"></textarea>
                </div>

                <!-- Sélection du manager -->
                <div>
                    <label class="label">
                        <span class="label-text">Manager de l'équipe *</span>
                    </label>
                    <select v-model="createTeamForm.manager_uuid" class="select select-bordered w-full">
                        <option value="">Sélectionnez un manager</option>
                        <option v-for="employee in allEmployees" :key="employee.uuid" :value="employee.uuid">
                            {{ employee.first_name }} {{ employee.last_name }} ({{ employee.email }})
                        </option>
                    </select>
                    <label class="label">
                        <span class="label-text-alt opacity-70">Cette personne sera définie comme manager de
                            l'équipe</span>
                    </label>
                </div>
            </div>

            <div class="modal-action">
                <button class="btn btn-ghost" @click="showCreateTeamModal = false; resetCreateTeamForm()">
                    Annuler
                </button>
                <button class="btn btn-primary" @click="createTeam">
                    Créer l'équipe
                </button>
            </div>
        </div>
        <div class="modal-backdrop" @click="showCreateTeamModal = false"></div>
    </dialog>

    <!-- Modal d'ajout de membre -->
    <dialog :open="showAddMemberModal" class="modal">
        <div class="modal-box">
            <h3 class="font-bold text-lg mb-4">Ajouter un membre</h3>

            <div class="space-y-4">
                <!-- Sélection employé -->
                <div>
                    <label class="label">
                        <span class="label-text">Employé *</span>
                    </label>
                    <select v-model="addMemberForm.user_uuid" class="select select-bordered w-full">
                        <option value="">Sélectionnez un employé</option>
                        <option v-for="employee in availableEmployees" :key="employee.uuid" :value="employee.uuid">
                            {{ employee.first_name }} {{ employee.last_name }} ({{ employee.email }})
                        </option>
                    </select>
                </div>

                <!-- Checkbox Manager -->
                <div>
                    <label class="label cursor-pointer justify-start gap-2">
                        <input v-model="addMemberForm.is_manager" type="checkbox" class="checkbox checkbox-primary" />
                        <span class="label-text">Définir comme manager de l'équipe</span>
                    </label>
                </div>
            </div>

            <div class="modal-action">
                <button class="btn btn-ghost" @click="showAddMemberModal = false; resetAddMemberForm()">
                    Annuler
                </button>
                <button class="btn btn-primary" @click="addMember">
                    Ajouter
                </button>
            </div>
        </div>
        <div class="modal-backdrop" @click="showAddMemberModal = false"></div>
    </dialog>

    <!-- ConfirmDialog pour supprimer l'équipe -->
    <ConfirmDialog v-model="showDeleteTeamConfirm" title="Supprimer l'équipe"
        :message="`Êtes-vous sûr de vouloir supprimer l'équipe <span class='font-bold'>${selectedTeam?.name}</span> ?<br><span class='text-sm opacity-70'>Tous les membres seront retirés de cette équipe.</span>`"
        confirm-text="Supprimer l'équipe" cancel-text="Annuler" variant="error" @confirm="deleteTeam" />

    <!-- ConfirmDialog pour retirer un membre -->
    <ConfirmDialog v-model="showRemoveMemberConfirm" title="Retirer le membre"
        :message="`Êtes-vous sûr de vouloir retirer <span class='font-bold'>${memberToRemove?.first_name} ${memberToRemove?.last_name}</span> de l'équipe ?`"
        confirm-text="Retirer" cancel-text="Annuler" variant="warning" @confirm="removeMember" />
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