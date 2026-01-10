<script setup lang="ts">
import { reactive, ref } from 'vue';
import type { UserLogin } from "@/types/userType";
import { useAuthStore } from '@/store/AuthStore';
import router from '@/router';
import { useNotificationsStore } from '@/store/NotificationsStore';

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
    notificationsStore.addNotification({
      status: 'success',
      title: 'Connexion réussie',
      description: 'Vous êtes maintenant connecté'
    })
    router.push('/dashboard')
  }
  catch (error) {
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
  <form @submit.prevent="handleSubmit" class="space-y-6">
    <!-- Username Input -->
    <div>
      <label class="block text-xs font-bold uppercase tracking-widest mb-3">
        USERNAME
      </label>
      <input 
        v-model="form.username" 
        type="text" 
        required 
        placeholder="Username" 
        pattern="[A-Za-z][A-Za-z0-9\-]*"
        minlength="3" 
        maxlength="30" 
        title="Only letters, numbers or dash"
        :disabled="loading"
        class="w-full border-2 border-black p-3 font-bold uppercase text-sm focus:outline-none focus:border-black bg-white"
      />
      
    </div>

    <!-- Password Input -->
    <div>
      <label class="block text-xs font-bold uppercase tracking-widest mb-3">
        PASSWORD
      </label>
      <input 
        v-model="form.password" 
        type="password" 
        required 
        placeholder="Password" 
        minlength="8"
        title="Must be more than 8 characters, including number, lowercase letter, uppercase letter"
        :disabled="loading"
        class="w-full border-2 border-black p-3 font-bold uppercase text-sm focus:outline-none focus:border-black bg-white"
      />
      
    </div>

    <!-- Submit Button -->
    <button 
      @click.prevent="handleSubmit" 
      class="brutal-btn brutal-btn-primary w-full mt-8"
      :disabled="loading"
    >
      <span v-if="!loading">CONNEXION</span>
      <span v-else class="loading loading-spinner"></span>
    </button>
  </form>
</template>

<style scoped>
input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

input:focus {
  box-shadow: none;
}
</style>