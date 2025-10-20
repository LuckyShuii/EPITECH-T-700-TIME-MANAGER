<script setup lang="ts">
import { reactive, ref } from 'vue';
import type { UserLogin } from "@/types/userType";
import { useAuthStore } from '@/store/AuthStore';
import router from '@/router';
import { useNotificationsStore } from '@/store/notificationsStore';

const loading = ref<boolean>(false)
const AuthStore = useAuthStore()
const notificationsStore = useNotificationsStore()

const form = reactive<UserLogin>({
  username: '',
  password: ''
})

const getPayload = () => {
  return {
    password: form.password ?? '',
    username: form.username ?? ''
  }
}

const handleSubmit = async () => {
  loading.value = true
  try {
    await AuthStore.login(getPayload())
    
    // Notification de succès uniquement si login réussi
    notificationsStore.addNotification({
      status: 'success',
      title: 'Connexion réussie',
      description: 'Vous êtes maintenant connecté'
    })
    
    router.push('/dashboard')
  }
  catch (error) {
    // Notification d'erreur uniquement si login échoué
    notificationsStore.addNotification({
      status: 'error',
      title: 'Erreur de connexion',
      description: 'Identifiants incorrects'
    })
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <label class="input validator w-full">
      <svg class="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
        <g stroke-linejoin="round" stroke-linecap="round" stroke-width="2.5" fill="none" stroke="currentColor">
          <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
          <circle cx="12" cy="7" r="4"></circle>
        </g>
      </svg>
      <input v-model="form.username" type="text" required placeholder="Username" pattern="[A-Za-z][A-Za-z0-9\-]*"
        minlength="3" maxlength="30" title="Only letters, numbers or dash" :disabled="loading" />
    </label>
    <p class="validator-hint">
      Must be 3 to 30 characters
      <br />containing only letters, numbers or dash
    </p>

    <input v-model="form.password" type="password" class="input validator w-full" required placeholder="Password" minlength="8"
      title="Must be more than 8 characters, including number, lowercase letter, uppercase letter"
      :disabled="loading" />
    <p class="validator-hint">
      Must be more than 8 characters, including
      <br />At least one number
      <br />At least one lowercase letter
      <br />At least one uppercase letter
    </p>

    <button @click.prevent="handleSubmit" class="btn btn-primary w-full" :disabled="loading" >
      <span v-if="!loading">Valider</span>
      <span v-else class="loading loading-spinner"></span>
    </button>
  </form>
</template>