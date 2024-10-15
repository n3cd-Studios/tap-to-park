<script lang="ts">
    import { authStore } from "$lib/auth";
    import Button from "../../../../components/form/Button.svelte";
    import type { Coords } from "$lib/models";
    import { get } from "$lib/api";
    import { get as storeGet } from 'svelte/store';
    import Input from "../../../../components/form/Input.svelte";
    import { toaster } from "../../../../components/toaster/toaster";

    let name: string = '';
    let longitude: string = '';
    let latitude: string = '';
    let inputCoordinates = true;

    function toggleCoordinates() {
        inputCoordinates = !inputCoordinates;
    }

    export const createSpot = async (name: string, longitude: number, latitude: number) => {
        const response = await get({ route: "spots/create", method: "POST", headers: { "Authentication": `Bearer ${storeGet(authStore).token}` }, body: { name, coords: { longitude: longitude, latitude: latitude }}});
        if (!response) throw "Failed to login.";
    }

    const promisifyGeolocation = (): Promise<Coords> =>
        new Promise((res, rej) => navigator.geolocation ? navigator.geolocation.getCurrentPosition(({ coords: { latitude, longitude } }) => res({ latitude, longitude })) : rej(null));

    const handleSpotCreation = async () => {
        if (!inputCoordinates){
            const coords = await promisifyGeolocation();
            longitude = String(coords.longitude);
            latitude = String(coords.latitude);
        }
        await createSpot(name, Number(longitude), Number(latitude))
            .then( () => {
                toaster.push({ type: "success", message: `Spot "${name}" created successfully` }, 5000);
                name = ''; longitude = ''; latitude = ''; // clear form inputs
            })
            .catch(() => toaster.push({ type: "error", message: "Failed to create spot." }, 5000));
    };
</script>

<div class="h-full w-full flex flex-col justify-center items-center">
    <form class="p-10 bg-white rounded-xl w-1/3" on:submit|preventDefault={handleSpotCreation}>
        <Input bind:value={name} name="Name" required/>
        {#if inputCoordinates}
            <Input bind:value={longitude} type="number" step="0.000000000001" name="Longitude" required/>
            <Input bind:value={latitude} type="number" step="0.000000000001" name="Latitude" required/>
        {/if}
            <div class="flex flex-row justify-between">  
            <Button type="button" on:click={toggleCoordinates} class="text-blue-800 underline">
                {inputCoordinates ? "Use My Current Location" : "Specify Coordinates"}
            </Button>
            <Button type="submit">Create Spot</Button>
        </div>
    </form>
</div>