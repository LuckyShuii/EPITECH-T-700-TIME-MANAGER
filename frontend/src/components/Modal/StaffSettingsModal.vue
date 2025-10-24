<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Employee, EmployeeUpdateData } from '@/types/Employee'
import userAPI from '@/services/routers/UserAPI'
import { useNotificationsStore } from '@/store/NotificationsStore'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

const notificationsStore = useNotificationsStore()

interface Props {
  modelValue: boolean
}

const props = defineProps<Props>()

const showDeleteConfirm = ref(false)
const employeeToDelete = ref<Employee | null>(null)

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()


// Mode édition pour les champs sensibles
const isEditingNames = ref(false)

// État : liste d'employés (plus de mock data)
const employees = ref<Employee[]>([])
const isLoading = ref(false)

// État : employé sélectionné
const selectedEmployee = ref<Employee | null>(null)

// Formulaire d'édition
const editForm = ref<EmployeeUpdateData | null>(null)

// Charger les employés depuis l'API
const loadEmployees = async () => {
  isLoading.value = true
  try {
    const response = await userAPI.getAll()
    employees.value = response.data // Ajustez selon la structure de votre APIHandler
  } catch (error) {
    console.error('Erreur lors du chargement des employés:', error)
    notificationsStore.addNotification({
      status: 'error',
      title: 'Erreur de chargement',
      description: 'Impossible de charger la liste des employés'
    })
  } finally {
    isLoading.value = false
  }
}

// Charger les employés à chaque ouverture de la modale
watch(() => props.modelValue, (newValue) => {
  if (newValue) {
    loadEmployees()
  }
})

// Sélectionner un employé
const selectEmployee = (employee: Employee) => {
  selectedEmployee.value = employee
  isEditingNames.value = false
  editForm.value = {
    first_name: employee.first_name,
    last_name: employee.last_name,
    email: employee.email,
    phone_number: employee.phone_number,
    roles: [...employee.roles],
    weekly_hours: employee.weekly_hours,
    status: employee.status
  }
}

// Sauvegarder les modifications (seulement le statut pour l'instant)
const saveChanges = async () => {
  if (!selectedEmployee.value || !editForm.value) return

  try {
    const payload = {
      ...editForm.value,
      username: selectedEmployee.value.username
    }

    await userAPI.update(selectedEmployee.value.uuid, payload)
    notificationsStore.addNotification({
      status: 'success',
      title: 'Modifications enregistrées',
      description: 'Les informations de l\'employé ont été mises à jour'
    })
    await loadEmployees()
    closeModal()
  } catch (error) {
    console.error('Erreur lors de la mise à jour:', error)
    notificationsStore.addNotification({
      status: 'error',
      title: 'Erreur de mise à jour',
      description: 'Impossible de sauvegarder les modifications'
    })
  }
}

// Supprimer un employé
// Préparer la suppression (ouvre la modale de confirmation)
const confirmDelete = () => {
  if (!selectedEmployee.value) return
  employeeToDelete.value = selectedEmployee.value
  showDeleteConfirm.value = true
}

// Supprimer réellement après confirmation
const deleteEmployee = async () => {
  if (!employeeToDelete.value) return

  try {
    await userAPI.deleteUser(employeeToDelete.value.uuid)
    notificationsStore.addNotification({
      status: 'success',
      title: 'Employé supprimé',
      description: `${employeeToDelete.value.first_name} ${employeeToDelete.value.last_name} a été supprimé`
    })
    await loadEmployees()
    closeModal()
  } catch (error) {
    console.error('Erreur lors de la suppression:', error)
    notificationsStore.addNotification({
      status: 'error',
      title: 'Erreur de suppression',
      description: 'Impossible de supprimer l\'employé'
    })
  } finally {
    employeeToDelete.value = null
  }
}



const closeModal = () => {
  emit('update:modelValue', false)
  selectedEmployee.value = null
  editForm.value = null
}

// Options pour le dropdown horaires
const weeklyHoursOptions = [24, 28, 35, 39]

// Changer le statut (active/inactive)
const toggleStatus = async () => {
  if (!selectedEmployee.value || !editForm.value) return

  const newStatus = editForm.value.status === 'active' ? 'inactive' : 'active'

  try {
    await userAPI.updateStatus(selectedEmployee.value.uuid, newStatus)

    // Mettre à jour localement
    editForm.value.status = newStatus
    if (selectedEmployee.value) {
      selectedEmployee.value.status = newStatus
    }

    notificationsStore.addNotification({
      status: 'success',
      title: 'Statut modifié',
      description: `L'employé est maintenant ${newStatus === 'active' ? 'actif' : 'inactif'}`
    })

    await loadEmployees()
  } catch (error) {
    console.error('Erreur lors du changement de statut:', error)
    notificationsStore.addNotification({
      status: 'error',
      title: 'Erreur de modification',
      description: 'Impossible de changer le statut'
    })
  }
}


</script>

<template>
  <Transition name="modal">
    <div v-if="modelValue" class="modal modal-open" @click.self="closeModal">
      <div class="modal-box max-w-6xl h-[80vh]">
        <!-- Header -->
        <div class="flex justify-between items-center mb-6">
          <h3 class="font-bold text-2xl">Gestion de l'effectif</h3>
          <button class="btn btn-sm btn-circle btn-ghost" @click="closeModal">
            ✕
          </button>
        </div>

        <!-- Layout split : Liste + Détails -->
        <div class="grid grid-cols-[300px_1fr] gap-6 h-full">
          <!-- Gauche : Liste des employés -->
          <div class="border-r border-base-300 pr-4 overflow-y-auto">
            <h4 class="font-semibold text-lg mb-4">Employés</h4>

            <!-- Indicateur de chargement -->
            <div v-if="isLoading" class="flex justify-center py-4">
              <span class="loading loading-spinner loading-md"></span>
            </div>

            <!-- Liste des employés -->
            <div v-else class="space-y-1">
              <div v-for="employee in employees" :key="employee.uuid" @click="selectEmployee(employee)"
                class="p-2 rounded-lg cursor-pointer transition-colors" :class="{
                  'bg-primary text-primary-content': selectedEmployee?.uuid === employee.uuid,
                  'bg-base-200 hover:bg-base-300': selectedEmployee?.uuid !== employee.uuid
                }">
                <p class="font-medium text-sm">{{ employee.first_name }} {{ employee.last_name }}</p>
                <div class="flex gap-1 mt-1">
                  <span class="badge badge-xs" :class="employee.status === 'active' ? 'badge-success' : 'badge-error'">
                    {{ employee.status === 'active' ? 'Actif' : 'Inactif' }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Droite : Détails et formulaire -->
          <div v-if="selectedEmployee && editForm" class="overflow-y-auto">
            <h4 class="font-semibold text-xl mb-6">
              {{ selectedEmployee.first_name }} {{ selectedEmployee.last_name }}
            </h4>

            <div class="space-y-4">
              <!-- Infos modifiables avec icône -->
              <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                <div class="flex items-center gap-2">
                  <span class="opacity-70">Prénom</span>
                  <button @click="isEditingNames = !isEditingNames" class="btn btn-ghost btn-xs btn-circle"
                    :class="{ 'text-primary': isEditingNames }">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                      <path
                        d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                    </svg>
                  </button>
                </div>
                <input v-if="isEditingNames" v-model="editForm.first_name" type="text"
                  class="input input-bordered w-full" />
                <span v-else class="font-medium">{{ editForm.first_name }}</span>
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                <span class="opacity-70">Nom</span>
                <input v-if="isEditingNames" v-model="editForm.last_name" type="text"
                  class="input input-bordered w-full" />
                <span v-else class="font-medium">{{ editForm.last_name }}</span>
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                <span class="opacity-70">Username</span>
                <span class="font-medium">{{ selectedEmployee.username }}</span>
              </div>

              <div class="divider"></div>

              <!-- Champs modifiables -->
              <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                <label class="label">
                  <span class="label-text">Statut</span>
                </label>
                <div class="flex items-center gap-3">
                  <input type="checkbox" class="toggle toggle-success" :checked="editForm.status === 'active'"
                    @change="toggleStatus" />
                  <span class="text-sm font-medium">
                    {{ editForm.status === 'active' ? 'Actif' : 'Inactif' }}
                  </span>
                </div>
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
                <label class="label pt-3">
                  <span class="label-text">Téléphone</span>
                </label>
                <input v-model="editForm.phone_number" type="tel" maxlength="10" class="input input-bordered w-full" />
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
                <label class="label pt-3">
                  <span class="label-text">Rôles</span>
                </label>
                <div>
                  <div class="flex gap-4 pt-3">
                    <label class="label cursor-pointer gap-2">
                      <input v-model="editForm.roles" type="checkbox" value="employee"
                        class="checkbox checkbox-primary" />
                      <span class="label-text">Employé</span>
                    </label>
                    <label class="label cursor-pointer gap-2">
                      <input v-model="editForm.roles" type="checkbox" value="manager"
                        class="checkbox checkbox-primary" />
                      <span class="label-text">Manager</span>
                    </label>
                    <label class="label cursor-pointer gap-2">
                      <input v-model="editForm.roles" type="checkbox" value="admin" class="checkbox checkbox-primary" />
                      <span class="label-text">Admin</span>
                    </label>
                  </div>
                </div>
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
                <label class="label pt-3">
                  <span class="label-text">Horaire hebdo</span>
                </label>
                <select v-model.number="editForm.weekly_hours" class="select select-bordered w-full">
                  <option v-for="hours in weeklyHoursOptions" :key="hours" :value="hours">
                    {{ hours }}h / semaine
                  </option>
                </select>
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
                <label class="label pt-3">
                  <span class="label-text">Statut</span>
                </label>
                <select v-model="editForm.status" class="select select-bordered w-full">
                  <option value="active">Actif</option>
                  <option value="inactive">Inactif</option>
                </select>
              </div>

              <!-- Boutons d'action -->
              <div class="flex gap-2 justify-end pt-6">
                <button @click="confirmDelete" class="btn btn-error btn-outline">
                  Supprimer
                </button>
                <button @click="saveChanges" class="btn bg-gradient-to-br from-primary-500 to-secondary-500 text-white">
                  Sauvegarder
                </button>
              </div>
            </div>
          </div>

          <!-- Placeholder si aucun employé sélectionné -->
          <div v-else class="flex items-center justify-center opacity-50">
            <p>Sélectionnez un employé dans la liste</p>
          </div>
        </div>
      </div>
    </div>
  </Transition>
  <!-- Dialog de confirmation de suppression -->
  <ConfirmDialog v-model="showDeleteConfirm" title="Confirmer la suppression"
    :message="`Êtes-vous sûr de vouloir supprimer <span class='font-bold'>${employeeToDelete?.first_name} ${employeeToDelete?.last_name}</span> ?<br><span class='text-sm opacity-70'>Cette action est irréversible.</span>`"
    confirm-text="Supprimer définitivement" cancel-text="Annuler" variant="error" @confirm="deleteEmployee" />
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