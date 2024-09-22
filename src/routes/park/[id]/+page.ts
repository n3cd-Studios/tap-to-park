import type { Page } from "@sveltejs/kit";

export interface Spot {
    name: string
}

export const load = ({ params }: Page): Spot => {
    return {
        name: params.id // should be provided in the slug
    }; 
}