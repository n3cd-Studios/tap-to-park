
<script lang="ts">
  import { ButtonType } from "$lib/utils";

  export let buttonType: ButtonType = ButtonType.DEFAULT;

  const buttonStyles = {
    [ButtonType.POSITIVE]: 'bg-green-600 hover:bg-green-500 disabled:bg-gray-300 disabled:text-gray-600',
    [ButtonType.CAUTION]: 'bg-yellow-600 hover:bg-yellow-500 disabled:bg-gray-300 disabled:text-gray-600',
    [ButtonType.NEGATIVE]: 'bg-red-600 hover:bg-red-500 disabled:bg-gray-300 disabled:text-gray-600',
    [ButtonType.DEFAULT]: 'bg-gray-600 hover:bg-gray-500 disabled:bg-gray-300 disabled:text-gray-600',
  }
</script>

<button on:click class={`py-2 px-4 text-white rounded-lg uppercase transition-colors duration-300 ${buttonStyles[buttonType]}`} {...$$restProps}>
    <slot />
</button>
