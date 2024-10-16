import { writable, type Writable } from "svelte/store";

export interface Toast {
    type?: "success" | "error";
    message: string;
}

const store = writable<Toast[]>([]);

export const toaster = {
    push: (toast: Toast, timeout = 3000) => {
        store.update(old => {
            old.push(toast);
            setTimeout(() => {
                store.update(list => {
                    const id = list.indexOf(toast);
                    if (id != -1) list.splice(id, 1);
                    return list;
                });
            }, timeout);
            return old;
        });
    },
    remove: (id: number) => store.update((old) => { old.splice(id, 1); return old; }),
    clear: () => store.set([]),
    subscribe: store.subscribe
}; 