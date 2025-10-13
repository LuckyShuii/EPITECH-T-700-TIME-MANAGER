<script setup lang="ts">
import { useAuthStore } from '@/store/AuthStore';
import { storeToRefs } from 'pinia';
import { RouterLink } from 'vue-router';

const authStore = useAuthStore();
const { user, isAuthenticated } = storeToRefs(authStore);
</script>

<template>
  <div class="hero min-h-screen bg-base-200">
    <div class="hero-content text-center">
      <div class="max-w-md">
        
        <!-- Si NON connecté -->
        <div v-if="!isAuthenticated">
          <h1 class="text-5xl font-bold mb-8">Bienvenue</h1>
          <p class="text-xl mb-6">Connectez-vous pour accéder à l'application</p>
          <RouterLink to="/login" class="btn btn-primary btn-lg">
            Se connecter
          </RouterLink>
        </div>

        <!-- Si connecté -->
        <div v-else>
          <h1 class="text-5xl font-bold mb-8">
            Bonjour {{ user?.first_name }} 👋
          </h1>
          <p class="text-xl mb-8">Que souhaitez-vous faire aujourd'hui ?</p>
          
          <div class="flex flex-col gap-4">
            <RouterLink to="/clock" class="btn btn-primary btn-lg">
              ⏱️ Pointer
            </RouterLink>
            <RouterLink to="/profile" class="btn btn-outline btn-lg">
              👤 Mon profil
            </RouterLink>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>