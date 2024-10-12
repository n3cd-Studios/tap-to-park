<script lang="ts">
  import Table from '../../components/Table.svelte';
  import type { GetParams } from '$lib/api';

  let fetchParams: GetParams = {
    route: "spots/info",
    method: "GET",
    params: { guid: "a0889b39-d087-4b3a-a5b9-8ed7e2db41b6"}
  };

  let columnLabels = {
    guid: "Spot GUID",
    name: "Spot Name",
    handicap: "Handicap Access",
    organization: "Organization ID"
  };
  
</script>
  
<Table {fetchParams} defaultValue={[]} {columnLabels}></Table>
  