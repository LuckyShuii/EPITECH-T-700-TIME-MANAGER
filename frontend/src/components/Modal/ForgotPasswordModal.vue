<script setup lang="ts">
import { ref } from 'vue'
import API from '@/services/API'
import { useNotificationsStore } from '@/store/NotificationsStore'

const notificationsStore = useNotificationsStore()

const emit = defineEmits<{
  close: []
}>()

const email = ref<string>('')
const isLoading = ref(false)
const isSuccess = ref(false)
const errors = ref<{ [key: string]: string }>({})

const validateEmail = (): boolean => {
  errors.value = {}
  
  if (!email.value.trim()) {
    errors.value.email = 'L\'email est requis'
    return false
  }
  
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(email.value)) {
    errors.value.email = 'Format d\'email invalide'
    return false
  }
  
  return true
}

const handleSubmit = async () => {
  if (!validateEmail()) {
    return
  }
  
  isLoading.value = true
  
  try {
    await API.userAPI.sendPasswordResetEmail(email.value)
    
    notificationsStore.addNotification({
      status: 'success',
      title: 'Email envoyé',
      description: 'Vérifiez votre email pour réinitialiser votre mot de passe'
    })
    
    isSuccess.value = true
    
    // Fermer la modal après 2 secondes
    setTimeout(() => {
      handleClose()
    }, 1000)
  } catch (error: any) {
    const errorMessage = error.response?.data?.message || 'Erreur lors de l\'envoi de l\'email'
    notificationsStore.addNotification({
      status: 'error',
      title: 'Erreur',
      description: errorMessage
    })
    errors.value.general = errorMessage
  } finally {
    isLoading.value = false
  }
}

const handleClose = () => {
  email.value = ''
  errors.value = {}
  isSuccess.value = false
  emit('close')
}
</script>

<template>
  <div class="fixed inset-0 bg-white bg-opacity-100 flex items-center justify-center z-50 p-4">
    <div class="bg-white border-2 border-black w-full max-w-md">
      <!-- Header -->
      <div class="border-b-2 border-black p-6">
        <h2 class="text-2xl font-black uppercase tracking-tight">MOT DE PASSE OUBLIÉ</h2>
      </div>

      <!-- Content -->
      <div class="p-6 space-y-4">
        <!-- Success Message -->
        <div v-if="isSuccess" class="border-2 border-green-600 bg-green-50 p-4">
          <p class="font-black uppercase text-green-700 text-sm tracking-wider">SUCCÈS</p>
          <p class="text-xs font-bold mt-2">Vérifiez votre email pour réinitialiser votre mot de passe.</p>
        </div>

        <!-- Form -->
        <template v-else>
          <!-- Error -->
          <div v-if="errors.general" class="border-2 border-red-700 bg-red-50 p-4">
            <p class="text-red-700 font-bold uppercase text-xs">ERREUR</p>
            <p class="text-red-600 text-xs font-bold mt-1">{{ errors.general }}</p>
          </div>

          <!-- Email Input -->
          <div>
            <label class="label pt-0">
              <span class="label-text font-bold uppercase text-xs">Email</span>
            </label>
            <input
              v-model="email"
              type="email"
              placeholder="email@example.com"
              class="input input-bordered w-full border-2 border-black rounded-none !bg-white !text-black hover:bg-gray-100"
              :class="{ 'border-red-700': errors.email }"
              :disabled="isLoading"
            />
            <label v-if="errors.email" class="label">
              <span class="label-text-alt text-red-700 font-bold text-xs">{{ errors.email }}</span>
            </label>
          </div>

          <!-- Info -->
          <p class="text-xs font-bold text-gray-600 uppercase tracking-widest">
            Un lien de réinitialisation vous sera envoyé par email.
          </p>
        </template>
      </div>

      <!-- Footer / Buttons -->
      <div class="border-t-2 border-black p-4 flex gap-2">
        <button
          type="button"
          @click="handleClose"
          :disabled="isLoading"
          class="brutal-btn flex-1"
        >
          Fermer
        </button>
        <button
          v-if="!isSuccess"
          type="button"
          @click="handleSubmit"
          :disabled="isLoading"
          class="brutal-btn brutal-btn-success flex-1"
        >
          {{ isLoading ? 'Envoi...' : 'Envoyer' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
button {
  transition: none;
}
</style>