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
  return Math.round((minutes / 60) * 100) / 100
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
    console.log('Historique chargé:', historyData.value)
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

// Fonction bouton aujourd'hui - MODIFIÉE pour recharger les données
const moveToday = () => {
  calendar.value?.move(new Date())
  displayedMonth.value = new Date()
  fetchHistory(new Date())
}

// Détecter le changement de mois - MODIFIÉE pour recharger les données à chaque changement
const onPageChange = (pages: any) => {
  if (pages && pages.length > 0) {
    const newMonth = pages[0].month - 1
    const newYear = pages[0].year
    const newDate = new Date(newYear, newMonth, 1)
    
    // Recharger si le mois a changé
    if (newDate.getMonth() !== displayedMonth.value.getMonth() || 
        newDate.getFullYear() !== displayedMonth.value.getFullYear()) {
      displayedMonth.value = newDate
      fetchHistory(newDate)
    }
  }
}

// Attribute pour aujourd'hui - RÉACTIF
const todayAttribute = computed(() => ({
  dates: new Date(),
  highlight: {
    color: 'blue',
    class: 'opacity-75'
  }
}))

// Transformation des données en attributes
const calendarAttributes = computed(() => {
  const workDayAttributes = historyData.value.map(entry => {
    const clockInDate = new Date(entry.clock_in)
    const clockOutDate = entry.clock_out ? new Date(entry.clock_out) : null
    
    const heureArrivee = `${clockInDate.getHours()}:${String(clockInDate.getMinutes()).padStart(2, '0')}`
    const heureDepart = clockOutDate 
      ? `${clockOutDate.getHours()}:${String(clockOutDate.getMinutes()).padStart(2, '0')}` 
      : "En cours"
    
    const totalFormatted = formatDuration(entry.duration_minutes)
    
    // Déterminer la couleur et le symbole du carré
    const isCompleted = entry.status === 'completed'
    const dotColor = isCompleted ? 'green' : 'orange'
    const colorSymbol = isCompleted ? '🟩' : '🟧'
    
    return {
      dates: clockInDate,
      dot: {
        color: dotColor,
        fillMode: 'light'
      },
      popover: {
        label: `${colorSymbol} Arrivée: ${heureArrivee}\n${colorSymbol} Départ: ${heureDepart}\n${colorSymbol} Total: ${totalFormatted}`
      }
    }
  })
  
  return [...workDayAttributes, todayAttribute.value]
})

</script>

<template>
  <div class="calendar-wrapper">
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
        <div class="w-full">
          <button @click="moveToday">Aujourd'hui</button>
        </div>
      </template>
    </Calendar>
  </div>
</template>

<style scoped>
/* === CONTENEUR PRINCIPAL === */
.calendar-wrapper {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  background-color: hsl(var(--b1));
}

/* En dark mode, fond blanc pour le conteneur */
[data-theme="dark"] .calendar-wrapper {
  background-color: white;
}

:deep(.vc-container) {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  background-color: hsl(var(--b1));
  border: 2px solid black;
  border-radius: 0;
  padding: 0;
  margin: 0;
}

/* En dark mode, le conteneur du calendrier aussi blanc */
[data-theme="dark"] :deep(.vc-container) {
  background-color: white;
  color: black;
  border-radius: 0;
}

/* === HEADER DU CALENDRIER === */
:deep(.vc-header) {
  border-bottom: 2px solid hsl(var(--bc));
  padding: 1rem;
  background-color: hsl(var(--b2));
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 1rem;
}

/* En dark mode, header blanc */
[data-theme="dark"] :deep(.vc-header) {
  background-color: white;
  border-bottom-color: black;
}

/* Flèches de navigation */
:deep(.vc-nav-arrow) {
  background: none;
  border: 2px solid hsl(var(--bc));
  color: hsl(var(--bc));
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-weight: 700;
  font-size: 1rem;
  transition: none;
  padding: 0;
  flex-shrink: 0;
}

/* En dark mode, flèches noires */
[data-theme="dark"] :deep(.vc-nav-arrow) {
  border-color: black;
  color: black;
}

:deep(.vc-nav-arrow:hover) {
  background-color: hsl(var(--bc));
  color: hsl(var(--b1));
}

/* En dark mode, hover flèches */
[data-theme="dark"] :deep(.vc-nav-arrow:hover) {
  background-color: black;
  color: white;
}

:deep(.vc-nav-arrow:active) {
  transform: scale(0.95);
}

/* Titre du mois/année - CENTRÉ avec CSS Grid */
:deep(.vc-title) {
  font-size: 1.125rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: hsl(var(--bc));
  font-family: monospace;
  text-align: center;
  grid-column: 1 / -1;
  order: 2;
}

/* En dark mode, titre noir */
[data-theme="dark"] :deep(.vc-title) {
  color: black;
}

/* === GRILLE DES JOURS DE LA SEMAINE === */
:deep(.vc-weekday) {
  background-color: hsl(var(--b2));
  border-bottom: 2px solid hsl(var(--bc));
  padding: 0.75rem;
  text-align: center;
  font-weight: 700;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: hsl(var(--bc));
  font-family: monospace;
  margin-top: 0.5rem;
}

/* En dark mode, weekday blanc */
[data-theme="dark"] :deep(.vc-weekday) {
  background-color: white;
  border-bottom-color: black;
  color: black;
}

/* === GRILLE DES JOURS === */
:deep(.vc-weeks) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0;
}

:deep(.vc-week) {
  display: flex;
  flex: 1;
  border-bottom: 1px solid hsl(var(--bc) / 0.3);
}

:deep(.vc-week:last-child) {
  border-bottom: none;
  display: none;
}

/* === CHAQUE JOUR === */
:deep(.vc-day) {
  flex: 1;
  border-right: 1px solid hsl(var(--bc) / 0.3);
  position: relative;
  display: flex;
  flex-direction: column;
  padding: 0.5rem;
  background-color: hsl(var(--b1));
}

/* En dark mode, jours blancs */
[data-theme="dark"] :deep(.vc-day) {
  background-color: white;
  border-right-color: hsl(0 0% 80%);
}

:deep(.vc-day:last-child) {
  border-right: none;
}

/* Jour hors du mois actuel */
:deep(.vc-day.is-not-in-month) {
  background-color: hsl(var(--b2));
  opacity: 0.5;
}

/* En dark mode, jour hors mois */
[data-theme="dark"] :deep(.vc-day.is-not-in-month) {
  background-color: hsl(0 0% 90%);
  opacity: 0.5;
  color: black;
}

/* Jour actuel (aujourd'hui) */
:deep(.vc-day.is-today) {
  background-color: hsl(var(--b2));
  border: 2px solid hsl(var(--bc));
}

/* En dark mode, aujourd'hui */
[data-theme="dark"] :deep(.vc-day.is-today) {
  background-color: hsl(0 0% 90%);
  border-color: black;
}

/* === CONTENU DU JOUR === */
:deep(.vc-day-content) {
  font-size: 1rem;
  font-weight: 700;
  color: hsl(var(--bc));
  font-family: monospace;
  text-align: left;
  padding: 0.25rem;
}

/* En dark mode, texte noir */
[data-theme="dark"] :deep(.vc-day-content) {
  color: black;
}

:deep(.vc-day.is-not-in-month .vc-day-content) {
  color: hsl(var(--bc) / 0.5);
}

/* En dark mode, jours hors mois */
[data-theme="dark"] :deep(.vc-day.is-not-in-month .vc-day-content) {
  color: hsl(0 0% 50%);
}

/* === REMPLACEMENT DES POINTS PAR DES CARRÉS === */
:deep(.vc-dots) {
  display: flex !important;
  flex-wrap: wrap;
  gap: 3px;
  padding: 0.5rem 0 0 0;
  justify-content: flex-start;
}

:deep(.vc-dot) {
  width: 8px !important;
  height: 8px !important;
  min-width: 8px !important;
  min-height: 8px !important;
  border-radius: 0 !important;
  border: 1px solid black !important;
  background-color: white !important;
  display: inline-block !important;
  flex-shrink: 0 !important;
  visibility: visible !important;
  opacity: 1 !important;
  margin: 0 !important;
  padding: 0 !important;
}

/* === POPOVER (bulle au survol) === */
:deep(.vc-popover-content) {
  background-color: white !important;
  border: 2px solid black !important;
  color: black !important;
  padding: 0.75rem !important;
  font-family: monospace;
  font-size: 0.875rem !important;
  white-space: pre-wrap;
  line-height: 1.6;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  opacity: 1 !important;
  backdrop-filter: none !important;
}

/* En dark mode, popover pareil (déjà blanc) */
[data-theme="dark"] :deep(.vc-popover-content) {
  background-color: white !important;
  border-color: black !important;
  color: black !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
}

:deep(.vc-popover-caret) {
  border-top-color: black !important;
}

/* En dark mode, caret noir */
[data-theme="dark"] :deep(.vc-popover-caret) {
  border-top-color: black !important;
}

/* === FOOTER (Bouton "Aujourd'hui") === */
:deep(.vc-footer) {
  border-top: 2px solid hsl(var(--bc));
  padding: 0;
  background-color: hsl(var(--b2));
}

/* En dark mode, footer blanc */
[data-theme="dark"] :deep(.vc-footer) {
  background-color: white;
  border-top-color: black;
}

/* Bouton "Aujourd'hui" */
button {
  width: 100%;
  border: none;
  background-color: hsl(var(--b1));
  border-top: 2px solid hsl(var(--bc));
  color: hsl(var(--bc));
  padding: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  font-size: 0.875rem;
  letter-spacing: 0.05em;
  cursor: pointer;
  font-family: monospace;
  transition: none;
}

/* En dark mode, bouton blanc */
[data-theme="dark"] button {
  background-color: white;
  color: black;
  border-top-color: black;
}

button:hover {
  background-color: hsl(var(--bc));
  color: hsl(var(--b1));
}

/* En dark mode, pas de hover ou très subtil */
[data-theme="dark"] button:hover {
  background-color: white;
  color: black;
}

button:active {
  transform: scale(0.98);
}

/* === GRISER LES WEEK-ENDS === */
:deep(.vc-day.weekday-7 .vc-day-content),
:deep(.vc-day.weekday-1 .vc-day-content) {
  color: hsl(var(--bc) / 0.6);
  opacity: 0.6;
}
</style>