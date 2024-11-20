import { apiURL, get } from "$lib/api";
import { authStore, getAuthHeader } from "$lib/auth";
import type { Spot } from "$lib/models";

export const load = async ({ params, fetch }) => {
  const spot = await fetch(apiURL`spots/${params.id}`, {
    headers: getAuthHeader(),
    method: "GET"
  }).then(res => res.json() as Spot)
    .catch(_ => null);
  return spot;
}
