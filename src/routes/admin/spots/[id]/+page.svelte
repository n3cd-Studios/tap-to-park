<script lang="ts">
    import { apiURL, get } from "$lib/api";
    import { authStore, getAuthHeader } from "$lib/auth";
    import { Region, type Point } from "$lib/geometry";
    import { daysOfWeek, Formats, properNoun } from "$lib/lang";
    import { ButtonType, IconType } from "$lib/utils";
    import { onMount } from "svelte";
    import Button from "../../../../components/form/Button.svelte";
    import Input from "../../../../components/form/Input.svelte";
    import { toaster } from "../../../../components/toaster/toaster";
    import type { Spot } from "$lib/models";
    import Modal from "../../../../components/Modal.svelte";
    import { goto } from "$app/navigation";

    // external
    export let data: Spot;

    // setup variables
    let maxHours = data.maxHours.toString();
    let name = data.name;

    // modal
    let deleting = false;

    // region management for price
    let price = "0";
    let region = new Region();
    let dragging = false;
    let times = Array(24).fill(0).map((_, num) => `${num % 12 === 0 ? 12 : num % 12}:00 ${num < 12 ? "AM" : "PM"}`);
    let schedule: number[][] = daysOfWeek.map((_) => times.map(_ => 0));

    onMount(async () => {
        if (!data) {
            toaster.push({ type: "error", message: "Failed to load spot information." });
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
    const namedItem = ([x, y]: Point) => `${properNoun(daysOfWeek[x])} at ${times[y]}`;

    const handleSave = async () => {
        updateItems(Number(price));
        await get<string>({
            route: `spots/${data.guid}`,
            headers: getAuthHeader(),
            body: {
              table: exportSchedule(),
              name,
              maxHours: Number(maxHours)
            },
            method: "PUT",
        });
        toaster.push({ type: "success", message: `Updated ${name}.` });
    };

    const handleDelete = async () => {
      await get<string>({
          route: `spots/${data.guid}`,
          headers: getAuthHeader(),
          body: {
            table: exportSchedule(),
            name,
            maxHours: Number(maxHours)
          },
          method: "DELETE",
      });
      toaster.push({ type: "error", message: `Deleted ${name}.` });
      await goto("/admin");
    }
</script>

<Modal
  visible={deleting}
  title={`Are you sure you want to delete "${name}?"`}
  icon={IconType.WARNING}
  on:close={() => deleting=false}
>
  <div slot="button" class="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6 gap-1" >
    <Button buttonType={ButtonType.CAUTION} aria-label="Cancel delete operation" on:click={() => deleting = false}>No</Button>
    <Button buttonType={ButtonType.CAUTION} aria-label="Confirm delete operation" on:click={handleDelete}>Yes</Button>
  </div>
</Modal>
<h1 class="text-xl font-bold text-center mb-2" aria-level="1">Managing "{name}" ({data.guid})</h1>
<div class="flex flex-col sm:flex-row gap-2">
    <div class="flex flex-col w-1/4">
        <Input bind:value={name} type="text" name="Name" id="name" aria-label="Name of the item"/>
        <Input bind:value={maxHours} type="number" name="Maximum reservation hours" id="max-hours" aria-label="Maximum reservation hours"/>
        <Input bind:value={price} type="number" name={`Price for ${namedItem(region.lower)} to ${namedItem(region.upper)}`} id="price" aria-label={`Price for reservations from ${namedItem(region.lower)} to ${namedItem(region.upper)}`}/>
        <Button on:click={handleSave} aria-label="Save">Save</Button>
        <div class="flex flex-col justify-center h-full">
            <p class="text-gray-700 text-sm font-bold">Spot QR Code</p>
            <img class="rounded-lg mt-2 w-2/3" src={apiURL`spots/${data.guid}/qr`} alt={`QR Code for managing ${name}`}/>
        </div>
        <Button buttonType={ButtonType.NEGATIVE} aria-label="Delete this item" on:click={() => deleting = true}>Delete</Button>
    </div>

    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
        class="grid grid-cols-8 bg-white text-center border-gray-200 border-2 w-full rounded-lg mb-10"
        aria-label="Reservation schedule"
        on:mouseenter={() => (dragging = false)}
    >
        <div class="grid grid-rows-12 border-r-2" role="rowgroup">
            <div class="bg-gray-200" role="rowheader">&nbsp;</div>
            {#each times as time}
                <div role="rowheader" aria-label={`Time: ${time}`}>{time}</div>
            {/each}
        </div>
        {#each daysOfWeek as day, x }
            <div class="grid grid-rows-12" role="columnheader" aria-label={day}>
                <h1 class="border-b-2 capitalize">{day}</h1>
                {#each times as _, y }
                    <!-- svelte-ignore a11y-mouse-events-have-key-events -->
                    <button
                        role="gridcell"
                        aria-label={`Price at ${daysOfWeek[x]} ${times[y]}: ${Formats.USDollar.format(schedule[x][y])}`}
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
