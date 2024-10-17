<script lang="ts">
    import { get } from "$lib/api";
    import { Formats, pluralize } from "$lib/lang";
    import type { Spot } from "$lib/models";
    import moment from "moment";
    import { onMount } from "svelte";
    import Button from "../../components/form/Button.svelte";

    export let data;

    let spot: Spot | null;
    let continued = false;

    let hours = 0;
    let minutes = 0;

    const costPerHour = 2;
    const costPerMinute = costPerHour / 60;
    $: price = (costPerHour * hours) + (costPerMinute * minutes);

    onMount(async () => {
        spot = await get<Spot>({ route: `spots/${data.id}` });
    })

    const checkout = async () => {
        const session = await get<{ url: string }>({ route: `stripe/${data.id}`, method: "POST", body: {
            start: moment().toISOString(),
            end: moment().add(hours, "hours").add(minutes, "minutes").toISOString()
        }});

        if (session) window.location.href = session.url;
    }

</script>

{#if spot}
    {#if continued}
        <div class="h-full flex flex-col justify-between items-center">
            <div class="m-10 flex flex-col justify-around h-1/2 items-start text-white text-lg font-bold">
                <p>Claim this space for:</p>
                <div class="flex flex-col gap-2">
                    <div class="flex flex-row gap-2 items-baseline text-4xl">
                        <input bind:value={hours} class="bg-white rounded-lg text-black text-center w-1/4 py-5" max="24"  />
                        <p>{hours == 1 ? "hour" : "hours"}</p>
                    </div>
                    <div class="flex flex-row gap-2 items-baseline text-4xl">
                        <input bind:value={minutes} class="bg-white rounded-lg text-black text-center w-1/4 py-5" max="60" />
                        <p>{minutes == 1 ? "minute" : "minutes"}</p>
                    </div>
                </div>
                <p class="text-4xl">Total: {Formats.USDollar.format(price)}</p>
            </div>
            <div class="mb-10">
                <Button on:click={() => checkout()}>Purchase</Button>
            </div>
        </div>
    {:else}
        <div class="h-full flex flex-col justify-between items-center">
            <div class="mt-10 flex flex-col justify-around h-1/2 items-center text-white text-lg font-bold">
                <p>This space is <span class="text-green-800">avaliable</span>.</p>
                <p>Claim this spot at the rate of:</p>
                <div class="flex flex-row items-baseline">
                    <p class="text-7xl font-bold">{Formats.USDollar.format(costPerHour)}</p>
                    <p>/hour</p>
                </div>
                <p>Maximum time: <span class="text-black">{pluralize(2, "hour")}</span></p>
            </div>
            <div class="mb-10">
                <Button on:click={() => continued = true}>Continue</Button>
            </div>
        </div>
    {/if}
{/if}
