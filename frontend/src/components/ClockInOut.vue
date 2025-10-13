<script setup lang="ts">
import { ref } from 'vue';

const isClockedIn = ref(false);

const toggleClock = () => {
  isClockedIn.value = !isClockedIn.value;
  
  if (isClockedIn.value) {
    alert('Clock In effectué !');
  } else {
    alert('Clock Out effectué !');
  }
  
  // TODO: Appeler l'API backend ici plus tard
};
</script>

<template>
  <div class="hero min-h-screen bg-base-200">
    <div class="hero-content text-center">
      <div class="max-w-md">
        <h1 class="text-5xl font-bold mb-8">Gestion du temps</h1>
        
        <div class="mb-8">
          <div v-if="isClockedIn" class="alert alert-success">
            <span>✅ Vous êtes actuellement en service</span>
          </div>
          <div v-else class="alert alert-info">
            <span>Vous n'êtes pas en service</span>
          </div>
        </div>

        <button 
          @click="toggleClock"
          :class="isClockedIn ? 'btn-error' : 'btn-success'"
          class="btn btn-lg w-full text-xl"
        >
          <span v-if="isClockedIn">Clock Out</span>
          <span v-else>Clock In</span>
        </button>

        <p class="mt-4 text-sm text-base-content/70">
          {{ isClockedIn ? 'Cliquez pour terminer votre service' : 'Cliquez pour commencer votre service' }}
        </p>
      </div>
    </div>
  </div>
</template>