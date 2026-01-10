<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="space-y-6">
    <!-- Error Container -->
    <div class="border-2 border-red-700 bg-red-50 p-8">
      <p class="text-red-700 font-black uppercase text-lg tracking-wider mb-4">LIEN EXPIRÉ</p>
      <p class="text-red-600 font-bold text-sm">
        Le lien de réinitialisation a expiré. Veuillez demander un nouveau lien.
      </p>
    </div>

    <!-- Action Button -->
    <div class="flex justify-center">
      <button
        @click="goHome"
        class="brutal-btn brutal-btn-primary px-8 py-3"
      >
        Retour à l'accueil
      </button>
    </div>
  </div>
</template>

<style scoped>
button {
  transition: none;
}
</style>