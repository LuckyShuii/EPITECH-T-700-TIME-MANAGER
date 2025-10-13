import type { UserLogin } from '@/types/userType';
import APIHandler from '../APIHandler';
const resource = 'auth';

export default {
    login(payload: UserLogin) {
        return APIHandler.post(`${resource}/login`, payload);
    }, 

    logout(){
        return APIHandler.post(`${resource}/logout`);
    },

    getUserInfo() {
        return APIHandler.get(`${resource}/me`);
    }
}