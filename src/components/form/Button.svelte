
<script lang="ts">
  export let color = "black";
</script>

<button on:click class="py-2 px-4 text-white bg-{color} rounded-lg uppercase transition-colors duration-300 hover:bg-white hover:text-black" {...$$restProps}>
    <slot />
</button>
