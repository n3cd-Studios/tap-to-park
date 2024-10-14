import type { Page } from "@sveltejs/kit";

export const load = ({ params }: Page) => {
    return {
        id: params.id
    }
} 