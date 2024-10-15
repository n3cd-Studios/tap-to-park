<script lang="ts">
  import { Paginator } from "$lib/api";
  import { authStore } from "$lib/auth";
  import type { Invite } from "$lib/models";
  import { faCancel } from "@fortawesome/free-solid-svg-icons";
  import { onMount } from "svelte";
  import Fa from "svelte-fa";
  import Button from "../../../components/form/Button.svelte";
  import Table from "../../../components/table/Table.svelte";
  import TableItem from "../../../components/table/TableItem.svelte";

  let loading = true;
  let data: Invite[] = [];

  const paginator = new Paginator<Invite>(
    {
      route: "organization/invites",
      method: "GET",
      headers: { Authentication: `Bearer ${$authStore.token}` },
    },
    10,
  );
  paginator.subscribe((items) => (data = items));

  onMount(async () => {
    await paginator.load();
    loading = false;
  });
</script>

<Table
  columns={["code", "expiration", "cancel"]}
  {data}
  {loading}
  let:code
  let:expiration
>
  <TableItem>{code}</TableItem>
  <TableItem>{expiration}</TableItem>
  <TableItem
    ><button on:click={() => console.log("Cancel")}
      ><Fa icon={faCancel} /></button
    ></TableItem
  >
</Table>
<div class="flex flex-row justify-center gap-2">
  <Button on:click={() => paginator.last()}>Last</Button>
  <Button on:click={() => paginator.next()}>Next</Button>
</div>
