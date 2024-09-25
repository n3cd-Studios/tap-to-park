<script lang="ts">
    import { onMount } from "svelte";
    import { browser } from "$app/environment";
    import { goto } from "$app/navigation";
    import { authStore, logout } from "$lib/auth";

    export let pass: boolean = !!$authStore.token; 
    export let logUserOut: boolean = true;
    export let redirect: string = "/auth/login";
    
    onMount(async () => {
        if (browser && !pass) {
            if (logUserOut) {
                logout();
                goto("/auth/login");
            } else goto(redirect);
        }
    });
</script>

{#if pass} <slot /> {/if}
