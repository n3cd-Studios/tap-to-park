<script lang="ts">
    import type { Spot } from "$lib/models";
    import { onMount } from "svelte";
    import Map from "../../../components/Map.svelte";
    import Button from "../../../components/form/Button.svelte";
    import { goto } from "$app/navigation";

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
    <div class="flex flex-col gap-2 text-center" role="status" aria-live="polite">
        <p>Successfully purchased spot!</p>
        {#if data.reservation}
            <p>Reservation expires at: {new Date(data.reservation.end).toLocaleString()}</p>
        {:else}
            <p>Reservation has expired</p>
        {/if}
        <div class="w-96 h-96 rounded-lg border-white border-4" role="region" aria-label="Map showing purchased spot location">
            <Map bind:map={map}/>
        </div>
        <Button on:click={() => goto("/")}>Main Page</Button>
    </div>
</div>
