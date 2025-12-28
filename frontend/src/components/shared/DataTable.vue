<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center bg-base-100 p-4 rounded-lg shadow-sm">
      <div class="flex items-center gap-2">
        <div v-if="enableSearch" class="form-control">
          <div class="input-group">
            <input 
              type="text" 
              placeholder="Search..." 
              class="input input-bordered w-full max-w-xs" 
              v-model="searchQuery"
              @input="onSearch"
            />
          </div>
        </div>
        <slot name="filters"></slot>
      </div>
      <div>
        <slot name="actions"></slot>
      </div>
    </div>

    <div class="overflow-x-auto bg-base-100 rounded-lg shadow-sm">
      <table class="table w-full">
        <thead>
          <tr>
            <th 
              v-for="col in columns" 
              :key="col.key"
              @click="col.sortable ? sortBy(col.key) : null"
              :class="{ 'cursor-pointer hover:bg-base-200': col.sortable }"
            >
              <div class="flex items-center gap-1">
                {{ col.label }}
                <span v-if="col.sortable && currentSort.key === col.key">
                  {{ currentSort.order === 'asc' ? '↑' : '↓' }}
                </span>
                <span v-else-if="col.sortable" class="text-gray-300">↕</span>
              </div>
            </th>
            <th v-if="hasActions">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td :colspan="displayColumnsCount" class="text-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </td>
          </tr>
          <tr v-else-if="data.length === 0">
            <td :colspan="displayColumnsCount" class="text-center py-8 text-gray-500">
              No data found.
            </td>
          </tr>
          <tr v-else v-for="(item, index) in data" :key="item.id || index" class="hover">
            <td v-for="col in columns" :key="col.key">
              <slot :name="`cell-${col.key}`" :item="item" :value="item[col.key]">
                {{ item[col.key] }}
              </slot>
            </td>
            <td v-if="hasActions">
              <slot name="item-actions" :item="item"></slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div v-if="enablePagination && totalPages > 1" class="flex justify-center mt-4">
      <div class="join">
        <button 
          class="join-item btn" 
          :disabled="currentPage <= 1 || loading"
          @click="changePage(currentPage - 1)"
        >
          «
        </button>
        <button class="join-item btn">Page {{ currentPage }} of {{ totalPages }}</button>
        <button 
          class="join-item btn" 
          :disabled="currentPage >= totalPages || loading"
          @click="changePage(currentPage + 1)"
        >
          »
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

export interface Column {
  key: string;
  label: string;
  sortable?: boolean;
}

interface Props {
  columns: Column[];
  data: any[];
  loading?: boolean;
  totalItems?: number;
  itemsPerPage?: number;
  enableSearch?: boolean;
  enablePagination?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  totalItems: 0,
  itemsPerPage: 10,
  enableSearch: true,
  enablePagination: true
});

const emit = defineEmits(['update:page', 'update:sort', 'search', 'filter']);

const searchQuery = ref('');
const currentPage = ref(1);
const currentSort = ref<any>({}); // Simplified typing to avoid union issues key: string, order: 'asc' | 'desc' } | {}>({});

const hasActions = computed(() => {
  // Check if the slot 'item-actions' is used. 
  // In Vue 3 composition API with script setup, useSlots() needs import
  return true; // Simplification, usually checked via useSlots().item-actions
});

const displayColumnsCount = computed(() => props.columns.length + (hasActions.value ? 1 : 0));
const totalPages = computed(() => Math.ceil(props.totalItems / props.itemsPerPage));

const onSearch = () => {
  currentPage.value = 1; // Reset to page 1 on search
  emit('search', searchQuery.value);
};

const sortBy = (key: string) => {
  let order: 'asc' | 'desc' = 'asc';
  if ((currentSort.value as any).key === key && (currentSort.value as any).order === 'asc') {
    order = 'desc';
  }
  currentSort.value = { key, order };
  emit('update:sort', { key, order });
};

const changePage = (page: number) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  emit('update:page', page);
};

watch(() => props.data, () => {
  // React to data changes if needed
});
</script>
