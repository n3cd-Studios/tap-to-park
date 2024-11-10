<script lang="ts">
    import { get, Paginator } from "$lib/api";
    import { authStore } from "$lib/auth";
    import { type Organization, type Spot } from "$lib/models";
    import { faDollar, faTrash } from "@fortawesome/free-solid-svg-icons";
    import { onMount } from "svelte";
    import { getAuthHeader } from "$lib/auth";
    import Fa from "svelte-fa";
    import Button from "../../components/form/Button.svelte";
    import Map from "../../components/Map.svelte";
    import Modal from "../../components/Modal.svelte";
    import Table from "../../components/table/Table.svelte";
    import TableItem from "../../components/table/TableItem.svelte";
    import { toaster } from "../../components/toaster/toaster";
    import { goto } from "$app/navigation";
    import { IconType } from "$lib/utils";

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

    let markers: [L.Marker, string][] = [];

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
        items.map(({ guid, coords }) => {
            const marker = leaflet
                .marker([coords.latitude, coords.longitude])
                .bindPopup("Loading...")
                .on("popupopen", ({ popup }) =>
                    getSpot(guid).then((spot) =>
                        popup.setContent(`${spot?.name}`),
                    ),
                )
                .addTo(map);
                
            markers.push([marker, guid])
        });
    });

    let showModal = false;
    let deletionMessage = "";
    let spotToDelete = "";
    
    const confirmDeletionModal = (name: string, guid: string) => {
        deletionMessage = `Are you sure you want to delete spot: "${name}"`
        showModal = true;
        spotToDelete = guid
    }

    const deleteSpot = async () => {
        const response = await get<{ message: string }>({
            route: `spots/${spotToDelete}`,
            headers: getAuthHeader(),
            method: "DELETE",
        });
        if (!response) throw "Failed to delete spot.";
    }

    const handleSpotDeletion = async () => {
        try {
            await deleteSpot();
            console.log(spotToDelete)
            items = items.filter(item => item.guid !== spotToDelete);
            markers = markers.filter(([marker, spotGuid]) => {
                if (spotGuid === spotToDelete) {
                    map.removeLayer(marker);
                    return false;
                }
                return true;
            });
            toaster.push({ type: "success", message: "Spot deleted successfully." }, 5000);
        } catch (error) {
            toaster.push({ type: "error", message: "Failed to delete spot." }, 5000);
        } finally {
            showModal = false;
        }
    }
</script>

<div class="flex flex-col gap-2">
    <h1 class="text-lg text-center">
        <span class="font-bold">{organization?.name}</span> organization
    </h1>
    <div class="w-full h-96 rounded-lg border-white border-4 z-0">
        <Map bind:map />
    </div>
    <Table
        columns={["name", "coords", "manage spots"]}
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
        <TableItem>
            <div class="flex flex-row justify-between gap-2">
                <a href={`/admin/spots/${guid}`}><Fa class="hover:text-green-500" icon={faDollar} /></a>
                <button on:click={() => confirmDeletionModal(name, guid)} class="hover:text-red-500">
                    <Fa icon={faTrash} />
                </button>
            </div>
        </TableItem>
    </Table>
    <Modal
        visible={showModal}
        title={"Spot Deletion Confirmation"}
        message={deletionMessage}
        buttonText={"Confirm"}
        buttonFunctionality={handleSpotDeletion}
        icon={IconType.WARNING}
        on:close={() => showModal=false}
    >
    </Modal>
    <div class="flex flex-row justify-center gap-2">
        <Button on:click={() => paginator.last()}>{"<"}</Button>
        <Button on:click={() => paginator.next()}>{">"}</Button>
    </div>
</div>
