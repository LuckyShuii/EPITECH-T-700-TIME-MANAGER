<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => router.push('/')
const goBack = () => router.back()
</script>

<template>
  <div class="min-h-screen bg-white text-black flex items-center justify-center p-8">
    <div class="text-center border-4 border-black p-12 max-w-2xl">
      <!-- Titre 404 -->
      <div class="mb-12">
        <h1 class="text-9xl font-black tracking-wider mb-4">
          404
        </h1>
        <div class="border-2 border-black p-4">
          <p class="font-black uppercase tracking-widest text-lg">
            PAGE INTROUVABLE
          </p>
        </div>
      </div>

      <!-- Message -->
      <div class="mb-12">
        <p class="text-base font-bold leading-relaxed mb-4">
          Cette page s'est perdue dans l'espace-temps.
        </p>
        <p class="text-sm font-bold">
          Elle est peut-être en pause café... ou elle n'a jamais existé.
        </p>
      </div>

      <!-- Actions -->
      <div class="flex gap-4 justify-center flex-wrap">
        <button 
          @click="goBack"
          class="brutal-btn brutal-btn-primary"
        >
          Retour
        </button>
        <button 
          @click="goHome" 
          class="brutal-btn brutal-btn-success"
        >
          Accueil
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