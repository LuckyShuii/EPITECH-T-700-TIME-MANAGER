export interface UserLogin {
    username: string;
    password: string;
}

export interface UserProfile{
    first_name : string; 
    email: string; 
    last_name : string, 
    phone_number: string,
    roles: string[], 
    user_uuid: string, 
    username: string
}
