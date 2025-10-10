import { defineStore } from "pinia";
import type { UserLogin } from "@/types/userType";
import API from "@/services/API";
import { ref, computed } from "vue"
import { error } from "console";

export const useAuthStore = defineStore("auth", () => {

    const user = ref(null)
    const isAuthentificated = computed(() => ! !user.value)

    const login = async(credentials : UserLogin) => {
        try{

            //API create the cookie

            await API.authAPI.login(credentials)

            //Request catch the cookie 

            //await fetchUserProfile()

            console.log("Connected")
        }
        catch(error){
            user.value = null
            throw error
        }
    }

    // const fetchUserProfile = async() => {
    //     try{
    //         //cookie's sending auto by the browser

    //         const response = await API.userAPI.me()
    //         user.value = null
    //         throw error
    //     }
    //     catch(error){
    //         console.error('Profile fetching failed:', error)
    //         user.value = null
    //         throw error
    //     }
    // }

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
        isAuthentificated, 
        login, 
        logout, 
        //fetchUserProfile
    }




})