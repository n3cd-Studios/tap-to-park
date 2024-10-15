<script lang="ts">
    import { get, Paginator } from "$lib/api";
    import { authStore } from "$lib/auth";
    import { type Organization, type Spot } from "$lib/models";
    import { faDollar } from "@fortawesome/free-solid-svg-icons";
    import { onMount } from "svelte";
    import Fa from "svelte-fa";
    import Button from "../../components/form/Button.svelte";
    import Map from "../../components/Map.svelte";
    import Table from "../../components/table/Table.svelte";
    import TableItem from "../../components/table/TableItem.svelte";
    import { toaster } from "../../components/toaster/toaster";

    let map: L.Map;

    let items: Spot[] = [];
    let paginator = new Paginator<Spot>(
        {
            route: "organization/spots",
            headers: { Authentication: `Bearer ${$authStore.token}` },
            method: "GET",
        },
        7,
    );
    paginator.subscribe((cb) => (items = cb));

    let loading = true;
    let organization: Organization | null;

    const getSpot = (guid: string) =>
        get<Spot>({ route: `spots/${guid}/info` });

    onMount(async () => {
        organization = await get<Organization>({
            route: "organization/me",
            headers: { Authentication: `Bearer ${$authStore.token}` },
            method: "GET",
        });
        if (!organization) {
            toaster.push({
                type: "error",
                message: "Failed to load spots for organization.",
            });
            return;
        }

        await paginator.load();
        loading = false;

        const leaflet = await import("leaflet");
        items.map(({ guid, coords }) =>
            leaflet
                .marker([coords.longitude, coords.latitude])
                .bindPopup("Loading...")
                .on("popupopen", ({ popup }) =>
                    getSpot(guid).then((spot) =>
                        popup.setContent(`${spot?.name}`),
                    ),
                )
                .addTo(map),
        );
    });
</script>

<div class="flex flex-col gap-2">
    <h1 class="text-lg text-center">
        <span class="font-bold">{organization?.name}</span> organization
    </h1>
    <div class="w-full h-96 rounded-lg border-white border-4">
        <Map bind:map />
    </div>
    <Table
        columns={["name", "coords", "manage pricing"]}
        data={items}
        {loading}
        showAddSpot={true}
        let:name
        let:coords
        let:guid
    >
        <TableItem>{name}</TableItem>
        <TableItem>({coords.longitude}, {coords.latitude})</TableItem>
        <TableItem><a href={`/admin/${guid}`}><Fa icon={faDollar} /></a></TableItem>
    </Table>
    <div class="flex flex-row justify-center gap-2">
        <Button on:click={() => paginator.last()}>{"<"}</Button>
        <Button on:click={() => paginator.next()}>{">"}</Button>
    </div>
</div>
