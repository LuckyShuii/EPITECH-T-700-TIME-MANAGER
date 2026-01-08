<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue';

interface ClockWidgetProps {
  clockInTime?: string | null
  status?: 'active' | 'paused' | 'completed' | 'no_active_session'
}

const props = withDefaults(defineProps<ClockWidgetProps>(), {
  clockInTime: null,
  status: 'no_active_session'
})

const currentTime = ref('')

const updateTime = () => {
  const now = new Date()
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  const seconds = String(now.getSeconds()).padStart(2, '0')
  currentTime.value = `${hours}:${minutes}:${seconds}`
}

const formatDuration = (milliseconds: number): string => {
  const totalMinutes = Math.floor(milliseconds / 1000 / 60)
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}`
}

const workedTime = computed(() => {
  // Force la réactivité avec currentTime
  currentTime.value

  
  // Vérifie si on est dans un état "clocké" (active ou paused)
  if (!props.clockInTime || (props.status !== 'active' && props.status !== 'paused')) {
    return '--:--'
  }
  
  const clockIn = new Date(props.clockInTime)
  const now = new Date()
  const duration = now.getTime() - clockIn.getTime()
  

  
  return formatDuration(duration)
})

const startTime = computed(() => {
  if (!props.clockInTime) {
    return '--:--'
  }
  
  const clockIn = new Date(props.clockInTime)
  const hours = String(clockIn.getHours()).padStart(2, '0')
  const minutes = String(clockIn.getMinutes()).padStart(2, '0')
  return `${hours}:${minutes}`
})

let intervalId: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  updateTime()
  intervalId = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId)
  }
})
</script>

<template>
<div class="brutal-container">
  <!-- Temps travaillé -->
  <div class="text-center">
    <div class="brutal-title">Temps travaillé</div>
    <div class="brutal-value brutal-success">
      {{ workedTime }}
    </div>
    <div class="brutal-text mt-4">
      Début : {{ startTime }}
    </div>
    <!-- Indicateur de pause (optionnel) -->
    <div v-if="status === 'paused'" class="mt-2">
      <span class="brutal-warning text-sm uppercase tracking-wider">En pause</span>
    </div>
  </div>
  <!-- Séparateur -->
  <div class="brutal-divider"></div>
  <!-- Horloge actuelle -->
  <div class="text-center mb-8">
    <div class="brutal-value brutal-error">
      {{ currentTime }}
    </div>
    <div class="brutal-text mt-2">Heure actuelle</div>
  </div>
</div>
</template>