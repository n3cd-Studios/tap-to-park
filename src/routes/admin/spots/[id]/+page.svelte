<script lang="ts">
    import { get } from "$lib/api";
    import { authStore } from "$lib/auth";
    import { Region, type Point } from "$lib/geometry";
    import { daysOfWeek, Formats } from "$lib/lang";
    import { onMount } from "svelte";
    import Button from "../../../../components/form/Button.svelte";
    import Input from "../../../../components/form/Input.svelte";
    import { toaster } from "../../../../components/toaster/toaster";
    import type { Spot } from "$lib/models";

    // external
    export let data: Spot;

    // setup variables
    let price = "0";
    let region = new Region();
    let dragging = false;
    let times = Array(24).fill(0).map((_, num) => `${num % 12}:00 ${num / 12 < 1 ? "AM" : "PM"}`);
    let schedule: number[][] = daysOfWeek.map((_) => times.map(_ => 0));

    onMount(async () => {
        if (!data) {
            toaster.push({
                type: "error",
                message: "Failed to load spot information.",
            });
            return;
        }

        const pricing = data.table;
        schedule.forEach((inner, x) => {
            inner.forEach((_, y) => schedule[x][y] = pricing[daysOfWeek[x]][y]);
        });
    });

    // TODO: this is odd, maybe fix??
    const updateItems = (val: number) => schedule.forEach((inner, x) => inner.forEach((_, y) => region.in([x, y]) ? schedule[x][y] = val : undefined));
    const exportSchedule = () => schedule.reduce((prev: any, item, x) => { prev[daysOfWeek[x]] = item; return prev; }, {});
    const namedItem = ([x, y]: Point) => `${daysOfWeek[x].toUpperCase()} at ${times[y]}`;

    const handleSave = async () => {
        updateItems(Number(price));
        await get<string>({
            route: `spots/${data.guid}`,
            headers: { Authentication: `Bearer ${$authStore.token}` },
            body: exportSchedule(),
            method: "PUT",
        });
        toaster.push({
            type: "success",
            message: "Updated pricing schedule.",
        });
    };
</script>

<h1 class="text-xl font-bold text-center">Pricing information for {data.name}</h1>
<div class="flex flex-col sm:flex-row gap-2">
    <div class="flex flex-col w-1/4">
        <Input
                bind:value={price}
                type="number"
                name={`Price for ${namedItem(region.lower)} to ${namedItem(region.upper)}`}
        />
        <Button on:click={handleSave}>Save</Button>
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
        {#each daysOfWeek as day, x }
            <div class="grid grid-rows-12">
                <h1 class="border-b-2 capitalize">{day}</h1>
                {#each times as _, y }
                    <!-- svelte-ignore a11y-mouse-events-have-key-events -->
                    <button
                        on:mousedown={() => {
                            dragging = true;
                            region.lower = [x, y];
                        }}
                        on:mouseover={() => {
                            if (dragging) region.upper = [x, y];
                        }}
                        on:click={() => {
                            region.lower = [x, y];
                            region.upper = [x, y];
                        }}
                        on:mouseup={() => (dragging = false)}
                        class={`bg-${region.in([x, y]) ? "green-500" : "white"} hover:bg-gray-400`}
                        >{Formats.USDollar.format(schedule[x][y])}</button
                    >
                {/each}
            </div>
        {/each}
    </div>
</div>
