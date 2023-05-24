export interface User {
    id: string;
    username: string;
    name: string;
    location: string;
}

export interface RegistrationRequest {
    username: string;
    password: string;
    confirmPassword: string;
    name: string;
    location: string;
}