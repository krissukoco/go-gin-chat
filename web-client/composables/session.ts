import { User } from "~/types/user";

export const useSessionToken = () => useState<string | null>('sessionToken', () => null); 
export const useUser = () => useState<User|null>(() => null);