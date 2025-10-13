import APIHandler from '../APIHandler';
const resource = 'work-session';

export default {
    getClockedStatus() {
        return APIHandler.get(`${resource}/is-clocked`);
    }, 
    updateClocking(isClockedIn:boolean){
        return APIHandler.post(`${resource}/update-clocking`,{
            is_clocked:isClockedIn
        });
    }
}