<script lang="ts">
    import { get, getWithDefault } from "$lib/api";
    import type { Marker } from "leaflet";
    import { onMount } from "svelte";
    import Button from "../components/form/Button.svelte";
    import Map from "../components/Map.svelte";
    import type { Coords, Spot } from "../lib/models";
    import Fa from 'svelte-fa';
    import { faFilter } from '@fortawesome/free-solid-svg-icons';
    import { Formats } from "$lib/lang";
    import moment from "moment";
    import { getUserInfo } from "$lib/auth";

    let map: L.Map;
    let spots: Marker<any>[];
    let handicapFilter = false;
    let loggedIn = false;
    let userEmail: string | null = null;

    onMount(async () => {
        try {
            let user = await getUserInfo();
            if (user && user.email) {
                userEmail = user.email;
                loggedIn = true;
            } else {
                loggedIn = false;
            }
        } catch (error) {
            loggedIn = false;
            console.error("Failed to fetch user info:", error);
        }

        const promisifyGeolocation = (): Promise<Coords> =>
            new Promise((res, rej) => navigator.geolocation ? navigator.geolocation.getCurrentPosition(({ coords: { latitude, longitude } }) => res({ latitude, longitude })) : rej(null));

        const leaflet = await import("leaflet");
        const { latitude, longitude } = await promisifyGeolocation();
        map.setView([latitude, longitude], 13);

        const nearbySpots = await getWithDefault<Pick<Spot, "guid" | "coords">[]>({ route: "spots/near", params: { lat: latitude.toString(), lng: longitude.toString() }}, []);
        spots = nearbySpots.map(({ guid, coords }) =>
            leaflet
                .marker([coords.latitude, coords.longitude])
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
                .addTo(map),
        );
    });

    let activeSpot = 0;
    const updateSpot = () => {
        if (activeSpot < 0) activeSpot = spots.length - 1;
        if (activeSpot >= spots.length) activeSpot = 0;
        spots[activeSpot].openPopup();
    };

    const handleLoginButton = () => {
            if (loggedIn) {
                location.href = "/user";
            } else {
                location.href = "/auth/login";
            }
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
                <Button on:click={() => handicapFilter = !handicapFilter}>
                    <Fa icon={faFilter} class="text-white w-4 h-4 mr-0.5" />
                </Button>
            </div>
            <div>
                <Button on:click={() => { activeSpot--; updateSpot(); }}>{"<"}</Button>
                <Button on:click={() => { activeSpot++; updateSpot(); }}>{">"}</Button>
            </div>
            <div class="absolute top-0 right-0 p-4">
                <Button on:click={handleLoginButton}>
                    {#if loggedIn}
                        {userEmail}
                    {:else}
                        Login
                    {/if}
                </Button>
            </div>
        </div>
    </div>
</div>
