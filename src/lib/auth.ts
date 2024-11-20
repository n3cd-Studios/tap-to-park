
import { persisted } from 'svelte-persisted-store'
import { apiURL, get } from "./api";
import type { User } from "./models";
import { get as storeGet } from 'svelte/store';

export interface AuthStore {
    token?: string
};

export interface TokenResponse {
    token: string;
}

export const authStore = persisted<AuthStore>("auth", {});

export const login = async (email: string, password: string) => {
    const response = await get<TokenResponse>({ route: "auth/login", method: "POST", body: { email, password } });
    if (!response) throw "Failed to login.";
    authStore.set({ token: response.token });
    return await getUserInfo();
}

export const register = async (email: string, password: string, invite: string = "") => {
    const response = await get<TokenResponse>({ route: "auth/register", method: "POST", body: { email, password }, params: { invite } });
    if (!response) throw "Failed to register you.";
    authStore.set({ token: response.token });
    return await getUserInfo();
}

export const getAuthHeader = () => ({ "Authentication": `Bearer ${storeGet(authStore).token}` })
export const getUserInfo = async (fetcher = fetch) =>
  fetcher(apiURL`auth/info`, { headers: getAuthHeader() })
    .then(async res => await res.json() as User)
    .catch(_ => null);

export const logout = () => {
    authStore.set({}); // clear
}
