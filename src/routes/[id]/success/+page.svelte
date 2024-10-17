<script lang="ts">
    import { onMount } from "svelte";
    import Map from "../../../components/Map.svelte";
    import type { Spot } from "$lib/models";
    import { get } from "$lib/api";

    export let data: Spot;

    let map: L.Map;

    onMount(async () => {
        const leaflet = await import("leaflet");

        if (data) {
            const { longitude, latitude } = data.coords;
            leaflet
                .marker([longitude, latitude])
                .bindPopup(`${data.name}`)
                .addTo(map)
                .openPopup();
            map.setView([longitude, latitude], 13);
        }
    });


</script>

<div class="flex h-full items-center justify-center">
    <div class="flex flex-col gap-2 text-center">
        <p>Successfully purchased spot!</p>
        <div class="w-96 h-96 rounded-lg border-white border-4">
            <Map bind:map={map}/>
        </div>
    </div>
</div>
