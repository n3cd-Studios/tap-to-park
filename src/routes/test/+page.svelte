<script lang="ts">
  import Table from '../../components/Table.svelte';
  import type { GetParams } from '$lib/api';

  let fetchParams: GetParams = {
    route: "admin/organization/invites",
    method: "GET",
    params: { guid: "8c6d4df8-31c5-493f-95b7-6b6e863958d9"}
  };
  
</script>
  
<Table {fetchParams}></Table>
  