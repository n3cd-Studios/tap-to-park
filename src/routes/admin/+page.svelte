
<script lang="ts">
    import { get } from "$lib/api";
    import { authStore } from "$lib/auth";
    import type { Organization } from "$lib/models";
    import { onMount } from "svelte";
    import Map from "../../components/Map.svelte";
    import { pluralize } from "$lib/lang";

    $: email = $authStore.user?.email

    let map: L.Map;
    let spots: L.Marker<any>[] = [];

    let organization: Organization | null;
    onMount(async () => {
        organization = await get<Organization>({ route: "organization/me", headers: { "Authentication": `Bearer ${$authStore.token}` }, method: "GET" });
        
        if (!organization) {
            return;
        }

        const leaflet = await import("leaflet");
        spots = organization.spots.map(({ name, coords }) =>
            leaflet
                .marker([coords.longitude, coords.latitude])
                .bindPopup(`<strong>${name}</strong>`)
                .addTo(map),
        );
    
    })

</script>

<div class="m-10 flex flex-col gap-2">
    <p>Hello, {email}! You are apart of <span class="font-bold">{organization?.name}</span>.</p>
    <p>Your organization has {pluralize(spots.length, "spot")}:</p>
    <div class="w-full h-96 rounded-lg border-white border-4">
        <Map bind:map={map}/>
    </div>
</div>
