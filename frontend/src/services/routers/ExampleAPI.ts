import APIHandler from '../APIHandler';
const resource = 'example';

export default {
    getExample() {
        return APIHandler.get(`${resource}/`);
    }
}