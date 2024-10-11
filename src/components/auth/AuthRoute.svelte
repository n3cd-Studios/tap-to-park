<script lang="ts">
    import { browser } from "$app/environment";
    import { goto } from "$app/navigation";
    import { authStore, getUserInfo, logout } from "$lib/auth";
    import { onMount } from "svelte";

    export let pass: boolean = !!$authStore.token; 
    export let signOut: boolean = true;
    export let redirect: string = "/auth/login";
    
    onMount(async () => {
        const user = await getUserInfo();
        if (browser && !(pass || user)) {
            if (signOut) {
                logout();
                goto("/auth/login");
            } else {
                goto(redirect);
            }
        }
    });
</script>

{#if pass} <slot /> {/if}