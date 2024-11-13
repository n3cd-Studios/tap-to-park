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
    <form class="p-10 bg-white rounded-xl w-1/3" on:submit|preventDefault={handleSpotCreation}>
        <Input bind:value={name} name="Name" required/>
        {#if inputCoordinates}
            <Input bind:value={latitude} type="number" step="0.000000000001" name="Latitude" required/>
            <Input bind:value={longitude} type="number" step="0.000000000001" name="Longitude" required/>
        {/if}
        <Input bind:value={price} type="number" step="0.01" name="Default Price" required/>
        <Input bind:value={maxHours} type="number" name="Maximum Reservation Hours" required/>
        <Toggle bind:checked={isHandicapped} name="Handicapped"/>
            <div class="flex flex-row justify-between">  
            <Button type="button" on:click={() => inputCoordinates = !inputCoordinates} class="text-blue-800 underline">
                {inputCoordinates ? "Use My Current Location" : "Specify Coordinates"}
            </Button>
            <Button type="submit">Create Spot</Button>
        </div>
    </form>
</div>
