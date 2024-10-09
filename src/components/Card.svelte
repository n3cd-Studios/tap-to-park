<script lang="ts">
    export let name: string;
    export let description: string;
    export let link: string = "#";
</script>

<a href={link} class="block max-w-sm p-6 bg-white border border-gray-200 rounded-lg shadow hover:bg-gray-100">
    <h5 class="mb-2 text-2xl font-bold tracking-tight text-black">{name}</h5>
    <p class="font-normal text-gray-700">{description}</p>
</a>
