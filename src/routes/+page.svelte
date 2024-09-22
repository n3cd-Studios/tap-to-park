<script lang="ts">
    import { onMount } from "svelte";
    import Button from "../components/Button.svelte";
    import Map from "../components/Map.svelte";
    import type { Coords } from "../lib/models";

    let map: L.Map;
    let coords: Coords= { latitude: 0, longitude: 0 };

    onMount(async () => {
        const promisifyGeolocation = (): Promise<Coords> =>
            new Promise((res, rej) => navigator.geolocation ? navigator.geolocation.getCurrentPosition(({ coords: { latitude, longitude } }) => res({ latitude, longitude })) : rej(null));
        coords = await promisifyGeolocation();
        map.flyTo([coords.latitude, coords.longitude], 10);
    });

    

</script>

<div class="flex h-full items-center justify-center">
    <div class="flex flex-col gap-2">
        <h1>You are at {coords.latitude}, {coords.longitude}</h1>
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
