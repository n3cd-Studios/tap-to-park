<script lang="ts">
    import type { Toast } from "./toaster";
    import { toaster } from "./toaster";
    import { blur } from "svelte/transition";

    let toasts: Toast[] = [];
    toaster.subscribe((val) => (toasts = val));
</script>

<div>
    <div class="absolute m-14 flex flex-col gap-1 bottom-0">
        {#each toasts as { message, type = "success" }, index }
            <div transition:blur={{ amount: 1 }} class="px-2 py-4 bg-white text-black rounded-md">
                <button on:click={() => toaster.remove(index)}>x</button> {message}
            </div>
        {/each}
    </div>
</div>
