
<script lang="ts">
    export let click: (() => void) | (() => Promise<void>) = () => {};

    type ButtonType = "submit" | "reset" | "button" | null | undefined; 
    export let type: ButtonType = "button"

</script>

<button on:click={click} {type} class="py-2 px-4 text-white bg-black rounded-lg uppercase transition-colors duration-300 hover:bg-white hover:text-black">
    <slot />
</button>