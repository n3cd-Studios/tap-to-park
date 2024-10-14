<script lang="ts">
    import TableItem from "./TableItem.svelte";

  
  export let data: any[] = [];
  export let error: string | null = null;
  export let loading = true;
  export let columns: string[] = [];

</script>

<div class="h-full w-full rounded-t-[3rem]">
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
          <slot {...row} />
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

