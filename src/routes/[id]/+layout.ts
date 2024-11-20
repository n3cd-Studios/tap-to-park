import { apiURL } from "$lib/api";
import type { Spot } from "$lib/models";

export const load = async ({ params, fetch }) =>
  fetch(apiURL`spots/${params.id}`)
    .then(res => res.json() as Spot)
    .catch(_ => null);
