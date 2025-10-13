import { defineStore } from "pinia";
import type { UserLogin, UserProfile } from "@/types/userType";
import API from "@/services/API";
import { ref, computed } from "vue"

export const useAuthStore = defineStore("auth", () => {

    const user = ref<UserProfile | null>(null)
    const isAuthenticated = computed(() =>  ! !user.value)
    const isClockedIn = ref<boolean | undefined>(false)

    const login = async(credentials : UserLogin) => {
        try{

            //API create the cookie

            await API.authAPI.login(credentials)

            //Request catch the cookie 

            await fetchUserProfile()
            await isClocked()

            console.log("Connected")
        }
        catch(error){
            user.value = null
            throw error
        }
    }
    const isClocked = async() =>{

        try{
            const response = (await API.WorkSession.getClockedStatus()).data
            isClockedIn.value = response.is_clocked

        }
        catch(error){
            console.error('ClockedIN failed:', error)
            isClockedIn.value = undefined
            throw error
            
            
        }
    }

   const updateClocking = async(clockIn: boolean) =>{

    try{
        //Call POST with param

        await API.WorkSession.updateClocking(clockIn);
        isClockedIn.value = clockIn

    }
    catch(error){
        console.error('Update clocking failed:', error)
        throw error

    }

   }



    const fetchUserProfile = async() => {
        try{
            //cookie's sending auto by the browser

            const response = (await API.authAPI.getUserInfo()).data
            user.value = response
        }
        catch(error){
            console.error('Profile fetching failed:', error)
            user.value = null
            throw error
        }
    }

    const logout = async() => {
        try{
            //call the API : server side deleting cookie

            await API.authAPI.logout()
        }
        catch(error){
            console.error('Erreur logout:', error)
        }
        finally{
            user.value = null
            console.log('Disconected')
        }
    }

    return{
        user, 
        isAuthenticated, 
        login, 
        logout, 
        fetchUserProfile, 
        isClockedIn, 
        updateClocking, 
        isClocked
    }




})