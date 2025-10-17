<script setup lang="ts">
import { onMounted, ref, computed } from 'vue';
import { useAuthStore } from '@/store/AuthStore';
import { storeToRefs } from 'pinia';

const authStore = useAuthStore();
const { isClockedIn } = storeToRefs(authStore)

// État local pour la pause
const isPaused = ref(false)

// Computed pour savoir quel état afficher
const currentState = computed(() => {
  if (isClockedIn.value === undefined) return 'ERROR'
  if (!isClockedIn.value) return 'NOT_CLOCKED'
  if (isPaused.value) return 'PAUSED'
  return 'CLOCKED_IN'
})

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

// Action 2 : Pause/Reprendre (À ÉCRIRE)
const togglePause = () => {
  isPaused.value = !isPaused.value
  // TODO : Appeler une API pour enregistrer la pause si besoin
  // Pour l'instant juste le toggle local
}

// Action 3 : Stop (Clock Out depuis la pause) (À ÉCRIRE)
const stopShift = async () => {
  isPaused.value = false
  // Appeler toggleClock ou directement updateClocking(false)
  await authStore.updateClocking(false)
  alert('Service terminé')
}

onMounted(async () => {
  try {
    await authStore.isClocked();
  } catch (error) {
    console.error('Erreur chargement statut:', error);
  }
});
</script>


<template>
  <div class="flex gap-2">
    <!-- Bouton gauche : Change selon l'état -->
    <Transition name="morph" mode="out-in">
      <!-- État NOT_CLOCKED : Clock In -->
      <button 
        v-if="currentState === 'NOT_CLOCKED'"
        key="clock-in"
        @click="toggleClock"
        class="btn btn-success btn-lg w-full"
      >
        Clock In
      </button>
      
      <!-- État CLOCKED_IN : Clock Out -->
      <button 
        v-else-if="currentState === 'CLOCKED_IN'"
        key="clock-out"
        @click="toggleClock"
        class="btn btn-error btn-lg flex-1"
      >
        Clock Out
      </button>
      
      <!-- État PAUSED : Reprendre -->
      <button 
        v-else-if="currentState === 'PAUSED'"
        key="resume"
        @click="togglePause"
        class="btn btn-warning btn-lg flex-1"
      >
        Reprendre
      </button>
    </Transition>

    <!-- Bouton droit : N'EXISTE QUE si clocké -->
    <Transition name="morph" mode="out-in">
      <!-- État CLOCKED_IN : Pause active -->
      <button 
        v-if="currentState === 'CLOCKED_IN'"
        key="pause"
        @click="togglePause"
        class="btn btn-warning btn-lg flex-1"
      >
        ⏸️ Pause
      </button>
      
      <!-- État PAUSED : Stop -->
      <button 
        v-else-if="currentState === 'PAUSED'"
        key="stop"
        @click="stopShift"
        class="btn btn-error btn-lg flex-1"
      >
        ⏹️ Stop
      </button>
    </Transition>
  </div>
</template>


<style scoped>
/* Transitions pour l'effet morphing */
.morph-enter-active,
.morph-leave-active {
  transition: all 0.3s ease;
}

.morph-enter-from {
  opacity: 0;
  transform: scale(0.8);
}

.morph-leave-to {
  opacity: 0;
  transform: scale(0.8);
}

/* Transition smooth de la largeur du conteneur */
.flex {
  transition: all 0.3s ease;
}

/* Animation de respiration pour pause (optionnel) */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}
</style>