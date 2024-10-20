<script lang="ts">
    import Fa from "svelte-fa";
import type { Toast } from "./toaster";
    import { toaster } from "./toaster";
    import { blur } from "svelte/transition";
    import { faRemove } from "@fortawesome/free-solid-svg-icons";

    let toasts: Toast[] = [];
    toaster.subscribe((val) => (toasts = val));
</script>

<div>
    <div class="absolute m-14 flex flex-col gap-1 bottom-0 text-white">
        {#each toasts as { message, type = "success" }, index}
            <div class={`px-2 py-4 rounded-md flex gap-2 ${type == "success" ? "bg-green-500" : "bg-red-500"}`} transition:blur={{ amount: 1 }}>
                <button class="" on:click={() => toaster.remove(index)}><Fa icon={faRemove}/></button>
                {message}
            </div>
        {/each}
    </div>
</div>
