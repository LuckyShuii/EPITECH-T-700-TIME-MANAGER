<script setup lang="ts">
import { computed, ref } from 'vue'
import { Calendar } from 'v-calendar'
import 'v-calendar/style.css'
import type { ClockHistoryEntry } from '@/types/ClockHistoryEntry'

// Mock data
const mockHistoryData: ClockHistoryEntry[] = [
  {
    clockInTime: "2025-10-15T08:15:00Z",
    clockOutTime: "2025-10-15T17:30:00Z",
    totalHours: 8.25,
    status: "completed"
  },
  {
    clockInTime: "2025-10-16T08:15:00Z",
    clockOutTime: "2025-10-16T17:30:00Z",
    totalHours: 9.25,
    status: "completed"
  },
  {
    clockInTime: "2025-10-17T08:15:00Z",
    clockOutTime: "2025-10-17T17:30:00Z",
    totalHours: 9.25,
    status: "completed"
  }, 
  {
    clockInTime: "2025-10-20T08:15:00Z",
    clockOutTime: "2025-10-20T17:30:00Z",
    totalHours: 9.25,
    status: "completed"
  }, 
  {
    clockInTime: "2025-10-21T08:15:00Z",
    clockOutTime: "2025-10-21T17:30:00Z",
    totalHours: 9.25,
    status: "completed"
  },
  {
    clockInTime: "2025-10-22T08:15:00Z",
    clockOutTime: "2025-10-22T17:30:00Z",
    totalHours: 8.25,
    status: "completed"
  },
  {
    clockInTime: "2025-10-23T08:15:00Z",
    clockOutTime: "2025-10-23T17:30:00Z",
    totalHours: 8.25,
    status: "completed"
  },
  {
    clockInTime: "2025-10-24T08:15:00Z",
    clockOutTime: "2025-10-24T17:30:00Z",
    totalHours: 8.25,
    status: "completed"
  }

]


// Référence au calendrier
const calendar = ref()

// Fonction bouton aujourd'hui
const moveToday = () => {
  calendar.value?.move(new Date())
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
  const workDayAttributes = mockHistoryData.map(entry => {
    const dateObj = new Date(entry.clockInTime)
    const clockInDate = new Date(entry.clockInTime)
    const clockOutDate = entry.clockOutTime ? new Date(entry.clockOutTime) : null
    
    const heureArrivee = `${clockInDate.getHours()}:${clockInDate.getMinutes().toString().padStart(2, '0')}`
    const heureDepart = clockOutDate ? `${clockOutDate.getHours()}:${clockOutDate.getMinutes().toString().padStart(2, '0')}` : "En cours"
    
    return {
      dates: dateObj,
      dot: {
        color: 'red',
        fillMode: 'light'
      },
      popover: {
        label: `Arrivée: ${heureArrivee}\nDépart: ${heureDepart}\nTotal: ${entry.totalHours}h`
      }
    }
  })
  
  return [...workDayAttributes, todayAttribute]
})



</script>
<template>
  <div>
    <Calendar 
      ref="calendar"
      :attributes="calendarAttributes as any"
      :is-dark="false"
      expanded
      locale="fr"
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