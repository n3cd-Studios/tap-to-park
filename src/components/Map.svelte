<script lang="ts">
    import { onMount, onDestroy } from "svelte";

    let mapElement: HTMLElement;
    export let map: L.Map;

    onMount(async () => {
        const leaflet = await import("leaflet");

        map = leaflet.map(mapElement).setView([51.505, -0.09], 13);

        leaflet
            .tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
                attribution:
                    '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
            })
            .addTo(map);

        leaflet
            .marker([51.5, -0.09])
            .addTo(map)
            .bindPopup("A pretty CSS3 popup.<br> Easily customizable.")
            .openPopup();
    });

    onDestroy(async () => {
        if (map) {
            console.log("Unloading Leaflet map.");
            map.remove();
        }
    });
</script>

<div class="w-full h-full" bind:this={mapElement}></div>