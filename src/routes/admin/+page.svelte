
<script lang="ts">
    import { get } from "$lib/api";
    import { authStore } from "$lib/auth";
    import { type Spot, type Organization } from "$lib/models";
    import { onMount } from "svelte";
    import Map from "../../components/Map.svelte";
    import { pluralize } from "$lib/lang";
    import Table from "../../components/Table.svelte";

    $: email = $authStore.user?.email;

    let map: L.Map;
    let spots: L.Marker<any>[] = [];

    let loading = true;
    let error: string;
    let organization: Organization | null;

    const getSpot = (guid: string) => get<Spot>({ route: `spots/${guid}/info` })

    onMount(async () => {
        organization = await get<Organization>({ route: "organization/me", headers: { "Authentication": `Bearer ${$authStore.token}` }, method: "GET" });
        if (!organization) {
            error = "Failed to load spots for organization.";
            return;
        }
        loading = false;

        const leaflet = await import("leaflet");
        spots = organization.spots.map(({ guid, coords }) =>
            leaflet
                .marker([coords.longitude, coords.latitude])
                .bindPopup("Loading...")
                .on('popupopen', ({ popup }) => getSpot(guid).then(spot => popup.setContent(`${spot?.name}`)))
                .addTo(map),
        );
    
    })

</script>

<div class="flex flex-col gap-2">
    <h1 class="text-lg text-center"><span class="font-bold">{organization?.name}</span> organization</h1>
    <div class="w-full h-96 rounded-lg border-white border-4">
        <Map bind:map={map}/>
    </div>
    <Table 
        columns={["name", "coords"]} 
        data={organization?.spots
            .map(({ name, coords }) => 
                ({ 
                    name, 
                    coords: `(${coords.longitude}, ${coords.latitude})` 
                })
            )} 
        {error} 
        {loading} />
</div>
