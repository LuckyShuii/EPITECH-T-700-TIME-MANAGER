import axios from 'axios'
import { snakeToCamel } from '@/utils/snakeToCamel'

const baseDomain = import.meta.VITE_FRONTEND_URL;
const baseURL = `${baseDomain}/api/`;
const APIHandler = axios.create({
    baseURL,
    timeout: 10000,
});

APIHandler.interceptors.request.use((config: any) => {
    /**
     * Check if the user is authenticated and add the token to the request headers
     * before the request is sent to other services
     * 
     * if no token is found, the request is sent without the Authorization header
     * if a token exists, add to the request in the header Authorization: config.headers.Authorization = `Bearer <token>`
     */
});

APIHandler.interceptors.response.use(
    (response: any) => {
        return response
    },
    (error: any) => {
        if (error.response?.status === 401 || error.response?.data?.message === "Invalid or expired token") {
            /**
             * If the user is not authenticated, redirect to the login page
             * Log him out and clear the local storage
             */
        }

        return Promise.reject(error)
    }
)

export default APIHandler;