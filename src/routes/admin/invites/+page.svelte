<script lang="ts">
  import { onMount } from 'svelte';
  import Table from '../../../components/table/Table.svelte';
  import { authStore } from '$lib/auth';
  import { getWithDefault } from '$lib/api';
    import type { Invite } from '$lib/models';
    import TableItem from '../../../components/table/TableItem.svelte';

  let loading = true;
  let error: string;
  let data: Invite[] = [];

  onMount(async () => {
    data = await getWithDefault<Invite[]>({ route: "organization/invites", method: "GET", headers: { "Authentication": `Bearer ${$authStore.token}` }}, []);
    if (!data) {
      error = "Failed to load table."
      return;
    }
    loading = false;
  });
  
</script>
  
<Table columns={["code", "expiration"]} {data} {loading} {error} 
  let:code let:expiration>
  <TableItem>{code}</TableItem>
  <TableItem>{expiration}</TableItem>
</Table>