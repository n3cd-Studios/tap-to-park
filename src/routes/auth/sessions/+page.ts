import { get, getWithDefault } from "$lib/api";
import { authStore } from "$lib/auth";
import type { Session } from "$lib/models";
import { get as getStore } from "svelte/store";

export const load = async ({}) => {
    const sessions = await getWithDefault<Session[]>({
        route: `auth/sessions`,
        headers: { Authentication: `Bearer ${getStore(authStore).token}` },
        method: "GET",
    }, []);
    return { sessions };
} 