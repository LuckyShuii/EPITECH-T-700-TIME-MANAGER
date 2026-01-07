import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './AuthStore'
import API from '@/services/API'
import type {
  IndividualPauseData,
  WorkingTimeIndividualDisplay,
  WorkingTimeTeamData,
  PresenceRateData
} from '@/types/Kpi'
import { startOfWeek, subWeeks, format, addDays, subDays } from 'date-fns'
import { fr } from 'date-fns/locale'

export const useKpiStore = defineStore('kpi', () => {
  const authStore = useAuthStore()

  // État de la semaine courante
  const currentWeekStart = ref<Date>(subWeeks(startOfWeek(new Date(), { weekStartsOn: 1 }), 1))

  // États de chargement
  const loading = ref<Record<string, boolean>>({})

  // Données KPI
  const individualPause = ref<IndividualPauseData | null>(null)
  const workingTimeIndividual = ref<WorkingTimeIndividualDisplay | null>(null)
  const workingTimeTeam = ref<WorkingTimeTeamData | null>(null)
  const presenceRate = ref<PresenceRateData | null>(null)

  // Équipes du manager
  const managerTeams = ref<WorkingTimeTeamData[]>([])
  const currentTeamIndex = ref<number>(0)

  // Timestamps pour le cache
  const lastFetch = ref<Record<string, number>>({})
  const CACHE_DURATION = 5 * 60 * 1000 // 5 minutes

  // === COMPUTED POUR LES DATES ===

  const weekDateRange = computed(() => {
    const start = currentWeekStart.value
    const end = addDays(start, 4)
    return { start, end }
  })

  const weekStartDate = computed(() => {
    return format(weekDateRange.value.start, 'yyyy-MM-dd')
  })

  const weekEndDate = computed(() => {
    return format(weekDateRange.value.end, 'yyyy-MM-dd')
  })

  const weekDisplayLabel = computed(() => {
    const start = weekDateRange.value.start
    return `Semaine du ${format(start, 'd MMMM yyyy', { locale: fr })}`
  })

  const currentTeam = computed(() => {
    if (managerTeams.value.length === 0) return null
    return managerTeams.value[currentTeamIndex.value] ?? null
  })

  // === PERMISSIONS ===

  const availableKpis = computed(() => {
    const kpis: string[] = []
    kpis.push('individualPause', 'workingTimeIndividual')
    if (authStore.isManager || authStore.isAdmin) {
      kpis.push('workingTimeTeam')
    }
    if (authStore.isAdmin) {
      kpis.push('presenceRate')
    }
    return kpis
  })

  const canAccessKpi = (kpiName: string): boolean => {
    return availableKpis.value.includes(kpiName)
  }

  // === HELPERS ===

  const isCacheValid = (key: string): boolean => {
    const timestamp = lastFetch.value[key]
    if (!timestamp) return false
    return Date.now() - timestamp < CACHE_DURATION
  }

  const setLoading = (key: string, value: boolean) => {
    loading.value[key] = value
  }

  // === NAVIGATION SEMAINES ===

  const goToPreviousWeek = () => {
    currentWeekStart.value = subWeeks(currentWeekStart.value, 1)
  }

  const goToNextWeek = () => {
    currentWeekStart.value = addDays(currentWeekStart.value, 7)
  }

  const goToCurrentWeek = () => {
    currentWeekStart.value = subWeeks(startOfWeek(new Date(), { weekStartsOn: 1 }), 1)
  }

  // === NAVIGATION ÉQUIPES ===

  const goToNextTeam = () => {
    if (managerTeams.value.length === 0) return
    currentTeamIndex.value = (currentTeamIndex.value + 1) % managerTeams.value.length
  }

  const goToPreviousTeam = () => {
    if (managerTeams.value.length === 0) return
    currentTeamIndex.value = (currentTeamIndex.value - 1 + managerTeams.value.length) % managerTeams.value.length
  }

  // === FONCTIONS FETCH ===

  const fetchIndividualPause = async (force = false) => {
    if (!canAccessKpi('individualPause')) return
    if (!force && isCacheValid('individualPause')) return

    if (!authStore.user) throw new Error('Utilisateur non authentifié')

    setLoading('individualPause', true)
    try {
      const response = await API.kpiAPI.getAverageBreakTime(
        authStore.user.user_uuid,
        weekStartDate.value,
        weekEndDate.value
      )
      individualPause.value = response.data
      lastFetch.value['individualPause'] = Date.now()
    } catch (error) {
      console.error('Erreur fetch individual pause:', error)
      throw error
    } finally {
      setLoading('individualPause', false)
    }
  }

  const fetchWorkingTimeIndividual = async (force = false) => {
    if (!canAccessKpi('workingTimeIndividual')) return
    if (!force && isCacheValid('workingTimeIndividual')) return

    if (!authStore.user) throw new Error('Utilisateur non authentifié')

    setLoading('workingTimeIndividual', true)
    try {
      const currentWeekResponse = await API.kpiAPI.getWorkingTimeIndividual(
        authStore.user.user_uuid,
        weekStartDate.value,
        weekEndDate.value
      )

      const previousWeekStart = format(subDays(weekDateRange.value.start, 7), 'yyyy-MM-dd')
      const previousWeekEnd = format(subDays(weekDateRange.value.end, 7), 'yyyy-MM-dd')

      const previousWeekResponse = await API.kpiAPI.getWorkingTimeIndividual(
        authStore.user.user_uuid,
        previousWeekStart,
        previousWeekEnd
      )

      const currentData = currentWeekResponse.data
      const previousData = previousWeekResponse.data

      workingTimeIndividual.value = {
        ...currentData,
        previousTotal: previousData.total_time,
        difference: currentData.total_time - previousData.total_time
      } as WorkingTimeIndividualDisplay

      lastFetch.value['workingTimeIndividual'] = Date.now()
    } catch (error) {
      console.error('Erreur fetch working time individual:', error)
      throw error
    } finally {
      setLoading('workingTimeIndividual', false)
    }
  }

  const fetchManagerTeams = async () => {
    if (!authStore.isManager && !authStore.isAdmin) return

    try {
      const allTeams = await API.teamAPI.getAll()

      if (authStore.isManager) {
        managerTeams.value = allTeams.data.filter((team: any) => {
          return team.team_members.some(
            (member: any) => member.user_uuid === authStore.user?.user_uuid && member.is_manager
          )
        })
      } else if (authStore.isAdmin) {
        managerTeams.value = allTeams.data
      }

      currentTeamIndex.value = 0
    } catch (error) {
      console.error('Erreur fetch manager teams:', error)
    }
  }

  const fetchWorkingTimeTeam = async (force = false) => {
    if (!canAccessKpi('workingTimeTeam')) return
    if (!currentTeam.value) return
    if (!force && isCacheValid('workingTimeTeam')) return

    setLoading('workingTimeTeam', true)
    try {
      const currentWeekResponse = await API.kpiAPI.getWorkingTimeTeam(
        currentTeam.value.team_uuid,
        weekStartDate.value,
        weekEndDate.value
      )

      const previousWeekStart = format(subDays(weekDateRange.value.start, 7), 'yyyy-MM-dd')
      const previousWeekEnd = format(subDays(weekDateRange.value.end, 7), 'yyyy-MM-dd')

      const previousWeekResponse = await API.kpiAPI.getWorkingTimeTeam(
        currentTeam.value.team_uuid,
        previousWeekStart,
        previousWeekEnd
      )

      const currentData = currentWeekResponse.data
      const previousData = previousWeekResponse.data

      workingTimeTeam.value = {
        ...currentData,
        previousTotal: previousData.total_time,
        difference: currentData.total_time - previousData.total_time
      } as any

      lastFetch.value['workingTimeTeam'] = Date.now()
    } catch (error) {
      console.error('Erreur fetch working time team:', error)
      throw error
    } finally {
      setLoading('workingTimeTeam', false)
    }
  }

  const fetchPresenceRate = async (force = false) => {
    if (!canAccessKpi('presenceRate')) return
    if (!force && isCacheValid('presenceRate')) return

    if (!authStore.user) throw new Error('Utilisateur non authentifié')

    setLoading('presenceRate', true)
    try {
      const response = await API.kpiAPI.getPresenceRate(
        authStore.user.user_uuid,
        weekStartDate.value,
        weekEndDate.value
      )
      presenceRate.value = response.data
      lastFetch.value['presenceRate'] = Date.now()
    } catch (error) {
      console.error('Erreur fetch presence rate:', error)
      throw error
    } finally {
      setLoading('presenceRate', false)
    }
  }

  // === RAFRAÎCHISSEMENT GLOBAL ===

  const refreshAllKpis = async () => {
    const promises = []

    if (canAccessKpi('individualPause')) {
      promises.push(fetchIndividualPause(true))
    }
    if (canAccessKpi('workingTimeIndividual')) {
      promises.push(fetchWorkingTimeIndividual(true))
    }
    if (canAccessKpi('workingTimeTeam')) {
      promises.push(fetchWorkingTimeTeam(true))
    }
    if (canAccessKpi('presenceRate')) {
      promises.push(fetchPresenceRate(true))
    }

    await Promise.allSettled(promises)
  }

  const changeWeek = async (direction: 'previous' | 'next' | 'current') => {
    if (direction === 'previous') goToPreviousWeek()
    else if (direction === 'next') goToNextWeek()
    else goToCurrentWeek()

    await refreshAllKpis()
  }

  // === RETURN STORE ===

  return {
    // States
    loading,
    individualPause,
    workingTimeIndividual,
    workingTimeTeam,
    presenceRate,
    managerTeams,
    currentTeam,

    // Computed dates
    weekStartDate,
    weekEndDate,
    weekDisplayLabel,

    // Computed permissions
    availableKpis,
    canAccessKpi,

    // Actions - Semaines
    changeWeek,

    // Actions - Équipes
    goToNextTeam,
    goToPreviousTeam,
    fetchManagerTeams,

    // Actions - KPI
    fetchIndividualPause,
    fetchWorkingTimeIndividual,
    fetchWorkingTimeTeam,
    fetchPresenceRate,
    refreshAllKpis
  }
})