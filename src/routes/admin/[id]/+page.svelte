<script lang="ts">
    import { get } from "$lib/api";
    import { authStore } from "$lib/auth";
    import { Region, type Point } from "$lib/geometry";
    import { DayOfWeek, daysOfWeek, Formats } from "$lib/lang";
    import type { Spot } from "$lib/models";
    import { onMount } from "svelte";
    import Button from "../../../components/form/Button.svelte";
    import Input from "../../../components/form/Input.svelte";
    import { toaster } from "../../../components/toaster/toaster";

    interface ScheduleItem {
        day: DayOfWeek;
        times: {
            time: string;
            price: number;
            point: Point;
        }[];
    }

    // external
    export let data;

    // setup variables
    let price = "0";
    let region = new Region();
    let dragging = false;
    let times = Array(24).fill(0).map((_, num) => `${num % 12}:00`);
    let schedule: ScheduleItem[] = daysOfWeek.map((day, x) => ({ day, times: times.map((time, y) => ({ point: [x, y], time, price: 0 })) }));

    onMount(async () => {
        const spot = await get<Spot>({
            route: `spots/${data.id}`,
            headers: { Authentication: `Bearer ${$authStore.token}` },
            method: "GET",
        });
        if (!spot) {
            toaster.push({
                type: "error",
                message: "Failed to load spot information.",
            });
            return;
        }

        const pricing = spot.table;
        schedule.forEach(({ day, times }, i) => {
            const prices = pricing[day];
            if (prices) {
                times.forEach((_, j) => {
                    schedule[i].times[j].price = prices[j];
                });
            }
        })

    });

    // TODO: this is odd, maybe fix??
    const updateItems = (val: number) =>
        schedule.forEach((item, x) =>
            item.times.forEach((time, y) =>
                region.in(time.point)
                    ? (schedule[x].times[y].price = val)
                    : undefined,
            ),
        );

    const namedItem = ([x, y]: Point) => {
        let item = schedule[x];
        let time = item.times[y];
        return `${item.day} at ${time.time}`;
    };

    const handleSave = async () => {
        updateItems(Number(price));
        await get<string>({
            route: `spots/${data.id}`,
            headers: { Authentication: `Bearer ${$authStore.token}` },
            body: schedule.reduce((prev: any, item) => {
                prev[item.day] = item.times.map(time => time.price);
                return prev;
            }, {}),
            method: "PUT",
        });
    };
</script>

<h1 class="text-xl font-bold text-center">Pricing information for {data.id}</h1>
<div class="flex flex-col sm:flex-row gap-2">
    <div class="flex flex-col w-1/4">
        {#if region.size() >= 0}
            <Input
                bind:value={price}
                type="number"
                name={`Price for ${namedItem(region.lower)} to ${namedItem(region.upper)}`}
            />
            <Button on:click={handleSave}>Save</Button>
        {:else}
            <p>Select a region</p>
        {/if}
    </div>

    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
        class="grid grid-cols-8 bg-white text-center border-gray-200 border-2 w-full"
        on:mouseenter={() => (dragging = false)}
    >
        <div class="grid grid-rows-12 border-r-2">
            <div class="bg-gray-200">&nbsp;</div>
            {#each times as time}
                <div>{time}</div>
            {/each}
        </div>
        {#each schedule as { day, times }}
            <div class="grid grid-rows-12">
                <h1 class="border-b-2 capitalize">{day}</h1>
                {#each times as { point, price }}
                    <!-- svelte-ignore a11y-mouse-events-have-key-events -->
                    <button
                        on:mousedown={() => {
                            dragging = true;
                            region.lower = point;
                        }}
                        on:mouseover={() => {
                            if (dragging) region.upper = point;
                        }}
                        on:click={() => {
                            region.lower = point;
                            region.upper = point;
                        }}
                        on:mouseup={() => (dragging = false)}
                        class={`bg-${region.in(point) ? "green-500" : "white"} hover:bg-gray-400`}
                        >{Formats.USDollar.format(price)}</button
                    >
                {/each}
            </div>
        {/each}
    </div>
</div>
