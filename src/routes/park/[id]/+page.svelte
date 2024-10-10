<script lang="ts">
    import { pluralize } from "$lib/lang";
    import Button from "../../../components/Button.svelte";

    export let data;
</script>

<div class="h-full flex flex-col justify-between items-center">
    <div class="mt-10 flex flex-col justify-around h-1/2 items-center text-white text-lg font-bold">
        <p>This space is <span class="text-green-800">avaliable</span>.</p>
        <p>Claim this spot at the rate of:</p>
        <div class="flex flex-row items-baseline">
            <p class="text-7xl font-bold">$0.75</p>
            <p>/hour</p>
        </div>
        <p>Maximum time: <span class="text-black">{pluralize(2, "hour")}</span></p>
    </div>
    <div class="mb-10">
        <Button>Continue</Button>
    </div>
</div>
