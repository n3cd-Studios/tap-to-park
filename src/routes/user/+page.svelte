
<script lang="ts">
    import type { Reservation, User } from "$lib/models";
    import { onMount } from "svelte";
    import Table from "../../components/table/Table.svelte";
    import TableItem from "../../components/table/TableItem.svelte";
    import { get } from "svelte/store";
    import { getWithDefault } from "$lib/api";
    import { getAuthHeader } from "$lib/auth";
    import { Formats } from "$lib/lang";

    export let data: User;

    let reservations: Reservation[] = [];
    onMount(async () => {
      reservations = await getWithDefault<Reservation[]>({ route: "reservations", headers: getAuthHeader() }, []);
    })

</script>

<div class="m-10">
    <Table
        data={reservations}
        columns={["spot", "start", "end", "price"]}
        loading={false}
        let:guid
        let:start
        let:end
        let:price
    >
        <TableItem>{guid}</TableItem>
        <TableItem>{Formats.Date(start)}</TableItem>
        <TableItem>{Formats.Date(end)}</TableItem>
        <TableItem>{Formats.USDollar.format(price / 100)}</TableItem>
    </Table>
</div>
