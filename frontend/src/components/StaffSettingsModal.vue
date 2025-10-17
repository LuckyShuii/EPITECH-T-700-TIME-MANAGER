<script setup lang="ts">
import { ref } from 'vue'
import type { Employee, EmployeeUpdateData } from '@/types/Employee'

interface Props {
  modelValue: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

// Mock data : liste d'employés
const employees = ref<Employee[]>([
  {
    uuid: '1',
    first_name: 'Max',
    last_name: 'Loris',
    email: 'max.loris@test.com',
    username: 'mloris',
    phone_number: '0123456789',
    roles: ['employee', 'admin'],
    weekly_hours: 35,
    status: 'active'
  },
  {
    uuid: '2',
    first_name: 'Sophie',
    last_name: 'Martin',
    email: 'sophie.martin@test.com',
    username: 'smartin',
    phone_number: '0234567890',
    roles: ['employee', 'manager'],
    weekly_hours: 39,
    status: 'active'
  },
  {
    uuid: '3',
    first_name: 'Jean',
    last_name: 'Dupont',
    email: 'jean.dupont@test.com',
    username: 'jdupont',
    phone_number: '0345678901',
    roles: ['employee'],
    weekly_hours: 35,
    status: 'inactive'
  }
])

// État : employé sélectionné
const selectedEmployee = ref<Employee | null>(null)

// Formulaire d'édition (copie de l'employé sélectionné)
const editForm = ref<EmployeeUpdateData | null>(null)

// Sélectionner un employé
const selectEmployee = (employee: Employee) => {
  selectedEmployee.value = employee
  // Copier les données modifiables dans le formulaire
  editForm.value = {
    email: employee.email,
    phone_number: employee.phone_number,
    roles: [...employee.roles],
    weekly_hours: employee.weekly_hours,
    status: employee.status
  }
}

// Sauvegarder les modifications
const saveChanges = async () => {
  if (!selectedEmployee.value || !editForm.value) return
  
  // TODO: Appel API PATCH /api/users/:uuid
  console.log('Sauvegarde pour:', selectedEmployee.value.uuid, editForm.value)
  alert('✅ Modifications sauvegardées (placeholder)')
  
  closeModal()
}

// Supprimer un employé
const deleteEmployee = async () => {
  if (!selectedEmployee.value) return
  
  const confirmed = confirm(`Supprimer ${selectedEmployee.value.first_name} ${selectedEmployee.value.last_name} ?`)
  if (!confirmed) return
  
  // TODO: Appel API DELETE /api/users/:uuid
  console.log('Suppression de:', selectedEmployee.value.uuid)
  alert('🗑️ Employé supprimé (placeholder)')
  
  closeModal()
}

const closeModal = () => {
  emit('update:modelValue', false)
  selectedEmployee.value = null
  editForm.value = null
}

// Options pour le dropdown horaires
const weeklyHoursOptions = [24, 28, 35, 39]
</script>

<template>
  <Transition name="modal">
    <div 
      v-if="modelValue"
      class="modal modal-open"
      @click.self="closeModal"
    >
      <div class="modal-box max-w-6xl h-[80vh]">
        <!-- Header -->
        <div class="flex justify-between items-center mb-6">
          <h3 class="font-bold text-2xl">Gestion de l'effectif</h3>
          <button 
            class="btn btn-sm btn-circle btn-ghost"
            @click="closeModal"
          >
            ✕
          </button>
        </div>

        <!-- Layout split : Liste + Détails -->
        <div class="grid grid-cols-[300px_1fr] gap-6 h-full">
          <!-- Gauche : Liste des employés -->
          <div class="border-r border-base-300 pr-4 overflow-y-auto">
            <h4 class="font-semibold text-lg mb-4">Employés</h4>
            <div class="space-y-2">
              <div
                v-for="employee in employees"
                :key="employee.uuid"
                @click="selectEmployee(employee)"
                class="p-3 rounded-lg cursor-pointer transition-colors"
                :class="{
                  'bg-primary text-primary-content': selectedEmployee?.uuid === employee.uuid,
                  'bg-base-200 hover:bg-base-300': selectedEmployee?.uuid !== employee.uuid
                }"
              >
                <p class="font-medium">{{ employee.first_name }} {{ employee.last_name }}</p>
                <p class="text-sm opacity-75">{{ employee.username }}</p>
                <div class="flex gap-1 mt-1">
                  <span 
                    class="badge badge-sm"
                    :class="employee.status === 'active' ? 'badge-success' : 'badge-error'"
                  >
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
              <!-- Infos non-modifiables -->
              <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                <span class="opacity-70">Prénom</span>
                <span class="font-medium">{{ selectedEmployee.first_name }}</span>
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                <span class="opacity-70">Nom</span>
                <span class="font-medium">{{ selectedEmployee.last_name }}</span>
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-center">
                <span class="opacity-70">Username</span>
                <span class="font-medium">{{ selectedEmployee.username }}</span>
              </div>

              <div class="divider"></div>

              <!-- Champs modifiables -->
              <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
                <label class="label pt-3">
                  <span class="label-text">Email</span>
                </label>
                <input 
                  v-model="editForm.email"
                  type="email"
                  class="input input-bordered w-full"
                />
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
                <label class="label pt-3">
                  <span class="label-text">Téléphone</span>
                </label>
                <input 
                  v-model="editForm.phone_number"
                  type="tel"
                  maxlength="10"
                  class="input input-bordered w-full"
                />
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
                <label class="label pt-3">
                  <span class="label-text">Rôles</span>
                </label>
                <div>
                  <div class="flex gap-4 pt-3">
                    <label class="label cursor-pointer gap-2">
                      <input 
                        v-model="editForm.roles"
                        type="checkbox"
                        value="employee"
                        class="checkbox checkbox-primary"
                      />
                      <span class="label-text">Employé</span>
                    </label>
                    <label class="label cursor-pointer gap-2">
                      <input 
                        v-model="editForm.roles"
                        type="checkbox"
                        value="manager"
                        class="checkbox checkbox-primary"
                      />
                      <span class="label-text">Manager</span>
                    </label>
                    <label class="label cursor-pointer gap-2">
                      <input 
                        v-model="editForm.roles"
                        type="checkbox"
                        value="admin"
                        class="checkbox checkbox-primary"
                      />
                      <span class="label-text">Admin</span>
                    </label>
                  </div>
                </div>
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
                <label class="label pt-3">
                  <span class="label-text">Horaire hebdo</span>
                </label>
                <select 
                  v-model.number="editForm.weekly_hours"
                  class="select select-bordered w-full"
                >
                  <option v-for="hours in weeklyHoursOptions" :key="hours" :value="hours">
                    {{ hours }}h / semaine
                  </option>
                </select>
              </div>

              <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
                <label class="label pt-3">
                  <span class="label-text">Statut</span>
                </label>
                <select 
                  v-model="editForm.status"
                  class="select select-bordered w-full"
                >
                  <option value="active">Actif</option>
                  <option value="inactive">Inactif</option>
                </select>
              </div>

              <!-- Boutons d'action -->
              <div class="flex gap-2 justify-end pt-6">
                <button 
                  @click="deleteEmployee"
                  class="btn btn-error btn-outline"
                >
                  Supprimer
                </button>
                <button 
                  @click="saveChanges"
                  class="btn bg-gradient-to-br from-primary-500 to-secondary-500 text-white"
                >
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