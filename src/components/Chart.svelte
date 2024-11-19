
<script lang="ts">
  import { afterUpdate, onMount } from "svelte";
  import { Chart } from "chart.js";
  import type { ChartConfiguration } from "chart.js";

  export let config: ChartConfiguration;
  let canvas: HTMLCanvasElement;
  let chart: Chart;

  onMount(() => {
    chart = new Chart(canvas, config);
  })

  afterUpdate(() => {
    if (!chart) return;

    chart.data = config.data;
    Object.assign(chart.options, config.options);
    chart.update();
  });

</script>


<div class="w-full">
    <canvas bind:this={canvas}></canvas>
</div>
