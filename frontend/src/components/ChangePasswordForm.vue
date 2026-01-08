<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import type { PasswordChangePayload } from '@/types/PasswordReset'
import API from '@/services/API'
import { useNotificationsStore } from '@/store/NotificationsStore'

const router = useRouter()
const route = useRoute()
const notificationsStore = useNotificationsStore()

const userUuid = ref<string>('')
const formData = ref({
  password: '',
  passwordConfirm: ''
})
const errors = ref<{ [key: string]: string }>({})
const isSubmitting = ref(false)

// Validation du mot de passe
const passwordValidation = computed(() => ({
  minLength: formData.value.password.length >= 8,
  hasUppercase: /[A-Z]/.test(formData.value.password),
  hasNumber: /\d/.test(formData.value.password),
  matches: formData.value.password === formData.value.passwordConfirm && formData.value.password.length > 0
}))

// Vérifier si tous les critères sont remplis
const isPasswordValid = computed(() => {
  return Object.values(passwordValidation.value).every(v => v)
})

onMounted(() => {
  // Extraire l'UUID de l'URL
  const publicKey = route.query.user_public_key as string
  if (!publicKey) {
    notificationsStore.addNotification({
      status: 'error',
      title: 'Lien invalide',
      description: 'Le lien de réinitialisation n\'est pas valide'
    })
    // Rediriger après 2 secondes
    setTimeout(() => {
      router.push('/')
    }, 2000)
    return
  }
  userUuid.value = publicKey
})

const handleSubmit = async () => {
  errors.value = {}

  if (!isPasswordValid.value) {
    errors.value.password = 'Le mot de passe ne respecte pas les critères'
    return
  }

  if (formData.value.password !== formData.value.passwordConfirm) {
    errors.value.passwordConfirm = 'Les mots de passe ne correspondent pas'
    return
  }

  isSubmitting.value = true

  try {
    const payload: PasswordChangePayload = {
      new_password: formData.value.password,
      user_uuid: userUuid.value
    }
    await API.userAPI.changePassword(payload)
    
    notificationsStore.addNotification({
      status: 'success',
      title: 'Mot de passe changé',
      description: 'Votre mot de passe a été modifié avec succès. Redirection vers la connexion...'
    })

    // Rediriger vers login après 3 secondes
    setTimeout(() => {
      router.push('/login')
    }, 3000)
  } catch (error: any) {
    const errorMessage = error.response?.data?.message || 'Erreur lors du changement de mot de passe'
    notificationsStore.addNotification({
      status: 'error',
      title: 'Erreur',
      description: errorMessage
    })
    errors.value.general = errorMessage
  } finally {
    isSubmitting.value = false
  }
}

const resetForm = () => {
  formData.value = {
    password: '',
    passwordConfirm: ''
  }
  errors.value = {}
}

const handleCancel = () => {
  resetForm()
  router.push('/login')
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-6">
    <!-- Erreur générale -->
    <div v-if="errors.general" class="brutal-container bg-red-50 border-red-700">
      <p class="text-red-700 font-bold">ERREUR</p>
      <p class="text-red-600 text-sm mt-1">{{ errors.general }}</p>
    </div>

    <!-- Mot de passe -->
    <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
      <label class="label pt-3">
        <span class="label-text font-bold uppercase text-xs">Mot de passe</span>
      </label>
      <div class="w-full">
        <input 
          v-model="formData.password" 
          type="password" 
          placeholder="••••••••"
          class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100"
          :class="{ 'border-red-700': errors.password }" 
        />
        <label v-if="errors.password" class="label">
          <span class="label-text-alt text-red-700 font-bold text-xs">{{ errors.password }}</span>
        </label>

        <!-- Critères de validation -->
        <div class="mt-4 space-y-2 p-4 border-2 border-black bg-gray-50">
          <p class="font-bold uppercase text-xs mb-3">Critères requis :</p>
          
          <!-- Critère : 8 caractères -->
          <div class="flex items-center gap-2">
            <div 
              class="w-4 h-4 border-2 border-black flex items-center justify-center"
              :class="passwordValidation.minLength ? 'bg-black' : 'bg-gray-200'"
            >
              <span v-if="passwordValidation.minLength" class="text-white font-bold text-xs">✓</span>
            </div>
            <span 
              class="text-xs font-bold uppercase"
              :class="passwordValidation.minLength ? 'text-black' : 'text-gray-400'"
            >
              8 caractères minimum
            </span>
          </div>

          <!-- Critère : Majuscule -->
          <div class="flex items-center gap-2">
            <div 
              class="w-4 h-4 border-2 border-black flex items-center justify-center"
              :class="passwordValidation.hasUppercase ? 'bg-black' : 'bg-gray-200'"
            >
              <span v-if="passwordValidation.hasUppercase" class="text-white font-bold text-xs">✓</span>
            </div>
            <span 
              class="text-xs font-bold uppercase"
              :class="passwordValidation.hasUppercase ? 'text-black' : 'text-gray-400'"
            >
              Une lettre majuscule
            </span>
          </div>

          <!-- Critère : Chiffre -->
          <div class="flex items-center gap-2">
            <div 
              class="w-4 h-4 border-2 border-black flex items-center justify-center"
              :class="passwordValidation.hasNumber ? 'bg-black' : 'bg-gray-200'"
            >
              <span v-if="passwordValidation.hasNumber" class="text-white font-bold text-xs">✓</span>
            </div>
            <span 
              class="text-xs font-bold uppercase"
              :class="passwordValidation.hasNumber ? 'text-black' : 'text-gray-400'"
            >
              Un chiffre
            </span>
          </div>

          <!-- Critère : Confirmation -->
          <div class="flex items-center gap-2">
            <div 
              class="w-4 h-4 border-2 border-black flex items-center justify-center"
              :class="passwordValidation.matches ? 'bg-black' : 'bg-gray-200'"
            >
              <span v-if="passwordValidation.matches" class="text-white font-bold text-xs">✓</span>
            </div>
            <span 
              class="text-xs font-bold uppercase"
              :class="passwordValidation.matches ? 'text-black' : 'text-gray-400'"
            >
              Les mots de passe correspondent
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Confirmation mot de passe -->
    <div class="grid grid-cols-[150px_1fr] gap-2 items-start">
      <label class="label pt-3">
        <span class="label-text font-bold uppercase text-xs">Confirmation</span>
      </label>
      <div>
        <input 
          v-model="formData.passwordConfirm" 
          type="password" 
          placeholder="••••••••"
          class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100"
          :class="{ 'border-red-700': errors.passwordConfirm }" 
        />
        <label v-if="errors.passwordConfirm" class="label">
          <span class="label-text-alt text-red-700 font-bold text-xs">{{ errors.passwordConfirm }}</span>
        </label>
      </div>
    </div>

    <!-- Boutons -->
    <div class="flex gap-2 justify-end pt-4">
      <button 
        type="button" 
        class="brutal-btn" 
        @click="handleCancel" 
        :disabled="isSubmitting"
      >
        Annuler
      </button>
      <button 
        type="submit" 
        class="brutal-btn brutal-btn-success" 
        :disabled="isSubmitting || !isPasswordValid"
      >
        {{ isSubmitting ? 'Changement...' : 'Changer le mot de passe' }}
      </button>
    </div>
  </form>
</template>

<style scoped>
button {
  transition: none;
}
</style>