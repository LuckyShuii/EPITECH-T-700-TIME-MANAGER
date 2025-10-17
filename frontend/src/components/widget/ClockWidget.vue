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
<div class="bg-gradient-dark rounded-2xl p-8 shadow-2xl">
  <!-- Temps travaillé -->
  <div class="text-center">
    <div class="text-gray-500 text-sm uppercase tracking-[0.15em] mb-2">
      Temps travaillé
    </div>
    <div class="font-mono text-[2rem] lg:text-[3rem] font-bold text-green-500 drop-shadow-[0_0_150px_rgba(0,255,0,0.3)] tracking-[0.1em]">
      {{ workedTime }}
    </div>
    <div class="text-gray-400 text-sm mt-4">
      Début : {{ startTime }}
    </div>
    <!-- Indicateur de pause (optionnel) -->
    <div v-if="status === 'paused'" class="mt-2">
      <span class="text-yellow-500 text-sm uppercase tracking-wider">⏸️ En pause</span>
    </div>
  </div>

  <!-- Séparateur -->
  <div class="h-px bg-gradient-to-r from-transparent via-gray-600 to-transparent my-6"></div>

  <!-- Horloge actuelle -->
  <div class="text-center mb-8">
    <div class="font-mono text-2xl lg:text-[2.5rem] font-bold text-red-500 drop-shadow-[0_0_100px_rgba(255,0,0,0.5)] tracking-[0.1em] leading-none">
      {{ currentTime }}
    </div>
    <div class="text-gray-500 text-sm mt-2 uppercase tracking-[0.15em]">
      Heure actuelle
    </div>
  </div>
</div>
</template>