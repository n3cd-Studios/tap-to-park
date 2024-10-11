import { writable } from "svelte/store";

export interface Toast {
    message: string;
}

const toaster = writable([]);