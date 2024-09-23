<script lang="ts">
    import { onMount } from "svelte";
    import Button from "../components/Button.svelte";
    import Map from "../components/Map.svelte";
    import type { Coords, Spot } from "../lib/models";
    import { get } from "$lib/api";

    let map: L.Map;

    onMount(async () => {
        const promisifyGeolocation = (): Promise<Coords> =>
            new Promise((res, rej) => navigator.geolocation ? navigator.geolocation.getCurrentPosition(({ coords: { latitude, longitude } }) => res({ latitude, longitude })) : rej(null));

        const leaflet = await import("leaflet");
        const { longitude, latitude } = await promisifyGeolocation();
        const spots = await get<Spot[]>({ route: "spots/near", params: new URLSearchParams({ lng: longitude.toString(), lat: latitude.toString() })}) ?? [];
        spots.forEach(({ coords }) => leaflet.marker([coords.longitude, coords.latitude]).addTo(map));
    });

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
                <Button>{"<"}</Button>
                <Button>{">"}</Button>
            </div>
        </div>
    </div>
</div>
