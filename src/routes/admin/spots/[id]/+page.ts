import { get } from "$lib/api";
import { authStore } from "$lib/auth";
import type { Spot } from "$lib/models";
import type { Page } from "@sveltejs/kit";
import { get as getStore } from "svelte/store";

export const load = async ({ params }: Page) => {
    const spot = await get<Spot>({
        route: `spots/${params.id}`,
        headers: { Authentication: `Bearer ${getStore(authStore)}` },
        method: "GET",
    });
    return spot;
} 