<script lang="ts">
    import { get } from "$lib/api";
    import { Formats, pluralize } from "$lib/lang";
    import type { Spot } from "$lib/models";
    import moment from "moment";
    import Button from "../../components/form/Button.svelte";
    import { getUserInfo } from "$lib/auth";

    export let data: Spot;
    let continued = false;

    let hours = 0;
    let minutes = 0;

    const reservation = data.reservation;
    const costPerHour = data.price ?? 0;
    const costPerMinute = costPerHour / 60;
    $: cost = (costPerHour * hours) + (costPerMinute * minutes);
    $: purchaseDisabled = (hours == 0 && minutes == 0) || (hours < 0 || minutes < 0) || (Number(hours)*60 + Number(minutes) > data.maxHours*60) || (cost < 0.50);

    const checkout = async () => {
        const user = await getUserInfo();
        const session = await get<{ url: string }>({ route: `stripe/${data.guid}`, method: "POST", body: {
            start: moment().toISOString(),
            end: moment().add(hours, "hours").add(minutes, "minutes").toISOString(),
            user: user?.guid
        }});

        if (session) window.location.href = session.url;
    }

</script>

{#if continued}
    <div class="h-full flex flex-col justify-between items-center">
        <div class="m-10 flex flex-col justify-around h-1/2 items-start text-white text-lg font-bold" role="form" aria-labelledby="claim-space">
            <p>Claim this space for:</p>
            <div class="flex flex-col gap-2">
                <label for="hours" class="sr-only">Select number of hours</label>
                <div class="flex flex-row gap-2 items-baseline text-4xl">
                    <input bind:value={hours} class="bg-white rounded-lg text-black text-center w-1/4 py-5" max="24" type="number" aria-label="Enter number of hours"/>
                    <p>{hours == 1 ? "hour" : "hours"}</p>
                </div>
                <div class="flex flex-row gap-2 items-baseline text-4xl">
                    <input bind:value={minutes} class="bg-white rounded-lg text-black text-center w-1/4 py-5" max="60" aria-label="Enter number of minutes"/>
                    <p>{minutes == 1 ? "minute" : "minutes"}</p>
                </div>
            </div>
            {#if (hours == 0 && minutes == 0) || (hours < 0 || minutes < 0)}
                <p class="text-4xl">Enter a valid time</p>
            {:else if (Number(hours)*60 + Number(minutes) > data.maxHours*60)}
                <p class="text-4xl">Time exceeds the spot's maximum of {data.maxHours} hours</p>
            {:else}
                <p class="text-4xl">Total: {Formats.USDollar.format(Math.floor(cost*100)/100)}</p>
                {#if cost < 0.50}
                    <p class="text-2xl">Cost must be at least 50¢ to purchase</p>
                {/if}
            {/if}
        </div>
        <div class="mb-10">
            <Button on:click={() => checkout()} disabled={purchaseDisabled}  aria-label="Purchase the selected spot">Purchase</Button>
        </div>
    </div>
{:else}
    {#if reservation}
        <div class="h-full flex flex-col justify-between items-center" aria-live="polite">
            <div class="mt-10 flex flex-col justify-around h-1/2 items-center text-white text-lg font-bold" role="status">
                <p>This space is <span class="text-red-800">reserved</span>.</p>
                <p>This spot will be free in:</p>
                <div class="flex flex-row items-baseline">
                    <p class="text-7xl font-bold" aria-live="polite">{moment(reservation.end).fromNow(true)}</p>
                </div>
            </div>
        </div>
    {:else}
        <div class="h-full flex flex-col justify-between items-center" aria-live="polite">
            <div class="mt-10 flex flex-col justify-around h-1/2 items-center text-white text-lg font-bold" role="status">
                <p>This space is <span class="text-green-800">avaliable</span>.</p>
                <p>Claim this spot at the rate of:</p>
                <div class="flex flex-row items-baseline">
                    <p class="text-7xl font-bold">{Formats.USDollar.format(costPerHour)}</p>
                    <p>/hour</p>
                </div>
                <p>Maximum time: <span class="text-black">{pluralize(data.maxHours, "hour")}</span></p>
            </div>
            <div class="mb-10">
                <Button on:click={() => continued = true} aria-label="Continue to claim this spot">Continue</Button>
            </div>
        </div>
    {/if}
{/if}
