import type { UserBase } from "./userType";

export interface WorkSession{
    user : UserBase,
    ws_uuid : string, 
    clock_in : Date, 
    clock_out: Date, 
    duration_minutes : number,
    status : string

}