<script setup lang="ts">
import { onMounted, ref, watch, nextTick } from 'vue';
import Chart from 'chart.js/auto';

const props = defineProps<{
  data: {
    labels: string[];
    datasets: any[];
  };
  options?: any;
}>();

const chartCanvas = ref<HTMLCanvasElement | null>(null);
let chartInstance: Chart | null = null;

const renderChart = () => {
  if (!chartCanvas.value) return;
  if (chartInstance) chartInstance.destroy();

  chartInstance = new Chart(chartCanvas.value, {
    type: 'bar',
    data: props.data,
    options: {
      responsive: true,
      maintainAspectRatio: false,
      ...props.options
    }
  });
};

watch(() => props.data, () => {
  renderChart();
}, { deep: true });

onMounted(() => {
  nextTick(() => {
      renderChart();
  });
});
</script>

<template>
  <div class="w-full h-full">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>
