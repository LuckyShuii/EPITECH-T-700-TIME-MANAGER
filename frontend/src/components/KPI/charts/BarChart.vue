<script setup lang="ts">
import { ref, onMounted, watch, onBeforeUnmount } from 'vue'
import { Chart, ChartConfiguration, registerables } from 'chart.js'

// Enregistrer tous les composants Chart.js nécessaires
Chart.register(...registerables)

interface Props {
  data: ChartConfiguration['data']
  options?: ChartConfiguration['options']
}

const props = withDefaults(defineProps<Props>(), {
  options: () => ({})
})

const canvasRef = ref<HTMLCanvasElement | null>(null)
let chartInstance: Chart | null = null

const createChart = () => {
  if (!canvasRef.value) return

  // Détruire l'ancienne instance si elle existe
  if (chartInstance) {
    chartInstance.destroy()
  }

  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) return

  chartInstance = new Chart(ctx, {
    type: 'bar',
    data: props.data,
    options: {
      responsive: true,
      maintainAspectRatio: true,
      ...props.options
    }
  })
}

onMounted(() => {
  createChart()
})

// Recréer le chart si les données changent
watch(() => props.data, () => {
  createChart()
}, { deep: true })

// Nettoyer à la destruction du composant
onBeforeUnmount(() => {
  if (chartInstance) {
    chartInstance.destroy()
  }
})
</script>

<template>
  <div class="chart-container">
    <canvas ref="canvasRef"></canvas>
  </div>
</template>

<style scoped>
.chart-container {
  position: relative;
  width: 100%;
  height: 100%;
}
</style>