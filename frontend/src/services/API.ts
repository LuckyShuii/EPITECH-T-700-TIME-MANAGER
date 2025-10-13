import ExampleAPI from "@/services/routers/ExampleAPI";
import AuthAPI from "./routers/AuthAPI";
import UserAPI from "./routers/UserAPI";

export default {
    example: ExampleAPI,
    authAPI: AuthAPI,
    userAPI: UserAPI
}