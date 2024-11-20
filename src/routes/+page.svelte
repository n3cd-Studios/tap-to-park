<script lang="ts">
    import { get, getWithDefault } from "$lib/api";
    import type { Marker } from "leaflet";
    import { onMount } from "svelte";
    import Button from "../components/form/Button.svelte";
    import Map from "../components/Map.svelte";
    import UserDropdown from "../components/userDropdown/UserDropdown.svelte";
    import type { Coords, PartialSpot, Spot, User } from "../lib/models";
    import Fa from 'svelte-fa';
    import { faBan, faWheelchair } from '@fortawesome/free-solid-svg-icons';
    import { Formats } from "$lib/lang";
    import moment from "moment";

    let map: L.Map;
    let spots: Marker<any>[] = [];
    let handicapFilter = false;

    const toggleHandicap = () => {
        handicapFilter = !handicapFilter;
        updateMarkers();
    }

    const updateMarkers = async () => {
        const promisifyGeolocation = (): Promise<Coords> =>
            new Promise((res, rej) => navigator.geolocation ? navigator.geolocation.getCurrentPosition(({ coords: { latitude, longitude } }) => res({ latitude, longitude })) : rej(null));

        const leaflet = await import("leaflet");
        const { latitude, longitude } = await promisifyGeolocation();
        map.setView([latitude, longitude], 13);

        if (spots.length > 0) {
            spots.forEach(marker => map.removeLayer(marker));
            spots = [];
        }

        const nearbySpots = await getWithDefault<PartialSpot[]>({ route: "spots/near", params: { lat: latitude.toString(), lng: longitude.toString(), handicap: handicapFilter.toString() }}, []);
        spots = nearbySpots.map(({ guid, coords, timeLeft = 0 }) => {
            const markerColor = timeLeft <= 0 ? 'green' : (timeLeft > 30 ? 'red' : 'yellow');
            const marker = leaflet
            .marker([coords.latitude, coords.longitude], {
                icon: leaflet.divIcon({
                    className: 'custom-marker',
                    html: `<div style="background-color: ${markerColor}; width: 16px; height: 16px; border-radius: 50%; border: 2px solid #fff;"></div>`,
                    iconSize: [16, 16],
                    iconAnchor: [8, 8],
                    popupAnchor: [0, -8]
                })
            })
                .bindPopup(`Loading`)
                .on("popupopen", async ({ popup }) => {
                    const spot = await get<Spot>({ route: `spots/${guid}` })
                    if (spot) {
                        if (spot.reservation) popup.setContent(
                           `<p>This spot is <span class="text-red-800">reserved</span>. It will become free in <span class="font-bold">${moment(spot.reservation.end).fromNow(true)}</span>.</p>`)
                        else popup.setContent(`
                            <p>This spot is <span class="text-green-800">available</span>, it costs <span class="font-bold">${Formats.USDollar.format(spot.price ?? 0)}</span></p>
                            <p>You can purchase this spot <a href="/${guid}">here</a>.</p>
                        `)
                    } else popup.setContent("Failed to load.");
                })
                .addTo(map)
            return marker;
        });
    }

    let timer: number;

    onMount(() => {
        updateMarkers();
        timer = setInterval(() => { updateMarkers(); }, 60 * 1000); // update markers every minute
        return () => { clearInterval(timer); };
    });

    let activeSpot = 0;
    const updateSpot = () => {
        if (activeSpot < 0) activeSpot = spots.length - 1;
        if (activeSpot >= spots.length) activeSpot = 0;
        spots[activeSpot].openPopup();
    };

</script>

<div class="flex h-full items-center justify-center" role="main" aria-label="Nearby parking spots interface">
    <div class="flex flex-col gap-2">
        <div class="w-96 h-96 rounded-lg border-white border-4" aria-label="Map showing nearby parking spots">
            <Map bind:map={map}/>
        </div>
        <div class="absolute top-6 right-12">
            <UserDropdown onLoginRedirect="/auth/login" />
        </div>
        <div class="flex flex-row justify-between">
            <div class="flex gap-2">
                <Button
                    aria-label="Find the nearest available parking spot"
                    on:click={() => { activeSpot = 0; updateSpot()}}>
                    Find Nearest
                </Button>
                <Button
                    aria-presed={handicapFilter}
                    aria-label="Toggle filter for handicap accessible spots"
                    on:click={toggleHandicap}>
                    <div class="relative">
                        <Fa icon={faWheelchair} class="w-4 h-4" />
                        {#if handicapFilter}
                            <Fa icon={faBan} class="absolute top-0 left-0 w-4 h-4 text-red-500 opacity-80" />
                        {/if}
                    </div>
                </Button>
            </div>
            <div>
                <Button
                    aria-label="Previous parking spot"
                    on:click={() => { activeSpot--; updateSpot(); }}>
                    {"<"}
                </Button>
                <Button
                    aria-label="Next parking spot"
                    on:click={() => { activeSpot++; updateSpot(); }}>
                    {">"}
                </Button>
            </div>
        </div>
    </div>
</div>
