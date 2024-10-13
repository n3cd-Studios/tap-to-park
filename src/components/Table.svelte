<script lang="ts">
  import { onMount } from 'svelte';
  import { getWithDefault } from "$lib/api";
  import type { GetParams } from "$lib/api";


  // Parent props
  export let fetchParams: GetParams;
  export let defaultValue: any[] = [];

  let data: any[] = [];
  let error: string | null = null;
  let loading = true;
  let columns: string[] = [];

  onMount(async () => {
    try {
      const result = await getWithDefault<any>(fetchParams, defaultValue);

      data = Array.isArray(result) ? result : [result];

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

<div class="h-full rounded-t-[3rem] p-4">
  <div class="max-w-5xl mx-auto">
    <table class="min-w-full divide-y divide-gray-200 rounded-lg overflow-hidden">
      <thead class="bg-gray-50">
        <tr>
          {#each columns as column}
            <th
              scope="col"
              class="px-4 sm:px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider"
            >
              {column}
            </th>
          {/each}
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        {#each data as row}
          <tr>
            {#each columns as column}
              <td class="px-4 sm:px-6 py-4 text-sm text-gray-500 whitespace-nowrap">
                {row[column] ?? 'N/A'}
              </td>
            {/each}
          </tr>
        {/each}
      </tbody>
    </table>

    {#if loading}
      <p class="mt-4 text-center text-gray-500">Loading...</p>
    {/if}

    {#if error}
      <p class="mt-4 text-center text-red-500">Error: {error}</p>
    {/if}
  </div>
</div>

