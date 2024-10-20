import { redirect } from '@sveltejs/kit';
import type Page from '../login/+page.svelte';
import { authStore } from '$lib/auth';

export const load = async ({ params, url }: Page) => {
    return {
        type: params.type,
        search: url.searchParams
    }
}