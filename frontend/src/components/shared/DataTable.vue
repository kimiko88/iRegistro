<script setup lang="ts">
import { computed } from 'vue';

interface Column {
  key: string;
  label: string;
  sortable?: boolean;
}

const props = defineProps<{
  columns: Column[];
  items: any[];
  loading?: boolean;
  actions?: boolean;
}>();

const emit = defineEmits<{
  (e: 'action', name: string, item: any): void;
}>();
</script>

<template>
  <div class="overflow-x-auto">
    <table class="table table-zebra w-full bg-base-100 rounded-lg shadow">
      <thead>
        <tr>
          <th v-for="col in columns" :key="col.key">{{ col.label }}</th>
          <th v-if="actions">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="loading">
          <td :colspan="columns.length + (actions ? 1 : 0)" class="text-center py-4">
            <span class="loading loading-spinner loading-md"></span>
          </td>
        </tr>
        <tr v-else-if="items.length === 0">
           <td :colspan="columns.length + (actions ? 1 : 0)" class="text-center py-4 text-gray-500">
            No data available
          </td>
        </tr>
        <tr v-for="(item, idx) in items" :key="idx" v-else>
          <td v-for="col in columns" :key="col.key">
            <slot :name="`cell-${col.key}`" :item="item">
              {{ item[col.key] }}
            </slot>
          </td>
          <td v-if="actions">
            <slot name="actions" :item="item">
               <!-- Default Action Placeholder -->
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
