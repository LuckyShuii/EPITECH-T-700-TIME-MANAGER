import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './AuthStore'
import API from '@/services/API'
import type {
    WorkingTimeData,
    TeamWorkingTimeData,
    ShiftAverageData,
    PresenceRateData,
    ClockingTimeData,
    IndividualPauseData,
    TeamPauseData,
    WeeklyProgressData,
    ShiftPauseData
} from '@/types/kpi'

export const useKpiStore = defineStore('kpi', () => {
    const authStore = useAuthStore()

    // États de chargement
    const loading = ref<Record<string, boolean>>({})

    // Données KPI
    const workingTimeIndividual = ref<WorkingTimeData | null>(null)
    const workingTimeTeam = ref<TeamWorkingTimeData | null>(null)
    const shiftAverage = ref<ShiftAverageData | null>(null)
    const presenceRate = ref<PresenceRateData | null>(null)
    const clockingTime = ref<ClockingTimeData | null>(null)
    const individualPause = ref<IndividualPauseData | null>(null)
    const teamPause = ref<TeamPauseData | null>(null)
    const weeklyProgress = ref<WeeklyProgressData | null>(null)
    const shiftPause = ref<ShiftPauseData | null>(null)

    // Timestamps pour le cache (Redis côté back + cache local)
    const lastFetch = ref<Record<string, number>>({})
    const CACHE_DURATION = 5 * 60 * 1000 // 5 minutes

    // Helper pour vérifier si le cache est valide
    const isCacheValid = (key: string): boolean => {
        const timestamp = lastFetch.value[key]
        if (!timestamp) return false
        return Date.now() - timestamp < CACHE_DURATION
    }

    // Helper pour setter le loading state
    const setLoading = (key: string, value: boolean) => {
        loading.value[key] = value
    }

    // Permissions : quels KPI l'utilisateur peut-il voir ?
    const availableKpis = computed(() => {
        const kpis: string[] = []

        // KPI pour tous
        kpis.push('workingTimeIndividual', 'weeklyProgress')

        // KPI pour Manager et Admin
        if (authStore.isManager || authStore.isAdmin) {
            kpis.push(
                'workingTimeTeam',
                'shiftAverage',
                'presenceRate',
                'individualPause'
            )
        }

        // KPI uniquement Admin
        if (authStore.isAdmin) {
            kpis.push('clockingTime', 'teamPause', 'shiftPause')
        }

        return kpis
    })

    // Checker si l'utilisateur peut accéder à un KPI
    const canAccessKpi = (kpiName: string): boolean => {
        return availableKpis.value.includes(kpiName)
    }

    // === Fonctions de fetch pour chaque KPI ===

    const fetchWorkingTimeIndividual = async (force = false) => {
        if (!canAccessKpi('workingTimeIndividual')) return
        if (!force && isCacheValid('workingTimeIndividual')) return

        setLoading('workingTimeIndividual', true)
        try {
            // TODO: Remplacer par ton vrai endpoint API
            const response = await API.kpi.getWorkingTimeIndividual()
            workingTimeIndividual.value = response.data
            lastFetch.value['workingTimeIndividual'] = Date.now()
        } catch (error) {
            console.error('Erreur fetch working time individual:', error)
            throw error
        } finally {
            setLoading('workingTimeIndividual', false)
        }
    }

    const fetchWorkingTimeTeam = async (force = false) => {
        if (!canAccessKpi('workingTimeTeam')) return
        if (!force && isCacheValid('workingTimeTeam')) return

        setLoading('workingTimeTeam', true)
        try {
            const response = await API.kpi.getWorkingTimeTeam()
            workingTimeTeam.value = response.data
            lastFetch.value['workingTimeTeam'] = Date.now()
        } catch (error) {
            console.error('Erreur fetch working time team:', error)
            throw error
        } finally {
            setLoading('workingTimeTeam', false)
        }
    }

    const fetchShiftAverage = async (force = false) => {
        if (!canAccessKpi('shiftAverage')) return
        if (!force && isCacheValid('shiftAverage')) return

        setLoading('shiftAverage', true)
        try {
            const response = await API.kpi.getShiftAverage()
            shiftAverage.value = response.data
            lastFetch.value['shiftAverage'] = Date.now()
        } catch (error) {
            console.error('Erreur fetch shift average:', error)
            throw error
        } finally {
            setLoading('shiftAverage', false)
        }
    }

    const fetchPresenceRate = async (force = false) => {
        if (!canAccessKpi('presenceRate')) return
        if (!force && isCacheValid('presenceRate')) return

        setLoading('presenceRate', true)
        try {
            const response = await API.kpi.getPresenceRate()
            presenceRate.value = response.data
            lastFetch.value['presenceRate'] = Date.now()
        } catch (error) {
            console.error('Erreur fetch presence rate:', error)
            throw error
        } finally {
            setLoading('presenceRate', false)
        }
    }

    const fetchWeeklyProgress = async (force = false) => {
  if (!canAccessKpi('weeklyProgress')) return
  if (!force && isCacheValid('weeklyProgress')) return

  setLoading('weeklyProgress', true)
  try {
    const response = await API.kpi.getWeeklyProgress()
    weeklyProgress.value = response.data
    lastFetch.value['weeklyProgress'] = Date.now()
  } catch (error) {
    console.error('Erreur fetch weekly progress:', error)
    throw error
  } finally {
    setLoading('weeklyProgress', false)
  }
}

    // TODO: Ajouter les autres fetch (clockingTime, individualPause, teamPause, shiftPause)

    // Fonction pour tout rafraîchir (utile quand on revient sur le dashboard)
    const refreshAllKpis = async () => {
        const promises = []

        if (canAccessKpi('workingTimeIndividual')) {
            promises.push(fetchWorkingTimeIndividual(true))
        }
        if (canAccessKpi('workingTimeTeam')) {
            promises.push(fetchWorkingTimeTeam(true))
        }
        if (canAccessKpi('shiftAverage')) {
            promises.push(fetchShiftAverage(true))
        }
        if (canAccessKpi('presenceRate')) {
            promises.push(fetchPresenceRate(true))
        }
        if (canAccessKpi('weeklyProgress')) {
  promises.push(fetchWeeklyProgress(true))
}

        await Promise.allSettled(promises)
    }

    return {
        // États
        loading,
        workingTimeIndividual,
        workingTimeTeam,
        shiftAverage,
        presenceRate,
        clockingTime,
        individualPause,
        teamPause,
        weeklyProgress,
        shiftPause,

        // Computed
        availableKpis,
        canAccessKpi,

        // Actions
        fetchWorkingTimeIndividual,
        fetchWorkingTimeTeam,
        fetchShiftAverage,
        fetchPresenceRate,
        fetchWeeklyProgress,
        refreshAllKpis
    }
})