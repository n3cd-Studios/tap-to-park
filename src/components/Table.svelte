<script lang="ts">
  import { onMount } from 'svelte';
  import { getWithDefault } from "$lib/api";
  import type { GetParams } from "$lib/api";

  // parent props
  export let fetchParams: GetParams;
  export let defaultValue: any[] = [];
  export let columnLabels: { [key: string]: string } = {};
  let data: any[] = [];
  let error: string | null = null;
  let loading = true;
  let columns: string[] = [];

  onMount(async () => {
    try {
      const result = await getWithDefault<any>(fetchParams, defaultValue);
      
      data = Array.isArray(result) ? result : [result];
      console.log(data);

      if (data.length > 0) {
        columns = Object.keys(data[0]); // Extract column names from the first object
      }
    } catch (err) {
      error = err instanceof Error ? err.message : "An unexpected error occurred.";
    } finally {
      loading = false;
    }
  });
</script>

<table class="min-w-full divide-y divide-gray-200">
  <thead class="bg-gray-50">
    <tr>
      {#each columns as column}
        <th
          scope="col"
          class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
        >
          {columnLabels[column] || column}
        </th>
      {/each}
    </tr>
  </thead>
  <tbody class="bg-white divide-y divide-gray-200">
    {#each data as row}
      <tr>
        {#each columns as column}
          <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
            {row[column] ?? 'N/A'}
          </td>
        {/each}
      </tr>
    {/each}
  </tbody>
</table>

{#if loading}
  <p>Loading...</p>
{/if}

{#if error}
  <p>Error: {error}</p>
{/if}
