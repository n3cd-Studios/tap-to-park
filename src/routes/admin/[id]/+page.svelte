<script lang="ts">
    import { Region, type Point } from "$lib/geometry";
    import Button from "../../../components/form/Button.svelte";
    import Input from "../../../components/form/Input.svelte";
    import { authStore } from "$lib/auth";
    import type { Spot } from "$lib/models";
    import { get } from "$lib/api";
    import { onMount } from "svelte";
    import { Formats } from "$lib/lang";

    type DayOfWeek = "Sunday" | "Monday" | "Tuesday" | "Wednesday" | "Thursday" | "Friday" | "Saturday";
    let daysOfWeek: DayOfWeek[] = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
    
    interface TimeItem {
        time: string;
        price: number;
        point: Point;
    };

    interface ScheduleItem {
        day: DayOfWeek;
        times: TimeItem[];
    }

    // external
    export let data;

    // setup variables
    let price = "0";
    let region = new Region();
    let dragging = false;
    let times = Array(24).fill(0).map((_, num) => `${num}:00`);
    let schedule: ScheduleItem[] = daysOfWeek.map((day, x) => ({
        day,
        times: times.map((time, y) => ({ point: [x, y], time, price: 0 }))
    }));

    onMount(async () => {
        const spot = await get<Spot>({ route: `spots/${data.id}`, headers: { "Authentication": `Bearer ${$authStore.token}` }, method: "GET" });
        console.log(spot);
    })

    // TODO: this is odd, maybe fix??
    const usingSelection = (updater: (time: TimeItem, x: number, y: number) => void) => 
        schedule.forEach((item, x) => 
            item.times.forEach((time, y) => 
                region.in(time.point) ? updater(time, x, y) : undefined)
            );

    const named = ([x, y]: Point) => {
        let item = schedule[x];
        let time = item.times[y];
        return `${item.day} at ${time.time}`;
    }

</script>

<h1 class="text-xl font-bold text-center">Pricing information for {data.id}</h1>
<div class="flex flex-col sm:flex-row gap-2">
    <div class="flex flex-col w-1/4">
        {#if region.size() >= 0}
            <Input bind:value={price} type="number" name={`Price for ${named(region.lower)} to ${named(region.upper)}`} />
            <Button on:click={() => usingSelection((_, x, y) => schedule[x].times[y].price = Number(price))}>Save</Button>
        {:else}
            <p>Select a region</p>
        {/if}
    </div>
    
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="grid grid-cols-8 bg-white text-center border-gray-200 border-2 w-full" on:mouseenter={() => dragging = false}>
        <div class="grid grid-rows-12 border-r-2">
            <div class="bg-gray-200">&nbsp;</div>
            {#each times as time}
                <div>{time}</div>
            {/each}
        </div>
        {#each schedule as { day, times }}
            <div class="grid grid-rows-12">
                <h1 class="border-b-2">{day}</h1>
                {#each times as { point, price }}  
                    <!-- svelte-ignore a11y-mouse-events-have-key-events -->
                    <button
                        on:mousedown={() => { dragging = true; region.lower = point; }}
                        on:mouseover={() => { if (dragging) region.upper = point; }}
                        on:click={() => { region.lower = point; region.upper = point; }}
                        on:mouseup={() => dragging = false}
                        class={`bg-${region.in(point) ? "green-500" : "white"} hover:bg-gray-400`}>{Formats.USDollar.format(price)}</button
                    >
                {/each}
            </div>
        {/each}
    </div>
</div>
