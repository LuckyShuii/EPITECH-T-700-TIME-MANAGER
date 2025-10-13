import APIHandler from '../APIHandler';
const resource = 'me';

export default {
    getUserInfo() {
        return APIHandler.get(`${resource}`);
    }
}