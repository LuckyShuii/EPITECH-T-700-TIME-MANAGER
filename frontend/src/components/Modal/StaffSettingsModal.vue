<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Employee, EmployeeUpdateData } from '@/types/Employee'
import userAPI from '@/services/routers/UserAPI'
import { useNotificationsStore } from '@/store/NotificationsStore'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import BaseModal from '@/components/Modal/BaseModal.vue'
import type { WeeklyRate } from '@/types/WeeklyRate'
import weeklyRateAPI from '@/services/routers/WeeklyRateAPI'

const modelValue = defineModel<boolean>({ default: false })
const notificationsStore = useNotificationsStore()

const showDeleteConfirm = ref(false)
const employeeToDelete = ref<Employee | null>(null)
const isEditingNames = ref(false)
const employees = ref<Employee[]>([])
const isLoading = ref(false)
const weeklyRates = ref<WeeklyRate[]>([])
const isLoadingRates = ref(false)
const weekStartDay = ref('lundi')
const selectedEmployee = ref<Employee | null>(null)
const editForm = ref<EmployeeUpdateData | null>(null)

const loadEmployees = async () => {
  isLoading.value = true
  try {
    const response = await userAPI.getAll()
    employees.value = response.data
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

const loadWeeklyRates = async () => {
  isLoadingRates.value = true
  try {
    const response = await weeklyRateAPI.getAll()
    weeklyRates.value = response.data
  } catch (error) {
    console.error('Erreur lors du chargement des taux hebdomadaires:', error)
    notificationsStore.addNotification({
      status: 'error',
      title: 'Erreur de chargement',
      description: 'Impossible de charger les taux hebdomadaires'
    })
  } finally {
    isLoadingRates.value = false
  }
}

watch(() => modelValue.value, (newValue) => {
  if (newValue) {
    loadEmployees()
    loadWeeklyRates()
  }
})

const selectEmployee = (employee: Employee) => {
  selectedEmployee.value = employee
  isEditingNames.value = false

  const matchingRate = weeklyRates.value.find(
    rate => rate.amount === employee.weekly_rate
  )

  editForm.value = {
    first_name: employee.first_name,
    last_name: employee.last_name,
    email: employee.email,
    phone_number: employee.phone_number,
    roles: [...employee.roles],
    weekly_rate_uuid: matchingRate?.uuid || '',
    status: employee.status,
    username: employee.username
  }
}

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

const confirmDelete = () => {
  if (!selectedEmployee.value) return
  employeeToDelete.value = selectedEmployee.value
  showDeleteConfirm.value = true
}

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
  modelValue.value = false
  selectedEmployee.value = null
  editForm.value = null
}

const weekDaysOptions = [
  'lundi', 'mardi', 'mercredi', 'jeudi', 'vendredi', 'samedi', 'dimanche'
]
</script>

<template>
  <BaseModal :model-value="modelValue" @update:model-value="modelValue = $event" title="Gestion de l'effectif">
    <div class="grid grid-cols-[300px_1fr] gap-6 h-[calc(80vh-240px)]">
      <!-- Gauche : Liste des employés -->
      <div class="border-r-2 border-black pr-4 overflow-y-auto">
        <h4 class="font-bold uppercase text-sm mb-4 tracking-wider">Employés</h4>

        <div v-if="isLoading" class="flex justify-center py-4">
          <span class="loading loading-spinner loading-md"></span>
        </div>

        <div v-else class="space-y-1">
          <div v-for="employee in employees" :key="employee.uuid" @click="selectEmployee(employee)"
            class="p-2 border-2 border-gray-400 cursor-pointer" :class="{
              'bg-black text-white border-black': selectedEmployee?.uuid === employee.uuid,
              'hover:bg-gray-100': selectedEmployee?.uuid !== employee.uuid
            }">
            <p class="font-bold text-sm">{{ employee.first_name }} {{ employee.last_name }}</p>
            <div class="flex gap-1 mt-1">
              <span class="text-xs font-bold" :class="employee.status === 'active' ? 'text-green-700' : 'text-red-700'">
                {{ employee.status === 'active' ? 'ACTIF' : 'INACTIF' }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Droite : Détails et formulaire -->
      <div v-if="selectedEmployee && editForm" class="overflow-y-auto">
        <div class="flex justify-between items-center mb-6">
          <h4 class="font-bold text-lg uppercase">
            {{ selectedEmployee.first_name }} {{ selectedEmployee.last_name }}
          </h4>
          <button @click="isEditingNames = !isEditingNames"
            class="border-2 border-black px-3 py-1 flex items-center justify-center text-xs font-bold hover:bg-gray-100"
            :class="{ 'bg-black text-white': isEditingNames }">
            {{ isEditingNames ? 'FAIT' : 'ÉDITER' }}
          </button>
        </div>

        <div class="space-y-4">
          <!-- Prénom -->
          <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
            <span class="font-bold uppercase text-xs">Prénom</span>
            <input v-if="isEditingNames" v-model="editForm.first_name" type="text"
              class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100" />
            <span v-else class="font-bold">{{ editForm.first_name }}</span>
          </div>

          <!-- Nom -->
          <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
            <span class="font-bold uppercase text-xs">Nom</span>
            <input v-if="isEditingNames" v-model="editForm.last_name" type="text"
              class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100" />
            <span v-else class="font-bold">{{ editForm.last_name }}</span>
          </div>

          <!-- Username -->
          <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
            <span class="font-bold uppercase text-xs">Username</span>
            <span class="font-bold">{{ selectedEmployee.username }}</span>
          </div>

          <div class="h-px bg-black my-4"></div>

          <!-- Email -->
          <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
              <span class="label-text font-bold uppercase text-xs">Email</span>
            </label>
            <input v-model="editForm.email" type="email"
              class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100" />
          </div>

          <!-- Téléphone -->
          <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
              <span class="label-text font-bold uppercase text-xs">Téléphone</span>
            </label>
            <input v-model="editForm.phone_number" type="tel" maxlength="10"
              class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100" />
          </div>

          <!-- Rôles -->
          <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
              <span class="label-text font-bold uppercase text-xs">Rôles</span>
            </label>
            <div>
              <div class="flex gap-4 pt-3">
                <label class="label cursor-pointer gap-2">
                  <input v-model="editForm.roles" type="checkbox" value="employee"
                    class="checkbox checkbox-primary border-2" />
                  <span class="label-text text-xs font-bold">Employé</span>
                </label>
                <label class="label cursor-pointer gap-2">
                  <input v-model="editForm.roles" type="checkbox" value="manager"
                    class="checkbox checkbox-primary border-2" />
                  <span class="label-text text-xs font-bold">Manager</span>
                </label>
                <label class="label cursor-pointer gap-2">
                  <input v-model="editForm.roles" type="checkbox" value="admin"
                    class="checkbox checkbox-primary border-2" />
                  <span class="label-text text-xs font-bold">Admin</span>
                </label>
              </div>
            </div>
          </div>

          <!-- Taux hebdomadaire -->
          <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
              <span class="label-text font-bold uppercase text-xs">Taux</span>
            </label>
            <div class="space-y-2">
              <div v-if="isLoadingRates" class="flex items-center gap-2">
                <span class="loading loading-spinner loading-sm"></span>
                <span class="text-xs">Chargement...</span>
              </div>

              <select v-else v-model="editForm.weekly_rate_uuid"
                class="select select-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100">
                <option value="" disabled>Sélectionner un taux</option>
                <option v-for="rate in weeklyRates" :key="rate.uuid" :value="rate.uuid">
                  {{ rate.amount }}h/semaine - {{ rate.rate_name }}
                </option>
              </select>

              <p v-if="selectedEmployee" class="text-xs font-bold">
                ACTUELLEMENT: {{ selectedEmployee.weekly_rate }}h - {{ selectedEmployee.weekly_rate_name }}
              </p>
            </div>
          </div>

          <!-- Début de semaine -->
          <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
              <span class="label-text font-bold uppercase text-xs">Début</span>
            </label>
            <select v-model="weekStartDay"
              class="select select-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100">
              <option v-for="day in weekDaysOptions" :key="day" :value="day">
                {{ day.toUpperCase() }}
              </option>
            </select>
          </div>

          <!-- Statut -->
          <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
              <span class="label-text font-bold uppercase text-xs">Statut</span>
            </label>
            <select v-model="editForm.status"
              class="select select-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100">
              <option value="active">ACTIF</option>
              <option value="inactive">INACTIF</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Placeholder -->
      <div v-else class="flex items-center justify-center opacity-50">
        <p class="font-bold">Sélectionnez un employé</p>
      </div>
    </div>

    <!-- Footer avec boutons d'action -->
    <template #footer>
      <button @click="confirmDelete" class="brutal-btn brutal-btn-error">
        Supprimer
      </button>
      <button @click="saveChanges" class="brutal-btn brutal-btn-success">
        Sauvegarder
      </button>
    </template>
  </BaseModal>

  <ConfirmDialog v-model="showDeleteConfirm" title="Confirmer la suppression"
    :message="`Êtes-vous sûr de vouloir supprimer <span class='font-bold'>${employeeToDelete?.first_name} ${employeeToDelete?.last_name}</span> ?`"
    confirm-text="Supprimer définitivement" cancel-text="Annuler" variant="error" @confirm="deleteEmployee" />
</template>

<style scoped>
input[type="checkbox"].checkbox-primary:checked {
  background-color: #16a34a;
  border-color: #16a34a;
}

input[type="checkbox"].checkbox-primary:checked::after {
  content: '✓';
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  color: white;
}

input[type="checkbox"].checkbox {
  width: 1.25rem;
  height: 1.25rem;
  border-radius: 0;
  border: 2px solid #000;
  cursor: pointer;
}

input[type="checkbox"].checkbox:checked {
  background-color: #16a34a;
  border-color: #16a34a;
}

input[type="checkbox"].checkbox:hover {
  background-color: #f3f4f6;
}
</style>