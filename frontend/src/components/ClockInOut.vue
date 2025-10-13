<script setup lang="ts">
import { onMounted } from 'vue';
import { useAuthStore } from '@/store/AuthStore';
import { storeToRefs } from 'pinia';

const authStore = useAuthStore();

const { isClockedIn } = storeToRefs(authStore)

const toggleClock = async () => {

  try {
    //calculate new state

    const newState = !isClockedIn.value;

    //Call the store function
    await authStore.updateClocking(newState)

    //Afficher l'alerte, tant qu'on a pas les toaster (je sais que c'est crado)
    if (newState) {
      alert('Clock In effectué !');
    } else {
      alert('Clock Out effectué !');
    }
  }
  catch (error) {
    alert('Erreur lors du pointage');
  }
};

onMounted(async () => {
  try {
    await authStore.isClocked();
  } catch (error) {
    console.error('Erreur chargement statut:', error);
  }
});


</script>

<template>
  <div class="hero min-h-screen bg-base-200">
    <div class="hero-content text-center">
      <div class="max-w-md">
        <h1 class="text-5xl font-bold mb-8">Gestion du temps</h1>

        <div class="mb-8">
          <!-- Cas d'erreur -->
          <div v-if="isClockedIn === undefined" class="alert alert-error">
            <span>Impossible de récupérer votre statut</span>
          </div>
          <!-- Clocked In -->
          <div v-else-if="isClockedIn" class="alert alert-success">
            <span>Vous êtes actuellement en service</span>
          </div>
          <!-- Clocked Out -->
          <div v-else class="alert alert-info">
            <span>Vous n'êtes pas en service</span>
          </div>
        </div>

        <button @click="toggleClock" :class="isClockedIn ? 'btn-error' : 'btn-success'"
          :disabled="isClockedIn === undefined" class="btn btn-lg w-full text-xl">
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