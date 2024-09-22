
<script lang="ts">
    export let click = () => {};

</script>

<button on:click={click} class="py-2 px-4 text-white bg-black rounded-lg uppercase transition-colors duration-300 hover:bg-white hover:text-black">
    <slot />
</button>