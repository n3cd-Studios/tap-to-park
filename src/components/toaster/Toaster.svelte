<script lang="ts">
    import type { Toast } from "./toaster";
    import { toaster } from "./toaster";
    import { blur } from "svelte/transition";

    let toasts: Toast[] = [];
    toaster.subscribe((val) => (toasts = val));
</script>

<div>
    <div class="absolute m-14 flex flex-col gap-1 bottom-0 text-white">
        {#each toasts as { message, type = "success" }, index}
            <div transition:blur={{ amount: 1 }}>
                {#if type == "success"}
                    <div class="px-2 py-4 rounded-md bg-green-500">
                        <button on:click={() => toaster.remove(index)}>x</button>
                        {message}
                    </div>
                {:else if type == "error"}
                    <div class="px-2 py-4 rounded-md bg-red-500">
                        <button on:click={() => toaster.remove(index)}>x</button>
                        {message}
                    </div>
                {/if}
            </div>
        {/each}
    </div>
</div>
