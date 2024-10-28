<script lang="ts">
  import { goto } from "$app/navigation";
  
  export let data: any[] = [];
  export let loading = true;
  export let columns: string[] = [];
  export let addRowItem: string | null = null;
  export let addRowFunctionality: (() => void) | null = null;

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
      {#if addRowItem && addRowFunctionality}
      <tr class="bg-blue-100 cursor-pointer" on:click={addRowFunctionality}>
        <td 
          class="px-4 sm:px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider bg-blue-100 cursor-pointer">
          Add {addRowItem}
        </td>
        {#each columns.slice(1) as column}
          <td 
            class="bg-blue-100 cursor-pointer" 
          >
        {/each}
      </tr>
    {/if}
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
</div>

