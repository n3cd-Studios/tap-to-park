<script>
    import Button from "../components/Button.svelte";
    import Map from "../components/Map.svelte";

</script>

<div class="flex h-full  items-center justify-center">
    <div class="flex flex-col gap-2">
        <div class="w-96 h-96 rounded-lg border-white border-4">
            <Map />
        </div>
        <div class="flex flex-row justify-between">
            <Button click={() => console.log("Find nearest")}>Find Nearest</Button>
            <div>
                <Button>{'<'}</Button>
                <Button>{'>'}</Button>
            </div>
        </div>
    </div>
</div>
