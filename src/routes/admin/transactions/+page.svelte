<script lang="ts">
    import { onMount } from "svelte";
    import Table from "../../../components/Table.svelte";

    let loading = true;
    let error: string;

    //TODO: implement backend portion

    onMount(() => {
        setTimeout(() => {
            try {
                loading = false;
            } catch (e) {
                error = "Failed to load transaction history.";
                loading = false;
            }
        }, 1000); // Simulate delay
    });

</script>

<div class="flex flex-col gap-4 p-4">
    <h1 class="text-2xl font-bold text-center">Transaction History</h1>
    
    {#if loading}
        <p class="text-center">Loading transactions...</p>
    {:else if error}
        <p class="text-center text-red-500">{error}</p>
    {:else}
        <p class="text-center"> X transactions found</p>
        <Table
            columns={['Date', 'Description', 'Amount', 'Type']}
            {error}
            {loading}
        />
    {/if}
</div>