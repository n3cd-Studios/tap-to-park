
<script lang="ts">
    import type { User } from "$lib/models";
    import Table from "../../components/table/Table.svelte";
    import TableItem from "../../components/table/TableItem.svelte";

    export let data: User;
</script>

<div class="m-10">
    <Table
        data={[]}
        columns={["spot", "start", "end", "price"]}
        loading={false}
        let:guid
        let:device
        let:ip
        let:expires
        let:lastUsed
    >
    </Table>
</div>
