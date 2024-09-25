
<script lang="ts">
    import { get } from "$lib/api";
    import { authStore } from "$lib/auth";
    import type { Organization } from "$lib/models";
    import { onMount } from "svelte";

    $: email = $authStore.user?.email

    let organization: Organization | null;
    onMount(async () => {
        organization = await get<Organization>({ route: "admin/organization", headers: { "Authentication": `Bearer ${$authStore.token}` }, method: "GET" });
    })

</script>

<div class="m-10">
    <p>Hello, {email}! You are apart of <span class="font-bold">{organization?.name}</span>.</p>
    <!-- <div class="pt-3 flex flex-row gap-1">
        {#each organizations as organization}
            <Card name={organization.name} description="This is an organization" />
        {/each}
    </div> -->
</div>
