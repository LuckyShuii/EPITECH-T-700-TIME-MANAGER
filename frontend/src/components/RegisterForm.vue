<script setup lang="ts">
import { ref, watch } from 'vue'
import type { RegisterFormData, RegisterFormErrors } from '@/types/RegisterForm'
import API from '@/services/API'

const emit = defineEmits<{
    success: []
    cancel: []
}>()

// État du formulaire
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

// Génération automatique du username
const generateUsername = () => {
    const firstName = formData.value.first_name.trim()
    const lastName = formData.value.last_name.trim()

    if (firstName && lastName) {
        const username = (firstName[0] + lastName).toLowerCase()
        formData.value.username = username
    }
}

// Watch pour générer le username automatiquement
watch([() => formData.value.first_name, () => formData.value.last_name], () => {
    generateUsername()
})

// Validation du formulaire
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

// Soumission du formulaire
const handleSubmit = async () => {
    if (!validateForm()) {
        return
    }

    isSubmitting.value = true
    errors.value = {}

    try {
        await API.userAPI.register(formData.value)
        alert('✅ Employé créé avec succès !')
        emit('success')
        resetForm()
    } catch (error: any) {
        const errorMessage = error.response?.data?.message || 'Erreur lors de la création'
        alert('❌ ' + errorMessage)
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
        <div v-if="errors.general" class="alert alert-error">
            <span>{{ errors.general }}</span>
        </div>

        <!-- Prénom -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text">Prénom *</span>
            </label>
            <div>
                <input v-model="formData.first_name" type="text" placeholder="Max" class="input input-bordered w-full"
                    :class="{ 'input-error': errors.first_name }" />
                <label v-if="errors.first_name" class="label">
                    <span class="label-text-alt text-error">{{ errors.first_name }}</span>
                </label>
            </div>
        </div>

        <!-- Nom -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text">Nom *</span>
            </label>
            <div>
                <input v-model="formData.last_name" type="text" placeholder="Nom de famille"
                    class="input input-bordered w-full" :class="{ 'input-error': errors.last_name }" />
                <label v-if="errors.last_name" class="label">
                    <span class="label-text-alt text-error">{{ errors.last_name }}</span>
                </label>
            </div>
        </div>

        <!-- Email -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text">Email *</span>
            </label>
            <div>
                <input v-model="formData.email" type="email" placeholder="email@email.com"
                    class="input input-bordered w-full" :class="{ 'input-error': errors.email }" />
                <label v-if="errors.email" class="label">
                    <span class="label-text-alt text-error">{{ errors.email }}</span>
                </label>
            </div>
        </div>

        <!-- Username -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text">Nom d'utilisateur *</span>
            </label>
            <div>
                <input disabled v-model="formData.username" type="text" class="input input-bordered w-full"
                    :class="{ 'input-error': errors.username }" />
                <label v-if="errors.username" class="label">
                    <span class="label-text-alt text-error">{{ errors.username }}</span>
                </label>
            </div>
        </div>

        <!-- Téléphone -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text">Téléphone *</span>
            </label>
            <div>
                <input v-model="formData.phone_number" type="tel" placeholder="0123456789" maxlength="10"
                    class="input input-bordered w-full" :class="{ 'input-error': errors.phone_number }" />
                <label v-if="errors.phone_number" class="label">
                    <span class="label-text-alt text-error">{{ errors.phone_number }}</span>
                </label>
            </div>
        </div>

        <!-- Rôles -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text">Rôles *</span>
            </label>
            <div>
                <div class="flex gap-4 pt-3">
                    <label class="label cursor-pointer gap-2">
                        <input v-model="formData.roles" type="checkbox" value="employee"
                            class="checkbox checkbox-primary" />
                        <span class="label-text">Employé</span>
                    </label>
                    <label class="label cursor-pointer gap-2">
                        <input v-model="formData.roles" type="checkbox" value="manager"
                            class="checkbox checkbox-primary" />
                        <span class="label-text">Manager</span>
                    </label>
                    <label class="label cursor-pointer gap-2">
                        <input v-model="formData.roles" type="checkbox" value="admin"
                            class="checkbox checkbox-primary" />
                        <span class="label-text">Admin</span>
                    </label>
                </div>
                <label v-if="errors.roles" class="label">
                    <span class="label-text-alt text-error">{{ errors.roles }}</span>
                </label>
            </div>
        </div>

        <!-- Mot de passe -->
        <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
            <label class="label pt-3">
                <span class="label-text">Mot de passe *</span>
            </label>
            <div>
                <input v-model="formData.password" type="password" placeholder="••••••••"
                    class="input input-bordered w-full" :class="{ 'input-error': errors.password }" />
                <label v-if="errors.password" class="label">
                    <span class="label-text-alt text-error">{{ errors.password }}</span>
                </label>
            </div>
        </div>

        <!-- Boutons -->
        <div class="flex gap-2 justify-end pt-4">
            <button type="button" class="btn btn-ghost" @click="handleCancel" :disabled="isSubmitting">
                Annuler
            </button>
            <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
                <span v-if="isSubmitting" class="loading loading-spinner"></span>
                {{ isSubmitting ? 'Création...' : 'Créer l\'employé' }}
            </button>
        </div>
    </form>
</template>