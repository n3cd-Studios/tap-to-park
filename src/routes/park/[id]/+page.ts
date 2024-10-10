import type { Page } from "@sveltejs/kit";
import type { PageData } from "./$types";

export const load = ({ params }: Page) => {
    return {
        id: params.id
    }
} 