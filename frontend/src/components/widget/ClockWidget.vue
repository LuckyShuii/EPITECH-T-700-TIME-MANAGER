<script setup lang="ts">
import type { ClockData } from '@/types/clockType';
import { ref, onMounted, onUnmounted, computed } from 'vue';

interface ClockWidgetProps extends Partial<ClockData> { }

const props = withDefaults(defineProps<ClockWidgetProps>(), {
  clockInTime: null,
  clockOutTime: null,
  status: 'NOT_CLOCKED'
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
  currentTime.value
  
  if (!props.clockInTime || props.status !== 'CLOCKED_IN') {
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
    </div>
    
    
    <!-- Séparateur -->
    <div class="h-px bg-gradient-to-r from-transparent via-gray-600 to-transparent my-6"></div>
    
    <!-- Temps travaillé -->
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