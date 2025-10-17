import APIHandler from '../APIHandler';
const resource = 'work-session';

export default {
    getClockedStatus() {
        return APIHandler.get(`${resource}/status`);
    },
    updateClocking(isClockedIn: boolean) {
        return APIHandler.post(`${resource}/update-clocking`, {
            is_clocked: isClockedIn
        });
    },
    updateBreaking(data: { is_breaking: boolean, work_session_uuid: string }) {
        return APIHandler.post(`${resource}/update-breaking`, data);
    }
}