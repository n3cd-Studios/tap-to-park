<script lang="ts">
    export let click: (() => void) | (() => Promise<void>) = () => {}; 

</script>

<button on:click={click} class="py-2 px-2 bg-black rounded-lg transition-colors duration-300 hover:bg-white">
    <img src="src/assets/filter.svg" alt="Filter Icon" class="w-4 h-4 mr-0.5" />
    <slot />
</button>