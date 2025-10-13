export interface UserLogin {
    username: string;
    password: string;
}

export interface UserBase{
    first_name : string; 
    email: string; 
    last_name : string, 
    phone_number: string
}

export interface UserProfile extends UserBase {
    roles: string[], 
    user_uuid: string, 
    username: string
}

