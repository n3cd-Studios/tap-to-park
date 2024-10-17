
import { persisted } from 'svelte-persisted-store'
import { get } from "./api";
import type { User } from "./models";
import { get as storeGet } from 'svelte/store';

export interface AuthStore {
    token?: string,
    user?: User,
};

export interface TokenResponse {
    token: string;
}

export const authStore = persisted<AuthStore>("auth", {});

export const login = async (email: string, password: string) => {
    const response = await get<TokenResponse>({ route: "auth/login", method: "POST", body: { email, password } });
    if (!response) throw "Failed to login.";

    const { token } = response;
    const user = await get<User>({ route: "auth/info", headers: { "Authentication": `Bearer ${token}` }, method: "GET" });
    if (!user) throw "Failed to login.";

    authStore.set({ token, user });
}

export const getAuthHeader = () => ({ "Authentication": `Bearer ${storeGet(authStore).token}` })

export const getUserInfo = async () => {
    return get<User>({ route: "auth/info", headers: getAuthHeader() });
}

export const logout = () => {
    authStore.set({}); // clear
}