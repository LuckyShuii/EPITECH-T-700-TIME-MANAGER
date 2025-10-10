import type { UserLogin } from '@/types/userType';
import APIHandler from '../APIHandler';
const resource = 'authenticate';

export default {
    login(payload: UserLogin) {
        return APIHandler.post(`${resource}/`, payload);
    }, 

    logout(){
        return APIHandler.post(`${resource}/`);
    }
}