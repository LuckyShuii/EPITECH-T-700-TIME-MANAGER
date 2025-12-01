<script setup lang="ts">
import type { WorkingTimeData } from '@/types/kpi'

interface Props {
  data: WorkingTimeData
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

// Événement placeholder pour l'interactivité future
const emit = defineEmits<{
  viewDetails: [data: WorkingTimeData]
}>()

const handleClick = () => {
  emit('viewDetails', props.data)
}
</script>
<template>
  <div class="card bg-base-100 shadow-xl cursor-pointer hover:shadow-2xl transition-shadow" @click="handleClick">
    <div class="card-body">
      <h2 class="card-title">Travail hebdomadaire</h2>
      
      <div v-if="loading" class="flex justify-center items-center h-32">
        <span class="loading loading-spinner loading-lg"></span>
      </div>

      <div v-else>
        <!-- Total de la semaine -->
        <div class="stat">
          <div class="stat-title">Total semaine</div>
          <div class="stat-value text-primary">{{ data.totalWeek }}h</div>
        </div>
      </div>
    </div>
  </div>
</template>