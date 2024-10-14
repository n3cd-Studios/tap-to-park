<script lang="ts">
    import { Region, type Point } from "$lib/geometry";
    import Button from "../../../components/form/Button.svelte";
    import Input from "../../../components/form/Input.svelte";

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

    // setup variables
    let price = "0";
    let region = new Region();
    let dragging = false;
    let times = Array(12).fill(0).map((_, num) => `${num}:00`);
    let schedule: ScheduleItem[] = daysOfWeek.map((day, x) => ({
        day,
        times: times.map((time, y) => ({ point: [x, y], time, price: 0 }))
    }));

    const startDragging = (point: Point) => { dragging = true; region.lower = point; };
    const whileDragging = (point: Point) => { if (dragging) region.upper = point; };
    const stopDragging = (point: Point) => { dragging = false; region.upper = point; };

    // TODO: this is goofy, maybe fix??
    const updatePricing = () => {
        schedule.forEach((item, i) => {
            item.times.forEach((time, j) => {
              if (region.in(time.point)) {
                schedule[i].times[j].price = Number(price); // hella goofy
              }  
            })
        })
    }

</script>

<div class="flex flex-col sm:flex-row gap-2">
    <div class="flex flex-col w-1/4">
        {#if region.size() >= 0}
            <Input bind:value={price} type="number" name="Price" />
            <Button on:click={updatePricing}>Update</Button>
        {:else}
            <p>Select a region</p>
        {/if}
    </div>
    <div class="grid grid-cols-8 bg-white text-center border-gray-200 border-2 w-full">
        <div>
            <div class="bg-gray-200">&nbsp;</div>
            {#each Array(12).fill(0).map((_, num) => `${num}:00`) as time}
                <div>{time}</div>
            {/each}
        </div>
        {#each schedule as { day, times }}
            <div class="grid grid-rows-12">
                <h1>{day}</h1>
                {#each times as { point, price }}
                    <button
                        on:mousedown={() => startDragging(point)}
                        on:mouseover={() => whileDragging(point)}
                        on:mouseup={() => stopDragging(point)}
                        class={`bg-${region.in(point) ? "green-500" : "white"} hover:bg-gray-400`}>{price}</button
                    >
                {/each}
            </div>
        {/each}
    </div>
</div>
