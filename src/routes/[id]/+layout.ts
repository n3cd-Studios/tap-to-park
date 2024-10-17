import { get } from "$lib/api";
import type { Spot } from "$lib/models";
import type { Page } from "@sveltejs/kit";

export const load = async ({ params }: Page) => get<Spot>({ route: `spots/${params.id}` });