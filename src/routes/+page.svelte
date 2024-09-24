<script lang="ts">
    import { onMount } from "svelte";
    import Button from "../components/Button.svelte";
    import Map from "../components/Map.svelte";
    import type { Coords, Spot } from "../lib/models";
    import { get, getWithDefault } from "$lib/api";
    import type { Marker } from "leaflet";

    let map: L.Map;
    let spots: Marker<any>[];

    onMount(async () => {
        const promisifyGeolocation = (): Promise<Coords> =>
            new Promise((res, rej) => navigator.geolocation ? navigator.geolocation.getCurrentPosition(({ coords: { latitude, longitude } }) => res({ latitude, longitude })) : rej(null));

        const leaflet = await import("leaflet");
        const { longitude, latitude } = await promisifyGeolocation();
        map.setView([latitude, longitude], 13);

        const nearbySpots = await getWithDefault<Spot[]>({ route: "spots/near", params: new URLSearchParams({ lng: longitude.toString(), lat: latitude.toString() })}, []);
        spots = nearbySpots.map(({ name, coords }) =>
            leaflet
                .marker([coords.longitude, coords.latitude])
                .bindPopup(`<strong>${name}</strong>`)
                .addTo(map),
        );
    });

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
            <Button click={() => console.log("Find nearest")}
                >Find Nearest</Button
            >
            <div>
                <Button click={() => { activeSpot--; updateSpot(); }}>{"<"}</Button>
                <Button click={() => { activeSpot++; updateSpot(); }}>{">"}</Button>
            </div>
        </div>
    </div>
</div>
