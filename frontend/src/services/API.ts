import ExampleAPI from "@/services/routers/ExampleAPI";
import AuthAPI from "./routers/AuthAPI";
import UserAPI from "./routers/UserAPI";
import WorkSession from "./routers/WorkSessionAPI";
import TeamAPI from "./routers/TeamAPI";
import KpiAPI from "./routers/KpiAPI"


export default {
  example: ExampleAPI,
  authAPI: AuthAPI,
  userAPI: UserAPI,
  WorkSession: WorkSession,
  teamAPI: TeamAPI,
  kpiAPI: KpiAPI
}