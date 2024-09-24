
<script lang="ts">
    export let value: string;
    export let name: string;
    $: id = name.toLowerCase()
</script>

<div class="mb-4 flex flex-col w-full">
    <label class="block text-gray-700 text-sm font-bold mb-2" for={id}>{name}</label>
    <input bind:value={value} class="appearance-none border rounded-lg w-full px-2 py-2 leading-tight focus:outline-none focus:shadow-outline" id={id} type="text" />
</div>