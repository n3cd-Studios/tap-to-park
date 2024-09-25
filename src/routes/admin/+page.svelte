
<script lang="ts">
    import { get } from "$lib/api";
    import { authStore } from "$lib/auth";
    import type { Organization } from "$lib/models";
    import { onMount } from "svelte";
    import Map from "../../components/Map.svelte";

    $: email = $authStore.user?.email

    let map: L.Map;

    let organization: Organization | null;
    onMount(async () => {
        organization = await get<Organization>({ route: "organization/me", headers: { "Authentication": `Bearer ${$authStore.token}` }, method: "GET" });
        
        if (!organization) {
            return;
        }

        const leaflet = await import("leaflet");
        const { spots } = organization;
        spots.map(({ name, coords }) =>
            leaflet
                .marker([coords.longitude, coords.latitude])
                .bindPopup(`<strong>${name}</strong>`)
                .addTo(map),
        );
    
    })

</script>

<div class="m-10">
    <p>Hello, {email}! You are apart of <span class="font-bold">{organization?.name}</span>.</p>
    <div class="w-96 h-96 rounded-lg border-white border-4">
        <Map bind:map={map}/>
    </div>
</div>
