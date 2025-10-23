<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => router.push('/')
const goBack = () => router.back()
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-primary p-8">
    <div class="text-center bg-white rounded-3xl p-12 shadow-card-hover max-w-2xl">
      <!-- Illustration -->
      <div class="mb-8">
        <span class="text-8xl animate-float inline-block">🔍</span>
        <h1 class="text-9xl font-black bg-gradient-primary bg-clip-text text-transparent mt-4">
          404
        </h1>
      </div>
      
      <!-- Message -->
      <div class="mb-10">
        <h2 class="text-3xl font-bold text-gray-800 mb-4">Page introuvable</h2>
        <p class="text-lg text-gray-600 leading-relaxed">
          Oups ! On dirait que cette page s'est perdue dans l'espace-temps.
          <br />
          Elle est peut-être en pause café... ou elle n'a jamais existé.
        </p>
      </div>
      
      <!-- Actions -->
      <div class="flex gap-4 justify-center flex-wrap mb-8">
        <button 
          @click="goHome" 
          class="btn bg-gradient-primary text-white border-none hover:shadow-glow hover:-translate-y-1 transition-all px-8"
        >
          <span class="text-xl">🏠</span>
          Retour à l'accueil
        </button>
      </div>
    </div>
  </div>
</template>