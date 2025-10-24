<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { Calendar } from 'v-calendar'
import 'v-calendar/style.css'
import type { ClockHistoryEntry } from '@/types/ClockHistoryEntry'
import API from '@/services/API'

const formatDateForAPI = (date: Date): string => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}
const formatDuration = (minutes: number): string => {
  if (minutes < 60) {
    return `${minutes}min`
  }
  
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  
  if (remainingMinutes === 0) {
    return `${hours}h`
  }
  
  return `${hours}h${String(remainingMinutes).padStart(2, '0')}`
}
const minutesToHours = (minutes: number): number => {
  return Math.round((minutes / 60) * 100) / 100 // Arrondi à 2 décimales
}

// Données de l'historique
const historyData = ref<ClockHistoryEntry[]>([])
const isLoading = ref(false)

// Date affichée dans le calendrier
const displayedMonth = ref(new Date())

// Référence au calendrier
const calendar = ref()

// Fonction pour récupérer le premier et dernier jour du mois
const getMonthRange = (date: Date) => {
  const year = date.getFullYear()
  const month = date.getMonth()
  
  const startDate = new Date(year, month, 1, 0, 0, 0)
  const endDate = new Date(year, month + 1, 0, 23, 59, 59)
  
  const now = new Date()
  const finalEndDate = endDate > now ? now : endDate
  
  return {
    start_date: formatDateForAPI(startDate),
    end_date: formatDateForAPI(finalEndDate)
  }
}

// Fonction pour charger l'historique
const fetchHistory = async (date: Date) => {
  isLoading.value = true
  try {
    const { start_date, end_date } = getMonthRange(date)
    
    const response = await API.WorkSession.getWorkSessionHistory({
      start_date,
      end_date
    })
    
    historyData.value = response.data
    console.log('📅 Historique chargé:', historyData.value)
  } catch (error) {
    console.error('Erreur chargement historique:', error)
    historyData.value = []
  } finally {
    isLoading.value = false
  }
}

// Charger l'historique au montage
onMounted(() => {
  fetchHistory(displayedMonth.value)
})

// Fonction bouton aujourd'hui
const moveToday = () => {
  calendar.value?.move(new Date())
  displayedMonth.value = new Date()
}

// Détecter le changement de mois
// v-calendar émet un événement 'update:pages' quand on change de mois
const onPageChange = (pages: any) => {
  if (pages && pages.length > 0) {
    const newMonth = pages[0].month - 1 // v-calendar month est 1-based
    const newYear = pages[0].year
    const newDate = new Date(newYear, newMonth, 1)
    
    // Ne recharger que si le mois a vraiment changé
    if (newDate.getMonth() !== displayedMonth.value.getMonth() || 
        newDate.getFullYear() !== displayedMonth.value.getFullYear()) {
      displayedMonth.value = newDate
      fetchHistory(newDate)
    }
  }
}

// Attribute pour aujourd'hui (point bleu)
const todayAttribute = {
  dates: new Date(),
  highlight: {
    color: 'blue',
    class: 'opacity-75'
  }
}

// Transformation des données en attributes
const calendarAttributes = computed(() => {
  const workDayAttributes = historyData.value.map(entry => {
    const clockInDate = new Date(entry.clock_in)
    const clockOutDate = entry.clock_out ? new Date(entry.clock_out) : null
    
    const heureArrivee = `${clockInDate.getHours()}:${String(clockInDate.getMinutes()).padStart(2, '0')}`
    const heureDepart = clockOutDate 
      ? `${clockOutDate.getHours()}:${String(clockOutDate.getMinutes()).padStart(2, '0')}` 
      : "En cours"
    
    const totalFormatted = formatDuration(entry.duration_minutes) // ICI
    
    return {
      dates: clockInDate,
      dot: {
        color: entry.status === 'completed' ? 'green' : 'orange',
        fillMode: 'light'
      },
      popover: {
        label: `Arrivée: ${heureArrivee}\nDépart: ${heureDepart}\nTotal: ${totalFormatted}` // ICI
      }
    }
  })
  
  return [...workDayAttributes, todayAttribute]
})



</script>

<template>
  <div>
    <div v-if="isLoading" class="text-center p-4">
      Chargement...
    </div>
    <Calendar
      ref="calendar"
      :attributes="calendarAttributes as any"
      :is-dark="false"
      expanded
      locale="fr"
      @update:pages="onPageChange"
    >
      <template #footer>
        <div class="w-full px-4 pb-3">
          <button
            class="bg-indigo-600 hover:bg-indigo-700 text-white font-bold w-full px-3 py-1 rounded-md"
            @click="moveToday"
          >
            Aujourd'hui
          </button>
        </div>
      </template>
    </Calendar>
  </div>
</template>

<style scoped>
/* Griser les week-ends (samedi et dimanche) */
:deep(.vc-day.weekday-7 .vc-day-content),
:deep(.vc-day.weekday-1 .vc-day-content) {
  color: #9ca3af !important;
  opacity: 0.5;
}

/* Griser les jours hors du mois actuel */
:deep(.vc-day.is-not-in-month .vc-day-content) {
  opacity: 0.3;
}
</style>