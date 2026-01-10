<script setup lang="ts">
import { onMounted, computed } from 'vue';
import { useAuthStore } from '@/store/AuthStore';
import { storeToRefs } from 'pinia';
import { useNotificationsStore } from '@/store/NotificationsStore'

const authStore = useAuthStore();
const { isClockedIn, sessionStatus } = storeToRefs(authStore)

// Computed pour savoir si on est en pause (basé sur le store maintenant)
const isPaused = computed(() => sessionStatus.value === 'paused')


const notificationsStore = useNotificationsStore()


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
    await updating(newState)

    if (newState) {
      notificationsStore.addNotification({
            status: 'success',
            title: 'Début de journée',
            description: 'Vous commencez votre journée'
        })
    } else {
      notificationsStore.addNotification({
            status: 'success',
            title: 'Fin de journée',
            description: 'Vous avez fini votre journée'
        })
    }
  }
  catch (error) {
    notificationsStore.addNotification({
            status: 'error',
            title: 'Erreur pointage',
            description: 'Oopsie'
        })
  }
};

const togglePause = async () => {
  try {
    // Si on est en pause, on reprend (is_breaking = false)
    // Si on est actif, on se met en pause (is_breaking = true)
    const isBreaking = !isPaused.value

    await authStore.updateBreaking(isBreaking)

    if (isBreaking) {
      notificationsStore.addNotification({
            status: 'success',
            title: 'Pause démarrée',
            description: 'Vous êtes en pause'
        })
    } else {
      notificationsStore.addNotification({
            status: 'success',
            title: 'Pause terminée',
            description: 'Votre pause est finie'
        })
    }
  } catch (error) {
    notificationsStore.addNotification({
            status: 'error',
            title: 'Echec de la mise en pause',
            description: 'Oopsie'
        })
  }
}


const updating = async (clockingStatus: boolean) => {
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
        class="brutal-btn brutal-btn-success flex-1">
        Clock In
      </button>

      <!-- État CLOCKED_IN : Clock Out -->
      <button v-else-if="currentState === 'CLOCKED_IN'" key="clock-out" @click="toggleClock"
        class="brutal-btn brutal-btn-error flex-1">
        Clock Out
      </button>

      <!-- État PAUSED : Reprendre -->
      <button v-else-if="currentState === 'PAUSED'" key="resume" @click="togglePause"
        class="brutal-btn brutal-btn-warning flex-1">
        Reprendre
      </button>
    </Transition>

    <!-- Bouton droit : N'EXISTE QUE si clocké -->
    <Transition name="morph" mode="out-in">
      <!-- État CLOCKED_IN : Pause active -->
      <button v-if="currentState === 'CLOCKED_IN'" key="pause" @click="togglePause"
        class="brutal-btn brutal-btn-warning flex-1">
        Pause
      </button>
    </Transition>
  </div>
</template>