<template>
  <div class="flex justify-center items-center p-4">
    <span :class="['loading', 'loading-spinner', sizeClass, colorClass]"></span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  size?: 'xs' | 'sm' | 'md' | 'lg';
  color?: 'primary' | 'secondary' | 'accent' | 'neutral' | 'info' | 'success' | 'warning' | 'error';
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  color: 'primary',
});

const sizeClass = computed(() => `loading-${props.size}`);
const colorClass = computed(() => `text-${props.color}`);
</script>
