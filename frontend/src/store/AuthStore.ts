import { defineStore } from "pinia";
import type { UserLogin, UserProfile } from "@/types/userType";
import API from "@/services/API";
import { ref, computed } from "vue"


export const useAuthStore = defineStore("auth", () => {

    const user = ref<UserProfile | null>(null)
    const avatarColor = ref<string>('bg-primary')
    const isAuthenticated = computed(() => !!user.value)

    // Retrieve user roles (empty if not logged in)
    const userRoles = computed(() => user.value?.roles ?? [])

    // Check specific roles
    const isEmployee = computed(() => hasRole('employee'))
    const isManager = computed(() => hasRole('manager'))
    const isAdmin = computed(() => hasRole('admin'))

    // Work Session data
    const workSessionUuid = ref<string | null>(null)
    const clockInTime = ref<string | null>(null)
    const sessionStatus = ref<'active' | 'paused' | 'completed' | 'no_active_session'>('no_active_session')
    
    // Computed pour savoir si l'utilisateur est pointé
    const isClockedIn = computed(() => 
        sessionStatus.value === 'active' || sessionStatus.value === 'paused'
    )

    const initAuth = async () => {
        try {
            await fetchUserProfile()
            await fetchWorkSessionStatus()
            console.log('✅ Session restaurée depuis cookie')
        } catch (error) {
            console.log('ℹ️ Pas de session active')
            user.value = null
        }
    }

    const login = async (credentials: UserLogin) => {
        try {
            await API.authAPI.login(credentials)
            await fetchUserProfile()
            await fetchWorkSessionStatus()
            console.log("Connected")
        }
        catch (error) {
            user.value = null
            throw error
        }
    }

const fetchWorkSessionStatus = async () => {
    try {
        const response = (await API.WorkSession.getClockedStatus()).data
        
        console.log('📥 Réponse /status:', response) // AJOUTE ÇA
        
        if (response.status === 'no_active_session') {
            workSessionUuid.value = null
            clockInTime.value = null
            sessionStatus.value = 'no_active_session'
        } else {
            workSessionUuid.value = response.work_session_uuid
            clockInTime.value = response.clock_in_time
            sessionStatus.value = response.status
            
            console.log('✅ UUID stocké:', workSessionUuid.value) // AJOUTE ÇA
        }
    }
    catch (error) {
        console.error('Fetch work session status failed:', error)
        workSessionUuid.value = null
        clockInTime.value = null
        sessionStatus.value = 'no_active_session'
        throw error
    }
}

    const updateClocking = async (clockIn: boolean) => {
        try {
            await API.WorkSession.updateClocking(clockIn)
            // Recharger le statut après avoir pointé/dépointé
            await fetchWorkSessionStatus()
        }
        catch (error) {
            console.error('Update clocking failed:', error)
            throw error
        }
    }

    const updateBreaking = async (isBreaking: boolean) => {
        try {
            if (!workSessionUuid.value) {
                throw new Error('Pas de session active')
            }
            
            const response = await API.WorkSession.updateBreaking({
                is_breaking: isBreaking,
                work_session_uuid: workSessionUuid.value
            })
            
            // Mettre à jour le status selon la réponse
            if (response.data.status === 'break_started') {
                sessionStatus.value = 'paused'
            } else if (response.data.status === 'break_ended') {
                sessionStatus.value = 'active'
            }
        }
        catch (error) {
            console.error('Update breaking failed:', error)
            throw error
        }
    }

    const fetchUserProfile = async () => {
        try {
            const response = (await API.authAPI.getUserInfo()).data
            user.value = response
        }
        catch (error) {
            console.error('Profile fetching failed:', error)
            user.value = null
            throw error
        }
    }

    const logout = async () => {
        try {
            await API.authAPI.logout()
        }
        catch (error) {
            console.error('Erreur logout:', error)
        }
        finally {
            user.value = null
            workSessionUuid.value = null
            clockInTime.value = null
            sessionStatus.value = 'no_active_session'
            console.log('Disconected')
        }
    }

    const hasRole = (role: string): boolean => {
        if (!user.value) return false
        return userRoles.value.some(r => r.toLowerCase() === role.toLowerCase())
    }

    const hasAnyRole = (roles: string[]): boolean => {
        if (!user.value) return false
        return roles.some(role => hasRole(role))
    }

    const canAccess = (requiredRole: string): boolean => {
        if (!user.value) return false
        if (isAdmin.value) return true
        if (isManager.value && (requiredRole === 'manager' || requiredRole === 'employee')) {
            return true
        }
        if (isEmployee.value && requiredRole === 'employee') {
            return true
        }
        return false
    }

    const defaultRoute = computed(() => {
        if (isAdmin.value) return '/dashboard-admin'
        if (isManager.value) return '/dashboard-manager'
        return '/dashboard'
    })

    return {
        user,
        isAuthenticated,
        login,
        logout,
        fetchUserProfile,
        isClockedIn,
        updateClocking,
        fetchWorkSessionStatus,
        updateBreaking,
        initAuth,
        avatarColor,
        
        // Work session data
        workSessionUuid,
        clockInTime,
        sessionStatus,

        userRoles,
        isEmployee,
        isAdmin,
        isManager,
        hasRole,
        hasAnyRole,
        canAccess,
        defaultRoute
    }
})