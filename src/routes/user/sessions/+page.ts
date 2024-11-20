import { apiURL } from "$lib/api";
import { getAuthHeader } from "$lib/auth";
import type { Session } from "$lib/models";

export const load = async ({ fetch }) => {
  const sessions = await fetch(apiURL`auth/sessions`, { headers: getAuthHeader() })
    .then(res => res.json() as Session[])
    .catch(_ => []);
    return { sessions };
}
