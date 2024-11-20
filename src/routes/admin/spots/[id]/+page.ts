import { get } from "$lib/api";
import { authStore, getAuthHeader } from "$lib/auth";
import type { Spot } from "$lib/models";

export const load = async ({ params, fetch }) => {
  const spot = await fetch(`spots/${params.id}`, {
    headers: getAuthHeader(),
    method: "GET"
  }) as Spot;
  return spot;
}
