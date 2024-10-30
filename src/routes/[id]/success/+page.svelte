<script lang="ts">
    import type { Spot } from "$lib/models";
    import { onMount } from "svelte";
    import Map from "../../../components/Map.svelte";

    export let data: Spot;

    let map: L.Map;

    onMount(async () => {
        const leaflet = await import("leaflet");

        if (data) {
            const { latitude, longitude } = data.coords;
            leaflet
                .marker([latitude, longitude])
                .bindPopup(`${data.name}`)
                .addTo(map)
                .openPopup();
            map.setView([latitude, longitude], 13);
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
