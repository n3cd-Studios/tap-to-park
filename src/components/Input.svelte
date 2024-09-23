
<script lang="ts">
    export let name: string;
    $: id = name.toLowerCase()
</script>

<div class="mb-4 flex flex-col w-max">
    <label class="block text-gray-700 text-sm font-bold mb-2" for={id}>{name}</label>
    <input class="appearance-none border rounded-lg px-2 py-3 focus:outline-none focus:shadow-outline" id={id} type="text" />
</div>