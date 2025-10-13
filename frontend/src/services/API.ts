import ExampleAPI from "@/services/routers/ExampleAPI";
import AuthAPI from "./routers/AuthAPI";
import UserAPI from "./routers/UserAPI";
import WorkSession from "./routers/WorkSessionAPI"

export default {
    example: ExampleAPI,
    authAPI: AuthAPI,
    userAPI: UserAPI,
    WorkSession:WorkSession
}