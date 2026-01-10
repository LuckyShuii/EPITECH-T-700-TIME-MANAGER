import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './AuthStore'
import API from '@/services/API'
import type {
  IndividualPauseData,
  WorkingTimeIndividualData,
  WorkingTimeIndividualDisplay,
  WorkingTimeTeamData,
  PresenceRateData
} from '@/types/Kpi'
import { startOfWeek, subWeeks, format, addDays, subDays } from 'date-fns'
import { fr } from 'date-fns/locale'

export const useKpiStore = defineStore('kpi', () => {
  const authStore = useAuthStore()

  // === ÉTAT ===

  // Semaine courante (S-1 par défaut)
  const currentWeekStart = ref<Date>(subWeeks(startOfWeek(new Date(), { weekStartsOn: 1 }), 1))

  // États de chargement
  const loading = ref<Record<string, boolean>>({})

  // Données KPI
  const individualPause = ref<IndividualPauseData | null>(null)
  const workingTimeIndividual = ref<WorkingTimeIndividualDisplay | null>(null)
  const workingTimeTeam = ref<WorkingTimeTeamData | null>(null)
  const presenceRate = ref<PresenceRateData | null>(null)
  const averageTimePerShift = ref<any>(null)

  // Équipes du manager
  const managerTeams = ref<any[]>([])
  const currentTeamIndex = ref<number>(0)

  // Cache
  const lastFetch = ref<Record<string, number>>({})
  const CACHE_DURATION = 5 * 60 * 1000 // 5 minutes

  // === COMPUTED ===

  const weekDateRange = computed(() => {
    const start = currentWeekStart.value
    const end = addDays(start, 4) // Lundi à Vendredi
    return { start, end }
  })

  const weekStartDate = computed(() => format(weekDateRange.value.start, 'yyyy-MM-dd'))
  const weekEndDate = computed(() => format(weekDateRange.value.end, 'yyyy-MM-dd'))

  const weekDisplayLabel = computed(() => {
    const start = weekDateRange.value.start
    return `Semaine du ${format(start, 'd MMMM yyyy', { locale: fr })}`
  })

  const currentTeam = computed(() => {
    if (managerTeams.value.length === 0) return null
    const team = managerTeams.value[currentTeamIndex.value] ?? null

    console.log('🎯 currentTeam computed:', team)

    return team
  })

  const availableKpis = computed(() => {
    const kpis: string[] = []
    kpis.push('individualPause', 'workingTimeIndividual')
    if (authStore.isManager || authStore.isAdmin) {
      kpis.push('workingTimeTeam', 'averageTimePerShift')
    }
    if (authStore.isAdmin) {
      kpis.push('presenceRate')
    }
    return kpis
  })

  // === HELPERS ===

  const canAccessKpi = (kpiName: string): boolean => {
    return availableKpis.value.includes(kpiName)
  }

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

  const changeWeek = async (direction: 'previous' | 'next' | 'current') => {
    if (direction === 'previous') goToPreviousWeek()
    else if (direction === 'next') goToNextWeek()
    else goToCurrentWeek()
    await refreshAllKpis()
  }

  // === NAVIGATION ÉQUIPES ===

  const goToNextTeam = async () => {
  if (managerTeams.value.length === 0) return
  currentTeamIndex.value = (currentTeamIndex.value + 1) % managerTeams.value.length
  await fetchWorkingTimeTeam(true)
  await fetchPresenceRate(true)
}

const goToPreviousTeam = async () => {
  if (managerTeams.value.length === 0) return
  currentTeamIndex.value = (currentTeamIndex.value - 1 + managerTeams.value.length) % managerTeams.value.length
  await fetchWorkingTimeTeam(true)
  await fetchPresenceRate(true)
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
      // Semaine courante
      const currentWeekResponse = await API.kpiAPI.getWorkingTimeIndividual(
        authStore.user.user_uuid,
        weekStartDate.value,
        weekEndDate.value
      )

      // Semaine précédente pour comparaison
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
      }

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
        managerTeams.value = allTeams.data.filter((team: any) =>
          team.team_members.some(
            (member: any) => member.user_uuid === authStore.user?.user_uuid && member.is_manager
          )
        )
      } else if (authStore.isAdmin) {
        managerTeams.value = allTeams.data
      }

      console.log('🔍 Teams chargées:', managerTeams.value)
      console.log('🔍 Première team:', managerTeams.value[0])

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
      const teamUuid = currentTeam.value.uuid

      console.log('🔍 Fetching KPI pour team:', teamUuid)

      // Semaine courante
      const currentWeekResponse = await API.kpiAPI.getWorkingTimeTeam(
        teamUuid,
        weekStartDate.value,
        weekEndDate.value
      )

      console.log('📦 Données reçues current week:', currentWeekResponse.data)

      // Semaine précédente
      const previousWeekStart = format(subDays(weekDateRange.value.start, 7), 'yyyy-MM-dd')
      const previousWeekEnd = format(subDays(weekDateRange.value.end, 7), 'yyyy-MM-dd')

      const previousWeekResponse = await API.kpiAPI.getWorkingTimeTeam(
        teamUuid,
        previousWeekStart,
        previousWeekEnd
      )

      console.log('📦 Données reçues previous week:', previousWeekResponse.data)

      const currentData = currentWeekResponse.data
      const previousData = previousWeekResponse.data

      workingTimeTeam.value = {
        ...currentData,
        previousTotal: previousData.total_time,
        difference: currentData.total_time - previousData.total_time
      } as any

      console.log('✅ workingTimeTeam.value final:', workingTimeTeam.value)

      lastFetch.value['workingTimeTeam'] = Date.now()
    } catch (error) {
      console.error('Erreur fetch working time team:', error)
      throw error
    } finally {
      setLoading('workingTimeTeam', false)
    }
  }

  const fetchPresenceRate = async (force = false) => {
    if (!canAccessKpi('presenceRate')) {
      console.log('❌ Pas accès à presenceRate')
      return
    }
    if (!force && isCacheValid('presenceRate')) return
    if (!currentTeam.value) return
    
    console.log('🔍 CurrentTeam pour presenceRate:', currentTeam.value.uuid)  // AJOUTE ÇA


    setLoading('presenceRate', true)
    try {
      const teamMembers = currentTeam.value.team_members || []

      if (teamMembers.length === 0) {
        presenceRate.value = []
        lastFetch.value['presenceRate'] = Date.now()
        return
      }

      console.log('🔍 Fetching presence rate pour', teamMembers.length, 'membres')

      const promises = teamMembers.map((member: any) =>
        API.kpiAPI.getPresenceRate(
          member.user_uuid,
          weekStartDate.value,
          weekEndDate.value
        )
      )

      const responses = await Promise.allSettled(promises)

      const results: PresenceRateData = []

      responses.forEach((response) => {
        if (response.status === 'fulfilled' && response.value.data) {
          results.push(response.value.data)
        }
      })

      console.log('📦 Données presence rate reçues:', results)

      presenceRate.value = results
      lastFetch.value['presenceRate'] = Date.now()
    } catch (error) {
      console.error('Erreur fetch presence rate:', error)
      throw error
    } finally {
      setLoading('presenceRate', false)
    }
  }

  const fetchAverageTimePerShift = async (userUuid: string, force = false) => {
    if (!authStore.user) throw new Error('Utilisateur non authentifié')

    setLoading('averageTimePerShift', true)
    try {
      const response = await API.kpiAPI.getAverageTimePerShift(
        userUuid,
        weekStartDate.value,
        weekEndDate.value
      )
      averageTimePerShift.value = response.data
      lastFetch.value['averageTimePerShift'] = Date.now()
    } catch (error) {
      console.error('Erreur fetch average time per shift:', error)
      throw error
    } finally {
      setLoading('averageTimePerShift', false)
    }
  }

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

  // === RETURN ===

  return {
    // States
    loading,
    individualPause,
    workingTimeIndividual,
    workingTimeTeam,
    presenceRate,
    averageTimePerShift,
    managerTeams,
    currentTeam,

    // Computed
    weekStartDate,
    weekEndDate,
    weekDisplayLabel,
    availableKpis,

    // Methods
    canAccessKpi,
    changeWeek,
    goToNextTeam,
    goToPreviousTeam,
    fetchManagerTeams,
    fetchIndividualPause,
    fetchWorkingTimeIndividual,
    fetchWorkingTimeTeam,
    fetchPresenceRate,
    fetchAverageTimePerShift,
    refreshAllKpis
  }
})