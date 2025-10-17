import { defineStore } from "pinia";
import type { UserLogin, UserProfile } from "@/types/userType";
import API from "@/services/API";
import { ref, computed } from "vue"


export const useAuthStore = defineStore("auth", () => {

    const user = ref<UserProfile | null>(null)
    const avatarColor = ref<string>('bg-primary')
    const isAuthenticated = computed(() => ! !user.value)

    // Retrieve user roles (empty if not logged in)
    const userRoles= computed(() => user.value?.roles ?? [])

    // Check specific roles
    const isEmployee = computed(()=> hasRole('employee'))
    const isManager = computed(() => hasRole('manager'))
    const isAdmin = computed(() => hasRole('admin'))


    const isClockedIn = ref<boolean | undefined>(false)


    const initAuth = async () => {
        try {
            await fetchUserProfile()
            await isClocked()
            console.log('✅ Session restaurée depuis cookie')
        } catch (error) {
            console.log('ℹ️ Pas de session active')
            user.value = null
        }
    }

    const login = async (credentials: UserLogin) => {
        try {

            //API create the cookie

            await API.authAPI.login(credentials)

            //Request catch the cookie 

            await fetchUserProfile()
            await isClocked()

            console.log("Connected")
        }
        catch (error) {
            user.value = null
            throw error
        }
    }
    const isClocked = async () => {
        return ''
        // try{
        //     const response = (await API.WorkSession.getClockedStatus()).data
        //     isClockedIn.value = response.is_clocked

        // }
        // catch(error){
        //     console.error('ClockedIN failed:', error)
        //     isClockedIn.value = undefined
        //     throw error


        // }
    }

    const updateClocking = async (clockIn: boolean) => {

        try {
            //Call POST with param

            await API.WorkSession.updateClocking(clockIn);
            isClockedIn.value = clockIn

        }
        catch (error) {
            console.error('Update clocking failed:', error)
            throw error

        }

    }

    const fetchUserProfile = async () => {
        try {
            //cookie's sending auto by the browser

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
            //call the API : server side deleting cookie

            await API.authAPI.logout()
        }
        catch (error) {
            console.error('Erreur logout:', error)
        }
        finally {
            user.value = null
            console.log('Disconected')
        }
    }

    const hasRole=(role:string):boolean=>{
        if(!user.value) return false
        return userRoles.value.some(r => r.toLowerCase()=== role.toLowerCase())
    }

    const hasAnyRole = (roles: string[]): boolean => {
        if (!user.value) return false
        return roles.some(role => hasRole(role))
    }

    const canAccess = (requiredRole:string):boolean =>{
        if(!user.value) return false

        if(isAdmin.value) return true

        if(isManager.value && (requiredRole === 'manager' || requiredRole === 'employee')){
            return true
        }
        if (isEmployee.value && requiredRole === 'employee'){
            return true
        }
         return false
    }

    const defaultRoute = computed(() => {
        if(isAdmin.value) return '/dashboard-admin'
        if(isManager.value) return '/dashboard-manager'
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
        isClocked,
        initAuth,
        avatarColor,

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