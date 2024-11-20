import type { Spot } from "$lib/models";

export const load = async ({ params, fetch }) => fetch(`spots/${params.id}`) as Spot;
