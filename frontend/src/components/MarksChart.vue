<template>
  <div class="w-full h-64">
    <canvas ref="chartRef"></canvas>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import Chart from 'chart.js/auto';

const props = defineProps<{
  marks: any[]; // Expect { date: string, grade: number, subject: string }
}>();

const chartRef = ref<HTMLCanvasElement | null>(null);
let chartInstance: Chart | null = null;

const renderChart = () => {
  if (!chartRef.value) return;
  if (chartInstance) chartInstance.destroy();

  // Process data: Group by date or just show trend
  // Sort by date
  const sortedMarks = [...props.marks].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());
  
  const labels = sortedMarks.map(m => new Date(m.date).toLocaleDateString());
  const data = sortedMarks.map(m => m.grade);

  chartInstance = new Chart(chartRef.value, {
    type: 'line',
    data: {
      labels,
      datasets: [{
        label: 'Grades Trend',
        data,
        borderColor: '#4f46e5',
        tension: 0.1,
        fill: false
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        y: {
          min: 0,
          max: 10,
          beginAtZero: true
        }
      }
    }
  });
};

onMounted(renderChart);
watch(() => props.marks, renderChart, { deep: true });
</script>
