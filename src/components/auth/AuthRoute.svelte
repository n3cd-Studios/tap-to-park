<script lang="ts">
    import { browser } from "$app/environment";
    import { goto } from "$app/navigation";
    import { getUserInfo, logout } from "$lib/auth";
    import { onMount } from "svelte";

    export let authorized: boolean = true;
    export let signOut: boolean = true;
    export let redirect: string = "/auth/login";
    
    let pass = false;
    onMount(async () => {
        const user = await getUserInfo();
        if (browser) {
            if (authorized && !user) {
                if (signOut) logout();
                goto(redirect);
            } else pass = true;
        }
    });
</script>

{#if pass} <slot /> {/if}