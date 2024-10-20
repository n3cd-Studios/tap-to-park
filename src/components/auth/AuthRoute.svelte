<script lang="ts">
    import { goto } from "$app/navigation";
    import { getUserInfo } from "$lib/auth";
    import { UserRole } from "$lib/models";
    import { onMount } from "svelte";

    export let roles: UserRole[] = [ UserRole.USER ];
    export let redirect = "/auth/login";
    
    let pass = false;
    onMount(async () => {
        const user = await getUserInfo();
        if (user) {
            if (!roles.includes(user.role)) goto(redirect);
            else pass = true;
        } else goto("/auth/login");
    });
</script>

{#if pass} <slot /> {/if}