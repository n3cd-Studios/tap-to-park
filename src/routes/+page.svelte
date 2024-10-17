<script lang="ts">
    import { getWithDefault } from "$lib/api";
    import type { Marker } from "leaflet";
    import { onMount } from "svelte";
    import Button from "../components/form/Button.svelte";
    import Map from "../components/Map.svelte";
    import type { Coords, Spot } from "../lib/models";
    import Fa from 'svelte-fa';
    import { faFilter } from '@fortawesome/free-solid-svg-icons';

    let map: L.Map;
    let spots: Marker<any>[];
    let handicapFilter = false;

    onMount(async () => {
        const promisifyGeolocation = (): Promise<Coords> =>
            new Promise((res, rej) => navigator.geolocation ? navigator.geolocation.getCurrentPosition(({ coords: { latitude, longitude } }) => res({ latitude, longitude })) : rej(null));

        const leaflet = await import("leaflet");
        const { longitude, latitude } = await promisifyGeolocation();
        map.setView([latitude, longitude], 13);

        const nearbySpots = await getWithDefault<Spot[]>({ route: "spots/near", params: { lng: longitude.toString(), lat: latitude.toString() }}, []);
        spots = nearbySpots.map(({ guid, name, coords }) =>
            leaflet
                .marker([coords.longitude, coords.latitude])
                .bindPopup(`<a href="/${guid}">${name}</a>`)
                .addTo(map),
        );
    });

    // TODO: logic to handle filter button
    const toggleHandicapFilter = async () => {
        handicapFilter = !handicapFilter;
    };

    let activeSpot = 0;
    const updateSpot = () => {
        if (activeSpot < 0) activeSpot = spots.length - 1;
        if (activeSpot >= spots.length) activeSpot = 0;
        spots[activeSpot].openPopup();
    };

</script>

<div class="flex h-full items-center justify-center">
    <div class="flex flex-col gap-2">
        <div class="w-96 h-96 rounded-lg border-white border-4">
            <Map bind:map={map}/>
        </div>
        <div class="flex flex-row justify-between">
            <div class="flex gap-2">
                <Button on:click={() => { activeSpot = 0; updateSpot()}}>Find Nearest</Button>
                <Button on:click={toggleHandicapFilter}>
                    <Fa icon={faFilter} class="text-white w-4 h-4 mr-0.5" />
                </Button>
            </div>
            <div>
                <Button on:click={() => { activeSpot--; updateSpot(); }}>{"<"}</Button>
                <Button on:click={() => { activeSpot++; updateSpot(); }}>{">"}</Button>
            </div>
        </div>
    </div>
</div>
