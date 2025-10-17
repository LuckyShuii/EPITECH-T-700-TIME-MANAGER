<script setup lang="ts">
import { onMounted, computed } from 'vue';
import { useAuthStore } from '@/store/AuthStore';
import { storeToRefs } from 'pinia';

const authStore = useAuthStore();
const { isClockedIn, sessionStatus } = storeToRefs(authStore)

// Computed pour savoir si on est en pause (basé sur le store maintenant)
const isPaused = computed(() => sessionStatus.value === 'paused')

// Computed pour savoir quel état afficher
const currentState = computed(() => {
  if (sessionStatus.value === 'no_active_session') return 'NOT_CLOCKED'
  if (sessionStatus.value === 'paused') return 'PAUSED'
  if (sessionStatus.value === 'active') return 'CLOCKED_IN'
  return 'ERROR'
})

const toggleClock = async () => {
  try {
    const newState = !isClockedIn.value;
    await pipiCaca(newState)

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

const togglePause = async () => {
  try {
    // Si on est en pause, on reprend (is_breaking = false)
    // Si on est actif, on se met en pause (is_breaking = true)
    const isBreaking = !isPaused.value

    await authStore.updateBreaking(isBreaking)

    if (isBreaking) {
      alert('Pause démarrée')
    } else {
      alert('Pause terminée')
    }
  } catch (error) {
    alert('Erreur lors de la pause')
  }
}

const stopShift = async () => {
  try {
    await pipiCaca(false)
    alert('Service terminé')
  } catch (error) {
    alert('Erreur lors de l\'arrêt')
  }
}

const pipiCaca = async (clockingStatus: boolean) => {
  await authStore.updateClocking(clockingStatus)
  await authStore.fetchWorkSessionStatus()
}

</script>

<template>
  <div class="flex gap-2">
    <!-- Bouton gauche : Change selon l'état -->
    <Transition name="morph" mode="out-in">
      <!-- État NOT_CLOCKED : Clock In -->
      <button v-if="currentState === 'NOT_CLOCKED'" key="clock-in" @click="toggleClock"
        class="btn btn-success btn-lg w-full">
        Clock In
      </button>

      <!-- État CLOCKED_IN : Clock Out -->
      <button v-else-if="currentState === 'CLOCKED_IN'" key="clock-out" @click="toggleClock"
        class="btn btn-error btn-lg flex-1">
        Clock Out
      </button>

      <!-- État PAUSED : Reprendre -->
      <button v-else-if="currentState === 'PAUSED'" key="resume" @click="togglePause"
        class="btn btn-warning btn-lg flex-1">
        Reprendre
      </button>
    </Transition>

    <!-- Bouton droit : N'EXISTE QUE si clocké -->
    <Transition name="morph" mode="out-in">
      <!-- État CLOCKED_IN : Pause active -->
      <button v-if="currentState === 'CLOCKED_IN'" key="pause" @click="togglePause"
        class="btn btn-warning btn-lg flex-1">
        ⏸️ Pause
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

  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.7;
  }
}
</style>