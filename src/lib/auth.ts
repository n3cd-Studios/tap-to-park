
import { writable } from "svelte/store";

export interface AuthStore {
    token?: string
};

export const store = writable<AuthStore>({});