<script lang="ts">
    import { authStore, getAuthHeader } from "$lib/auth";
    import Button from "../../../../components/form/Button.svelte";
    import type { Coords } from "$lib/models";
    import { get } from "$lib/api";
    import { get as storeGet } from 'svelte/store';
    import Input from "../../../../components/form/Input.svelte";
    import Toggle from "../../../../components/form/Toggle.svelte";
    import { toaster } from "../../../../components/toaster/toaster";

    let name: string = '';
    let inputCoordinates = true;
    let latitude: string = '';
    let longitude: string = '';
    let price: string = '';
    let maxHours: string = '';
    let isHandicapped = false;

    const createSpot = async (name: string, latitude: number, longitude: number, price: number, maxHours: number, handicap: boolean) => {
        const response = await get({ route: "spots", method: "POST", headers: { "Authentication": `Bearer ${storeGet(authStore).token}` }, body: { name, coords: { latitude: latitude, longitude: longitude }, price, maxHours, handicap}});
        if (!response) throw "Failed to login.";
    }

    const promisifyGeolocation = (): Promise<Coords> =>
        new Promise((res, rej) => navigator.geolocation ? navigator.geolocation.getCurrentPosition(({ coords: { latitude, longitude } }) => res({ latitude, longitude })) : rej(null));

    const handleSpotCreation = async () => {
        if (!inputCoordinates){
            const coords = await promisifyGeolocation();
            latitude = String(coords.latitude);
            longitude = String(coords.longitude);
        }
        await createSpot(name, Number(latitude), Number(longitude), Number(price), Number(maxHours), Boolean(isHandicapped))
            .then( () => {
                toaster.push({ type: "success", message: `Spot "${name}" created successfully` }, 5000);
                name = ''; latitude =''; longitude=''; price=''; maxHours=''; isHandicapped=false; // clear form inputs
            })
            .catch(() => toaster.push({ type: "error", message: "Failed to create spot." }, 5000));
    };
</script>

<div class="h-full w-full flex flex-col justify-center items-center">
    <form class="p-10 bg-white rounded-xl w-1/3" aria-labelledby="spot-creation-form-title" on:submit|preventDefault={handleSpotCreation}>
        <h2 id="spot-creation-form" class="sr-only">Create a Parking Spot</h2>
        <Input bind:value={name} name="Name" required aria-label="Spot name"/>
        {#if inputCoordinates}
            <Input bind:value={latitude} type="number" step="0.000000000001" name="Latitude" id="latitude" required aria-label="Latitude of the spot"/>
            <Input bind:value={longitude} type="number" step="0.000000000001" name="Longitude" id="longitude" required aria-label="Longitude of the spot"/>
        {/if}
        <Input bind:value={price} type="number" step="0.01" name="Default Price" id="default-price" required aria-label="Default price for the spot"/>
        <Input bind:value={maxHours} type="number" name="Maximum Reservation Hours" id="max-reservation-hours" required aria-label="Maximum reservation hours"/>
        <Toggle bind:checked={isHandicapped} name="Handicapped"/>
            <div class="flex flex-row justify-between" aria-live="polite">  
            <Button type="button" on:click={() => inputCoordinates = !inputCoordinates} class="text-blue-800 underline" aria-label={inputCoordinates ? "Switch to use current location" : "Switch to specify coordinates"}>
                {inputCoordinates ? "Use My Current Location" : "Specify Coordinates"}
            </Button>
            <Button type="submit" aria-label="Create parking spot">Create Spot</Button>
        </div>
    </form>
</div>
