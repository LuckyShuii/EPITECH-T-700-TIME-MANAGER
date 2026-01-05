<script setup lang="ts">
import { ref, watch, computed, defineModel } from 'vue'
import type { Team, TeamMember } from '@/types/Team'
import type { Employee } from '@/types/Employee'
import API from '@/services/API'
import { useNotificationsStore } from '@/store/NotificationsStore'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import BaseModal from '@/components/Modal/BaseModal.vue'
import { PencilIcon, CheckIcon, XMarkIcon } from '@heroicons/vue/24/outline'

const modelValue = defineModel<boolean>()
const notificationsStore = useNotificationsStore()

const isLoading = ref(false)
const teams = ref<Team[]>([])
const selectedTeam = ref<Team | null>(null)
const allEmployees = ref<Employee[]>([])

const showCreateTeamModal = ref(false)
const showAddMemberModal = ref(false)
const showDeleteTeamConfirm = ref(false)
const showRemoveMemberConfirm = ref(false)

const memberToRemove = ref<TeamMember | null>(null)
const isEditingTeam = ref(false)
const editTeamForm = ref({ name: '', description: '' })
const createTeamForm = ref({ name: '', description: '', manager_uuid: '' })
const addMemberForm = ref({ user_uuid: '', is_manager: false })

const loadTeams = async () => {
    isLoading.value = true
    try {
        const response = await API.teamAPI.getAll()
        teams.value = response.data
        if (selectedTeam.value) {
            const updatedTeam = teams.value.find(t => t.uuid === selectedTeam.value?.uuid)
            if (updatedTeam) selectedTeam.value = updatedTeam
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

watch(() => modelValue.value, (newValue) => {
    if (newValue) {
        loadTeams()
        loadAllEmployees()
    }
})

const selectTeam = (team: Team) => {
    selectedTeam.value = team
}

const closeModal = () => {
    modelValue.value = false
    selectedTeam.value = null
    teams.value = []
    resetCreateTeamForm()
    resetAddMemberForm()
}

const resetCreateTeamForm = () => {
    createTeamForm.value = { name: '', description: '', manager_uuid: '' }
}

const resetAddMemberForm = () => {
    addMemberForm.value = { user_uuid: '', is_manager: false }
}

const createTeam = async () => {
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

const confirmRemoveMember = (member: TeamMember) => {
    memberToRemove.value = member
    showRemoveMemberConfirm.value = true
}

const removeMember = async () => {
    if (!selectedTeam.value || !memberToRemove.value) return

    try {
        await API.teamAPI.removeMember(selectedTeam.value.uuid, memberToRemove.value.user_uuid)
        notificationsStore.addNotification({
            status: 'success',
            title: 'Membre retiré',
            description: `${memberToRemove.value.first_name} ${memberToRemove.value.last_name} a été retiré`
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

const availableEmployees = computed(() => {
    if (!selectedTeam.value || !selectedTeam.value.team_members) {
        return allEmployees.value
    }
    const memberUuids = selectedTeam.value.team_members.map(m => m.user_uuid)
    return allEmployees.value.filter(emp => emp.uuid && !memberUuids.includes(emp.uuid))
})

const startEditingTeam = () => {
    if (!selectedTeam.value) return
    editTeamForm.value = {
        name: selectedTeam.value.name,
        description: selectedTeam.value.description
    }
    isEditingTeam.value = true
}

const cancelEditingTeam = () => {
    isEditingTeam.value = false
    editTeamForm.value = { name: '', description: '' }
}

const saveTeamEdit = async () => {
    if (!selectedTeam.value) return

    if (!editTeamForm.value.name.trim()) {
        notificationsStore.addNotification({
            status: 'warning',
            title: 'Champ requis',
            description: 'Le nom de l\'équipe est obligatoire'
        })
        return
    }

    try {
        await API.teamAPI.updateTeam(selectedTeam.value.uuid, {
            name: editTeamForm.value.name,
            description: editTeamForm.value.description
        })

        notificationsStore.addNotification({
            status: 'success',
            title: 'Équipe modifiée',
            description: 'Les informations ont été mises à jour'
        })

        isEditingTeam.value = false
        await loadTeams()
    } catch (error) {
        console.error('Erreur lors de la modification:', error)
        notificationsStore.addNotification({
            status: 'error',
            title: 'Erreur de modification',
            description: 'Impossible de modifier l\'équipe'
        })
    }
}
</script>

<template>
    <BaseModal v-model="modelValue" title="Gestion des équipes">
        <div v-if="isLoading && teams.length === 0" class="flex justify-center items-center h-full">
            <span class="loading loading-spinner loading-lg"></span>
        </div>

        <div v-else class="grid grid-cols-[300px_1fr] gap-6" style="height: calc(80vh - 240px);">
            <!-- Gauche : Liste des équipes -->
            <div class="border-r-2 border-black pr-4 overflow-y-auto">
                <h4 class="font-bold uppercase text-sm mb-3 tracking-wider">
                    Équipes ({{ teams.length }})
                </h4>

                <div class="space-y-1">
                    <div v-for="team in teams" :key="team.uuid" @click="selectTeam(team)"
                        class="p-3 border-2 border-gray-400 cursor-pointer" :class="{
                            'bg-black text-white border-black': selectedTeam?.uuid === team.uuid,
                            'hover:bg-gray-500': selectedTeam?.uuid !== team.uuid
                        }">
                        <p class="font-bold text-sm">{{ team.name }}</p>
                        <p class="text-xs mt-1">{{ team.team_members.length }} membres</p>
                    </div>
                </div>

                <div v-if="teams.length === 0" class="text-center py-8 opacity-50">
                    <p class="font-bold text-xs">Aucune équipe créée</p>
                </div>
            </div>

            <!-- Droite : Détails de l'équipe sélectionnée -->
            <div v-if="selectedTeam" class="overflow-y-auto">
                <!-- En-tête équipe -->
                <div class="mb-6">
                    <div v-if="!isEditingTeam" class="flex items-start gap-2">
                        <div class="flex-1">
                            <div class="flex items-center gap-2">
                                <h4 class="font-bold text-lg uppercase">{{ selectedTeam.name }}</h4>
                                <button @click="startEditingTeam" class="border-2 border-black w-6 h-6 flex items-center justify-center" title="Modifier">
                                    <PencilIcon class="w-3 h-3" />
                                </button>
                            </div>
                            <p class="text-xs mt-1 font-bold">{{ selectedTeam.description || 'Aucune description' }}</p>
                        </div>
                    </div>

                    <div v-else class="space-y-3">
                        <div class="flex items-center gap-2">
                            <input v-model="editTeamForm.name" type="text"
                                class="input input-bordered input-sm flex-1 border-2" placeholder="Nom de l'équipe" />
                            <button @click="saveTeamEdit" class="border-2 border-black w-6 h-6 flex items-center justify-center bg-green-600 text-white" title="Sauvegarder">
                                <CheckIcon class="w-3 h-3" />
                            </button>
                            <button @click="cancelEditingTeam" class="border-2 border-black w-6 h-6 flex items-center justify-center bg-red-600 text-white" title="Annuler">
                                <XMarkIcon class="w-3 h-3" />
                            </button>
                        </div>
                        <textarea v-model="editTeamForm.description"
                            class="textarea textarea-bordered textarea-sm w-full border-2"
                            placeholder="Description de l'équipe" rows="2"></textarea>
                    </div>
                </div>

                <!-- Liste des membres -->
                <div class="mb-6">
                    <div class="flex justify-between items-center mb-3">
                        <h5 class="font-bold uppercase text-xs">
                            Membres ({{ selectedTeam.team_members.length }})
                        </h5>
                        <button class="brutal-btn brutal-btn-primary text-xs" @click="showAddMemberModal = true">
                            Ajouter
                        </button>
                    </div>

                    <div class="space-y-2">
                        <div v-for="member in selectedTeam.team_members" :key="member.user_uuid"
                            class="flex items-center justify-between p-3 border-2 border-black">
                            <div class="flex-1">
                                <p class="font-bold text-sm">
                                    {{ member.first_name }} {{ member.last_name }}
                                </p>
                                <p class="text-xs">{{ member.email }}</p>
                            </div>

                            <div class="flex items-center gap-2">
                                <span v-if="member.is_manager" class="font-bold text-xs border-2 border-black px-2 py-1">
                                    MANAGER
                                </span>

                                <span class="font-bold text-xs border-2 px-2 py-1" :class="{
                                    'border-green-700 text-green-700': member.status === 'active',
                                    'border-red-700 text-red-700': member.status === 'inactive',
                                    'border-yellow-700 text-yellow-700': member.status === 'pending'
                                }">
                                    {{ member.status === 'active' ? 'ACTIF' : member.status === 'inactive' ? 'INACTIF' : 'ATTENTE' }}
                                </span>

                                <button class="border-2 border-black w-5 h-5 flex items-center justify-center text-xs font-bold" @click="confirmRemoveMember(member)">
                                    X
                                </button>
                            </div>
                        </div>
                    </div>

                    <div v-if="selectedTeam.team_members.length === 0" class="text-center py-8 opacity-50">
                        <p class="font-bold text-xs">Aucun membre dans cette équipe</p>
                    </div>
                </div>

                <div class="h-px bg-black my-4"></div>

                <!-- Actions équipe -->
                <div class="flex justify-end">
                    <button class="brutal-btn brutal-btn-error" @click="showDeleteTeamConfirm = true">
                        Supprimer l'équipe
                    </button>
                </div>
            </div>

            <div v-else class="flex items-center justify-center opacity-50 h-full">
                <p class="font-bold">Sélectionnez une équipe</p>
            </div>
        </div>

        <!-- Footer -->
        <template #footer>
            <button class="brutal-btn brutal-btn-primary" @click="showCreateTeamModal = true">
                Créer une équipe
            </button>
        </template>
    </BaseModal>

    <!-- Modal de création d'équipe -->
    <BaseModal v-model="showCreateTeamModal" title="Créer une nouvelle équipe">
        <div class="space-y-4">
            <div>
                <label class="label">
                    <span class="label-text font-bold uppercase text-xs">Nom de l'équipe</span>
                </label>
                <input v-model="createTeamForm.name" type="text" placeholder="Ex: Équipe développement"
                    class="input input-bordered w-full border-2" />
            </div>

            <div>
                <label class="label">
                    <span class="label-text font-bold uppercase text-xs">Description</span>
                </label>
                <textarea v-model="createTeamForm.description" placeholder="Description de l'équipe..."
                    class="textarea textarea-bordered w-full border-2" rows="3"></textarea>
            </div>

            <div>
                <label class="label">
                    <span class="label-text font-bold uppercase text-xs">Manager de l'équipe</span>
                </label>
                <select v-model="createTeamForm.manager_uuid" class="select select-bordered w-full border-2">
                    <option value="">Sélectionnez un manager</option>
                    <option v-for="employee in allEmployees" :key="employee.uuid" :value="employee.uuid">
                        {{ employee.first_name }} {{ employee.last_name }} ({{ employee.email }})
                    </option>
                </select>
            </div>
        </div>

        <template #footer>
            <button class="brutal-btn" @click="showCreateTeamModal = false; resetCreateTeamForm()">
                Annuler
            </button>
            <button class="brutal-btn brutal-btn-success" @click="createTeam">
                Créer l'équipe
            </button>
        </template>
    </BaseModal>

    <!-- Modal d'ajout de membre -->
    <BaseModal v-model="showAddMemberModal" title="Ajouter un membre">
        <div class="space-y-4">
            <div>
                <label class="label">
                    <span class="label-text font-bold uppercase text-xs">Employé</span>
                </label>
                <select v-model="addMemberForm.user_uuid" class="select select-bordered w-full border-2">
                    <option value="">Sélectionnez un employé</option>
                    <option v-for="employee in availableEmployees" :key="employee.uuid" :value="employee.uuid">
                        {{ employee.first_name }} {{ employee.last_name }} ({{ employee.email }})
                    </option>
                </select>
            </div>

            <div>
                <label class="label cursor-pointer justify-start gap-2">
                    <input v-model="addMemberForm.is_manager" type="checkbox" class="checkbox checkbox-primary border-2" />
                    <span class="label-text font-bold uppercase text-xs">Manager de l'équipe</span>
                </label>
            </div>
        </div>

        <template #footer>
            <button class="brutal-btn" @click="showAddMemberModal = false; resetAddMemberForm()">
                Annuler
            </button>
            <button class="brutal-btn brutal-btn-success" @click="addMember">
                Ajouter
            </button>
        </template>
    </BaseModal>

    <ConfirmDialog v-model="showDeleteTeamConfirm" title="Supprimer l'équipe"
        :message="`Êtes-vous sûr de vouloir supprimer l'équipe <span class='font-bold'>${selectedTeam?.name}</span> ?`"
        confirm-text="Supprimer l'équipe" cancel-text="Annuler" variant="error" @confirm="deleteTeam" />

    <ConfirmDialog v-model="showRemoveMemberConfirm" title="Retirer le membre"
        :message="`Êtes-vous sûr de vouloir retirer <span class='font-bold'>${memberToRemove?.first_name} ${memberToRemove?.last_name}</span> ?`"
        confirm-text="Retirer" cancel-text="Annuler" variant="warning" @confirm="removeMember" />
</template>