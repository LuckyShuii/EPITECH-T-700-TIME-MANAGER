<script setup lang="ts">
import { ref, watch } from 'vue'
import type { RegisterFormData, RegisterFormErrors } from '@/types/RegisterForm'
import API from '@/services/API'
import { useNotificationsStore } from '@/store/NotificationsStore'

const notificationsStore = useNotificationsStore()

const emit = defineEmits<{
    success: []
    cancel: []
}>()

const formData = ref<RegisterFormData>({
    first_name: '',
    last_name: '',
    email: '',
    username: '',
    phone_number: '',
    roles: [],
    password: ''
})

const errors = ref<RegisterFormErrors>({})
const isSubmitting = ref(false)

const generateUsername = () => {
    const firstName = formData.value.first_name.trim()
    const lastName = formData.value.last_name.trim()

    if (firstName && lastName) {
        const username = (firstName[0] + lastName).toLowerCase()
        formData.value.username = username
    }
}

watch([() => formData.value.first_name, () => formData.value.last_name], () => {
    generateUsername()
})

const validateForm = (): boolean => {
    errors.value = {}
    let isValid = true

    if (!formData.value.first_name.trim()) {
        errors.value.first_name = 'Le prénom est requis'
        isValid = false
    }

    if (!formData.value.last_name.trim()) {
        errors.value.last_name = 'Le nom est requis'
        isValid = false
    }

    if (!formData.value.email.trim()) {
        errors.value.email = "L'email est requis"
        isValid = false
    } else {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        if (!emailRegex.test(formData.value.email)) {
            errors.value.email = "Format d'email invalide"
            isValid = false
        }
    }

    if (!formData.value.username.trim()) {
        errors.value.username = "Le nom d'utilisateur est requis"
        isValid = false
    }

    if (!formData.value.phone_number.trim()) {
        errors.value.phone_number = 'Le numéro de téléphone est requis'
        isValid = false
    } else {
        const phoneRegex = /^\d{10}$/
        if (!phoneRegex.test(formData.value.phone_number)) {
            errors.value.phone_number = 'Le numéro doit contenir 10 chiffres'
            isValid = false
        }
    }

    if (formData.value.roles.length === 0) {
        errors.value.roles = 'Sélectionnez au moins un rôle'
        isValid = false
    }

    if (!formData.value.password) {
        errors.value.password = 'Le mot de passe est requis'
        isValid = false
    }

    return isValid
}

const handleSubmit = async () => {
  if (!validateForm()) {
    return
  }

  isSubmitting.value = true
  errors.value = {}

  try {
    // Étape 1 : Créer l'utilisateur
    await API.userAPI.register(formData.value)
    
    // Étape 2 : Envoyer l'email de réinitialisation
    try {
      await API.userAPI.sendPasswordResetEmail(formData.value.email)
      notificationsStore.addNotification({
        status: 'success',
        title: 'Employé créé',
        description: 'Le nouvel employé a été ajouté avec succès et un email de réinitialisation a été envoyé'
      })
    } catch (emailError: any) {
      // L'utilisateur est créé mais l'email n'a pas pu être envoyé
      notificationsStore.addNotification({
        status: 'warning',
        title: 'Employé créé, mais erreur email',
        description: 'L\'utilisateur a été créé mais l\'email de réinitialisation n\'a pas pu être envoyé'
      })
    }
    
    emit('success')
    resetForm()
  } catch (error: any) {
    // L'utilisateur n'a pas pu être créé
    const errorMessage = error.response?.data?.message || 'Erreur lors de la création'
    notificationsStore.addNotification({
      status: 'error',
      title: 'Erreur de création',
      description: errorMessage
    })
    errors.value.general = errorMessage
  } finally {
    isSubmitting.value = false
  }
}

const resetForm = () => {
    formData.value = {
        first_name: '',
        last_name: '',
        email: '',
        username: '',
        phone_number: '',
        roles: [],
        password: ''
    }
    errors.value = {}
}

const handleCancel = () => {
    resetForm()
    emit('cancel')
}
</script>

<template>
    <form @submit.prevent="handleSubmit" class="space-y-4">
        <!-- Erreur générale -->
        <div v-if="errors.general" class="brutal-container bg-red-50 border-red-700">
            <p class="text-red-700 font-bold">ERREUR</p>
            <p class="text-red-600 text-sm mt-1">{{ errors.general }}</p>
        </div>

        <!-- Prénom -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text font-bold uppercase text-xs">Prénom</span>
            </label>
            <div>
                <input v-model="formData.first_name" type="text" placeholder="Max" 
                    class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100" 
                    :class="{ 'border-red-700': errors.first_name }" />
                <label v-if="errors.first_name" class="label">
                    <span class="label-text-alt text-red-700 font-bold text-xs">{{ errors.first_name }}</span>
                </label>
            </div>
        </div>

        <!-- Nom -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text font-bold uppercase text-xs">Nom</span>
            </label>
            <div>
                <input v-model="formData.last_name" type="text" placeholder="Loris" 
                    class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100"
                    :class="{ 'border-red-700': errors.last_name }" />
                <label v-if="errors.last_name" class="label">
                    <span class="label-text-alt text-red-700 font-bold text-xs">{{ errors.last_name }}</span>
                </label>
            </div>
        </div>

        <!-- Email -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text font-bold uppercase text-xs">Email</span>
            </label>
            <div>
                <input v-model="formData.email" type="email" placeholder="email@example.com"
                    class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100" 
                    :class="{ 'border-red-700': errors.email }" />
                <label v-if="errors.email" class="label">
                    <span class="label-text-alt text-red-700 font-bold text-xs">{{ errors.email }}</span>
                </label>
            </div>
        </div>

        <!-- Username -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text font-bold uppercase text-xs">Username</span>
            </label>
            <div>
                <input disabled v-model="formData.username" type="text" 
                    class="input input-bordered w-full border-2 border-black rounded-none opacity-50 !bg-white !text-black"
                    :class="{ 'border-red-700': errors.username }" />
                <label v-if="errors.username" class="label">
                    <span class="label-text-alt text-red-700 font-bold text-xs">{{ errors.username }}</span>
                </label>
            </div>
        </div>

        <!-- Téléphone -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text font-bold uppercase text-xs">Téléphone</span>
            </label>
            <div>
                <input v-model="formData.phone_number" type="tel" placeholder="0123456789" maxlength="10"
                    class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100" 
                    :class="{ 'border-red-700': errors.phone_number }" />
                <label v-if="errors.phone_number" class="label">
                    <span class="label-text-alt text-red-700 font-bold text-xs">{{ errors.phone_number }}</span>
                </label>
            </div>
        </div>

        <!-- Rôles -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text font-bold uppercase text-xs">Rôles</span>
            </label>
            <div>
                <div class="flex gap-4 pt-3">
                    <label class="label cursor-pointer gap-2">
                        <input v-model="formData.roles" type="checkbox" value="employee"
                            class="checkbox checkbox-primary border-2" />
                        <span class="label-text">Employé</span>
                    </label>
                    <label class="label cursor-pointer gap-2">
                        <input v-model="formData.roles" type="checkbox" value="manager"
                            class="checkbox checkbox-primary border-2" />
                        <span class="label-text">Manager</span>
                    </label>
                    <label class="label cursor-pointer gap-2">
                        <input v-model="formData.roles" type="checkbox" value="admin" 
                            class="checkbox checkbox-primary border-2" />
                        <span class="label-text">Admin</span>
                    </label>
                </div>
                <label v-if="errors.roles" class="label">
                    <span class="label-text-alt text-red-700 font-bold text-xs">{{ errors.roles }}</span>
                </label>
            </div>
        </div>

        <!-- Mot de passe -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text font-bold uppercase text-xs">Mot de passe</span>
            </label>
            <div>
                <input v-model="formData.password" type="password" placeholder="••••••••"
                    class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100" 
                    :class="{ 'border-red-700': errors.password }" />
                <label v-if="errors.password" class="label">
                    <span class="label-text-alt text-red-700 font-bold text-xs">{{ errors.password }}</span>
                </label>
            </div>
        </div>

        <!-- Boutons -->
        <div class="flex gap-2 justify-end pt-4">
            <button type="button" class="brutal-btn" @click="handleCancel" :disabled="isSubmitting">
                Annuler
            </button>
            <button type="submit" class="brutal-btn brutal-btn-success" :disabled="isSubmitting">
                {{ isSubmitting ? 'Création...' : 'Créer l\'employé' }}
            </button>
        </div>
    </form>
</template>