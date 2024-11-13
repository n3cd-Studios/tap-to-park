<!-- modified code from https://flowbite.com/docs/forms/toggle/  -->
<script lang="ts">
   export let checked: boolean = false;
   export let name: string;
   $: id = name.toLowerCase()
</script>

<label class="inline-flex items-center cursor-pointer">
   <input bind:value={checked} type="checkbox" class="sr-only peer">
   <div class="relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-green-300 dark:peer-focus:ring-green-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-green-600"></div>
   <span class="ms-3 text-sm font-bold text-gray-700">{name}</span>
</label>