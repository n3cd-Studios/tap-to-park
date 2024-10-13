<script lang="ts">
  import Table from '../../../components/Table.svelte';
  import type { GetParams } from '$lib/api';
  import { authStore } from "$lib/auth";

  let fetchParams: GetParams = {
    route: "organization/invites",
    method: "GET",
    // params: { guid: "8c6d4df8-31c5-493f-95b7-6b6e863958d9"},
    headers: { "Authentication": `Bearer ${$authStore.token}` }
  };
  
</script>
  
<Table {fetchParams}></Table>