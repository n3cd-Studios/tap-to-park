<script lang="ts">
    import { get, Paginator } from "$lib/api";
    import { authStore, getAuthHeader } from "$lib/auth";
    import { type Organization, type Spot } from "$lib/models";
    import { faEdit } from "@fortawesome/free-solid-svg-icons";
    import { onMount } from "svelte";
    import Fa from "svelte-fa";
    import Button from "../../components/form/Button.svelte";
    import Map from "../../components/Map.svelte";
    import Table from "../../components/table/Table.svelte";
    import TableItem from "../../components/table/TableItem.svelte";
    import { toaster } from "../../components/toaster/toaster";
    import { goto } from "$app/navigation";
    import { IconType } from "$lib/utils";
    import { Formats } from "$lib/lang";
    import moment from "moment";

    let map: L.Map;

    let items: Spot[] = [];
    let paginator = new Paginator<Spot>(
        {
            route: "organization/spots",
            headers: getAuthHeader(),
            method: "GET",
        },
        7,
    );
    paginator.subscribe((cb) => (items = cb));

    let loading = true;
    let organization: Organization | null;

    const getSpot = (guid: string) =>
        get<Spot>({ route: `spots/${guid}` });

    let markers: [L.Marker, string][] = [];

    onMount(async () => {
        organization = await get<Organization>({
            route: "organization/me",
            headers: getAuthHeader(),
            method: "GET",
        });

        if (!organization) return;

        await paginator.load();
        loading = false;

        const leaflet = await import("leaflet");
        items.map(({ guid, coords }) => {
            const marker = leaflet
                .marker([coords.latitude, coords.longitude])
                .bindPopup("Loading...")
                .on("popupopen", async ({ popup }) => {
                    const spot = await get<Spot>({ route: `spots/${guid}` })
                    if (spot) {
                        if (spot.reservation) popup.setContent(
                           `<p>${spot.name} is <span class="text-red-800">reserved</span> with <span class="font-bold">${moment(spot.reservation.end).fromNow(true)}</span> remaining.</p>`)
                        else popup.setContent(`
                            <p>${spot.name} is <span class="text-green-800">available</span>.</p>
                            <p>You can purchase this spot <a href="/${guid}">here</a>.</p>
                        `)
                    } else popup.setContent("Failed to load.");
                })
                .addTo(map);

            markers.push([marker, guid])
        });
    });

</script>

<div class="flex flex-col gap-2">
    <h1 class="text-lg text-center" aria-level="1">
        <span class="font-bold">{organization?.name}</span> organization
    </h1>
    <div class="w-full h-96 rounded-lg border-white border-4 z-0" role="region" aria-label="Map showing organization locations">
        <Map bind:map />
    </div>
    <Table
        columns={["name", "coords", ""]}
        data={items}
        {loading}
        addRowItem={"spot"}
        addRowFunctionality={ () => goto("/admin/spots/create") }
        let:name
        let:coords
        let:guid
    >
        <TableItem>{name}</TableItem>
        <TableItem>({coords.latitude}, {coords.longitude})</TableItem>
        <TableItem><a class="inline-flex justify-center min-w-full hover:text-blue-600" href={`/admin/spots/${guid}`} aria-label={`Edit details for ${name}`}><Fa icon={faEdit} /></a></TableItem>
    </Table>
    <div class="flex flex-row justify-center gap-2">
        <Button on:click={() => paginator.last()} aria-label="Go to previous page">{"<"}</Button>
        <Button on:click={() => paginator.next()} aria-label="Go to next page">{">"}</Button>
    </div>
</div>
